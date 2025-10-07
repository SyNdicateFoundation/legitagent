package legitagent

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	utls "github.com/refraction-networking/utls"
)

var (
	ErrUnsupportedBrowser = errors.New("unsupported browser could not be parsed")
	ErrUnsupportedOS      = errors.New("unsupported os could not be parsed")
	ErrUnsupportedVersion = errors.New("no suitable profile found for the browser version")

	uaRegexes = []struct {
		Browser Browser
		Regex   *regexp.Regexp
	}{
		{BrowserSafari, regexp.MustCompile(`Version/(\d+)\..*Safari/`)},
		{BrowserEdge, regexp.MustCompile(`Edg/(\d+)\.\d+`)},
		{BrowserOpera, regexp.MustCompile(`OPR/(\d+)\.\d+`)},
		{BrowserBrave, regexp.MustCompile(`Brave/(\d+)\.\d+`)},
		{BrowserChrome, regexp.MustCompile(`Chrome/(\d+)\.\d+`)},
		{BrowserFirefox, regexp.MustCompile(`Firefox/(\d+)\.\d+`)},
	}
)

var parserStableChromeProfiles = map[int]utls.ClientHelloID{
	120: utls.HelloChrome_120,
	131: utls.HelloChrome_131,
	133: utls.HelloChrome_133,
}
var parserStableChromeVersions []int

func init() {
	if len(parserStableChromeVersions) == 0 {
		for v := range parserStableChromeProfiles {
			parserStableChromeVersions = append(parserStableChromeVersions, v)
		}
		sort.Ints(parserStableChromeVersions)
	}
}

func findClosestChromeProfileForParser(targetVersion int) utls.ClientHelloID {
	bestVersion := -1
	for _, v := range parserStableChromeVersions {
		if v <= targetVersion && v > bestVersion {
			bestVersion = v
		}
	}
	if bestVersion != -1 {
		return parserStableChromeProfiles[bestVersion]
	}
	return parserStableChromeProfiles[parserStableChromeVersions[0]]
}

type parsedUA struct {
	Browser Browser
	Version int
	OS      OperatingSystem
}

func parseUserAgentString(ua string) (*parsedUA, error) {
	var p parsedUA

	for _, re := range uaRegexes {
		if match := re.Regex.FindStringSubmatch(ua); len(match) > 1 {
			p.Browser = re.Browser
			v, err := strconv.Atoi(match[1])
			if err != nil {
				return nil, fmt.Errorf("could not parse version from ua string: %w", err)
			}
			p.Version = v
			break
		}
	}
	if p.Browser == "" {
		return nil, ErrUnsupportedBrowser
	}

	switch {
	case strings.Contains(ua, "Windows NT 10.0"):
		p.OS = OSWindows11
	case strings.Contains(ua, "iPhone"), strings.Contains(ua, "iPad"):
		p.OS = OSiOS
	case strings.Contains(ua, "Macintosh"):
		p.OS = osMacIntel
	case strings.Contains(ua, "Linux"):
		p.OS = OSLinux
	case strings.Contains(ua, "Android"):
		p.OS = OSAndroid
	case strings.Contains(ua, "CrOS"):
		p.OS = OSChromeOS
	default:
		return nil, ErrUnsupportedOS
	}

	return &p, nil
}

func FromUserAgentString(userAgentString string, requestType RequestType) (*Agent, error) {
	ua, err := parseUserAgentString(userAgentString)
	if err != nil {
		return nil, err
	}

	profile, ok := browserProfiles[ua.Browser]
	if !ok {
		return nil, ErrUnsupportedBrowser
	}

	osProf, ok := osProfiles[ua.OS]
	if !ok {
		return nil, ErrUnsupportedOS
	}

	platform := PlatformDesktop

	if ua.OS == OSAndroid || ua.OS == OSiOS {
		platform = PlatformMobile
	}

	platformProf, ok := platformProfiles[platform]
	if !ok {
		return nil, fmt.Errorf("internal error: no platform profile for %s", platform)
	}

	versionProf, _, err := findClosestVersionProfile(profile.Versions, ua.Version)
	if err != nil {
		return nil, err
	}

	fullVersion := ""
	if profile.ChromiumBased {
		fullVersion = fmt.Sprintf("%d.0.%d.0", ua.Version, versionProf.BuildNumber)
	}

	headers, headerOrder := buildStaticHeaders(profile, osProf, platformProf, ua.Version, fullVersion, versionProf, requestType)

	helloID := findClosestChromeProfileForParser(ua.Version)

	return &Agent{
		UserAgent:       userAgentString,
		Headers:         headers,
		HeaderOrder:     headerOrder,
		ClientHelloSpec: nil,
		ClientHelloID:   helloID,
		H2Settings:      GetChromiumH2Settings(),
	}, nil
}

func findClosestVersionProfile(versions map[int]versionProfile, targetVersion int) (versionProfile, int, error) {
	closestVersion := -1
	for v := range versions {
		if v <= targetVersion && v > closestVersion {
			closestVersion = v
		}
	}

	if closestVersion == -1 {
		return versionProfile{}, 0, ErrUnsupportedVersion
	}

	return versions[closestVersion], closestVersion, nil
}

func buildStaticHeaders(browser browserProfile, os osProfile, platform platformProfile, version int, fullVersion string, versionProf versionProfile, requestType RequestType) (http.Header, []string) {
	headerMap := make(map[string]string)

	sb := builderPool.Get().(*strings.Builder)
	defer func() {
		sb.Reset()
		builderPool.Put(sb)
	}()

	var acceptTemplate []AcceptHeaderPart
	switch requestType {
	case RequestTypeXHR:
		if len(versionProf.AcceptHeaderPatternsXHR) > 0 {
			acceptTemplate = versionProf.AcceptHeaderPatternsXHR[0]
		}
	default:
		if len(versionProf.AcceptHeaderPatterns) > 0 {
			acceptTemplate = versionProf.AcceptHeaderPatterns[0]
		}
	}

	for i, part := range acceptTemplate {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(part.Value)
		for _, extra := range part.Extras {
			sb.WriteString(";")
			sb.WriteString(extra)
		}
		if part.Q > 0 {
			sb.WriteString(";q=")
			sb.WriteString(strconv.FormatFloat(part.Q, 'f', 1, 64))
		}
	}

	headerMap["accept"] = sb.String()
	sb.Reset()

	headerMap["accept-encoding"] = "gzip, deflate, br"
	headerMap["accept-language"] = "en-US,en;q=0.9"

	if browser.ChromiumBased {
		headerMap["sec-ch-ua"] = buildSecChUa(browser.Brand, strconv.Itoa(version), false, false)
		headerMap["sec-ch-ua-mobile"] = platform.MobileHint
		headerMap["sec-ch-ua-platform"] = fmt.Sprintf(`"%s"`, os.Name)
		headerMap["sec-ch-ua-full-version-list"] = buildSecChUa(browser.Brand, fullVersion, true, false)
		if os.Version != "" {
			headerMap["sec-ch-ua-platform-version"] = fmt.Sprintf(`"%s"`, os.Version)
		}
		if os.Arch != "" {
			headerMap["sec-ch-ua-arch"] = fmt.Sprintf(`"%s"`, os.Arch)
		}
		if os.BitnessHint != "" {
			headerMap["sec-ch-ua-bitness"] = fmt.Sprintf(`"%s"`, os.BitnessHint)
		}
	}

	switch requestType {
	case RequestTypeNavigate:
		headerMap["sec-fetch-dest"] = "document"
		headerMap["sec-fetch-mode"] = "navigate"
		headerMap["sec-fetch-site"] = "none"
		headerMap["sec-fetch-user"] = "?1"
		headerMap["upgrade-insecure-requests"] = "1"
	case RequestTypeXHR:
		headerMap["sec-fetch-dest"] = "empty"
		headerMap["sec-fetch-mode"] = "cors"
		headerMap["sec-fetch-site"] = "same-origin"
	}

	header := http.Header{}
	keys := make([]string, 0, len(headerMap))
	for k := range headerMap {
		keys = append(keys, k)
	}

	PriorityHeaderSorter(keys)
	orderedKeys := append([]string{":method", ":authority", ":scheme", ":path"}, keys...)
	for _, k := range keys {
		header.Set(k, headerMap[k])
	}

	return header, orderedKeys
}
