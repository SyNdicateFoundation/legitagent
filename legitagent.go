package legitagent

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/SyNdicateFoundation/fastrand"
	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
)

type Agent struct {
	UserAgent       string
	Headers         http.Header
	HeaderOrder     []string
	ClientHelloSpec *utls.ClientHelloSpec
	ClientHelloID   utls.ClientHelloID
	H2Settings      map[http2.SettingID]uint32
}

type Generator struct {
	browsers               []Browser
	platforms              []Platform
	os                     []OperatingSystem
	minVersion             int
	maxVersion             int
	languageProfiles       [][]AcceptHeaderPart
	requestType            RequestType
	headerSorter           HeaderSorter
	fullFingerprint        bool
	h2Only                 bool
	fingerprintProfile     FingerprintProfile
	h2RandomizationProfile H2RandomizationProfile
	useBotAgents           bool
	botAgentTypes          []string
	acceptEncodingEnabled  bool
	acceptEnabled          bool
	agentPool              sync.Pool
	zeroHeader             bool
}

var allRealOS = []OperatingSystem{
	OSWindows,
	OSWindows11,
	OSLinux,
	OSMac,
	OSAndroid,
	OSiOS,
	OSChromeOS,
	osUbuntu,
	osFedora,
}
var allRealBrowsers = []Browser{
	BrowserChrome,
	BrowserOpera,
	BrowserEdge,
	BrowserBrave,
	BrowserFirefox,
	BrowserSafari,
}
var allRealPlatforms = []Platform{
	PlatformDesktop,
	PlatformMobile,
}
var macArchitectures = []OperatingSystem{
	osMacIntel,
	osMacAppleSilicon,
}

var builderPool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

func NewGenerator(opts ...Option) *Generator {
	defaultLanguages := [][]AcceptHeaderPart{
		{{Value: "en-US"}, {Value: "en", Q: 0.9}},
		{{Value: "de-DE"}, {Value: "de", Q: 0.9}},
		{{Value: "fa-IR"}, {Value: "fa", Q: 0.9}},
		{{Value: "fr-FR"}, {Value: "fr", Q: 0.9}},
		{{Value: "es-ES"}, {Value: "es", Q: 0.9}},
		{{Value: "ja-JP"}, {Value: "ja", Q: 0.9}},
		{{Value: "ko-KR"}, {Value: "ko", Q: 0.9}},
		{{Value: "pt-BR"}, {Value: "pt", Q: 0.9}},
		{{Value: "ru-RU"}, {Value: "ru", Q: 0.9}},
		{{Value: "tr-TR"}, {Value: "tr", Q: 0.9}},
		{{Value: "it-IT"}, {Value: "it", Q: 0.9}},
		{{Value: "pl-PL"}, {Value: "pl", Q: 0.9}},
		{{Value: "nl-NL"}, {Value: "nl", Q: 0.9}},
		{{Value: "sv-SE"}, {Value: "sv", Q: 0.9}},
		{{Value: "ar-EG"}, {Value: "ar", Q: 0.9}},
		{{Value: "cs-CZ"}, {Value: "cs", Q: 0.9}},
	}

	g := &Generator{
		browsers:               []Browser{BrowserRandom},
		platforms:              []Platform{PlatformRandom},
		os:                     []OperatingSystem{OSRandom},
		minVersion:             114,
		maxVersion:             141,
		languageProfiles:       defaultLanguages,
		requestType:            RequestTypeNavigate,
		headerSorter:           PriorityHeaderSorter,
		fullFingerprint:        false,
		h2Only:                 true,
		fingerprintProfile:     FingerprintProfileNormal,
		h2RandomizationProfile: H2RandomizationProfileNone,
		useBotAgents:           false,
		botAgentTypes:          nil,
		acceptEncodingEnabled:  false,
		acceptEnabled:          true,
	}

	g.agentPool.New = func() any {
		return new(Agent)
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

func (g *Generator) Generate() (*Agent, error) {
	agent := g.agentPool.Get().(*Agent)
	agent.Headers = make(http.Header)

	if g.useBotAgents {
		var eligibleBots []botProfile
		if len(g.botAgentTypes) == 0 {
			eligibleBots = allBotProfiles
		} else {
			for _, botName := range g.botAgentTypes {
				if profiles, ok := botProfileCategories[botName]; ok {
					eligibleBots = append(eligibleBots, profiles...)
				}
			}
		}

		if len(eligibleBots) == 0 {
			return nil, fmt.Errorf("legitagent: no bot profiles found for the specified types: %v", g.botAgentTypes)
		}

		chosenProfile := fastrand.Choice(eligibleBots)

		agent.UserAgent = chosenProfile.UserAgent
		agent.ClientHelloID = chosenProfile.HelloID

		for k, v := range chosenProfile.Headers {
			agent.Headers.Set(k, v)
		}

		keys := make([]string, 0, len(chosenProfile.Headers))
		for k := range chosenProfile.Headers {
			keys = append(keys, k)
		}
		PriorityHeaderSorter(keys)
		agent.HeaderOrder = append([]string{":method", ":authority", ":scheme", ":path"}, keys...)

		if g.h2Only {
			agent.H2Settings = GetChromiumH2Settings()
		} else {
			agent.H2Settings = nil
		}

		return agent, nil
	}

	browser, err := g.resolveBrowser()
	if err != nil {
		g.ReleaseAgent(agent)
		return nil, err
	}

	profile := browserProfiles[browser]

	chosenPlatform, chosenOS, err := g.resolvePlatformAndOS(browser)
	if err != nil {
		g.ReleaseAgent(agent)
		return nil, err
	}

	platformProf := platformProfiles[chosenPlatform]
	osProf := osProfiles[chosenOS]

	allVersions := getVersionKeys(profile.Versions)
	var possibleVersions []int

	if g.minVersion != 114 || g.maxVersion != 141 {
		possibleVersions = make([]int, 0, len(allVersions))
		for _, v := range allVersions {
			if v >= g.minVersion && v <= g.maxVersion {
				possibleVersions = append(possibleVersions, v)
			}
		}
	} else {
		possibleVersions = allVersions
	}

	var finalVersions []int
	if g.h2Only {
		finalVersions = make([]int, 0, len(possibleVersions))
		for _, v := range possibleVersions {
			if profile.Versions[v].SupportsH2 {
				finalVersions = append(finalVersions, v)
			}
		}
	} else {
		finalVersions = possibleVersions
	}

	if len(finalVersions) == 0 {
		return nil, fmt.Errorf("legitagent: no available browser versions for %s that meet the specified criteria", browser)
	}

	version := fastrand.Choice(finalVersions)
	versionProf := profile.Versions[version]

	fullVersion := ""
	if profile.ChromiumBased {
		fullVersion = fmt.Sprintf("%d.0.%d.%d", version, versionProf.BuildNumber, fastrand.IntN(999))
	}

	sb := builderPool.Get().(*strings.Builder)
	defer func() {
		sb.Reset()
		builderPool.Put(sb)
	}()

	componentGenerators := platformProf.ComponentGenerators[profile.Family]
	firstPart := true

	for _, componentGenerator := range componentGenerators {
		part := componentGenerator(profile, osProf, versionProf, fullVersion)
		if part != "" {
			if !firstPart {
				sb.WriteByte(' ')
			}
			sb.WriteString(part)
			firstPart = false
		}
	}

	agent.UserAgent = sb.String()

	headerSorter := g.headerSorter

	if g.fingerprintProfile == FingerprintProfileMaximum {
		headerSorter = ShuffledPriorityHeaderSorter
	}

	if !g.zeroHeader {
		agent.Headers, agent.HeaderOrder = g.buildHeaders(
			profile,
			osProf,
			platformProf,
			version,
			fullVersion,
			versionProf,
			headerSorter,
		)
	} else {
		agent.Headers = nil
		agent.HeaderOrder = nil
	}

	if g.h2Only {
		agent.H2Settings = profile.H2Settings()
		if g.h2RandomizationProfile != H2RandomizationProfileNone {
			agent.H2Settings = randomizeH2Settings(agent.H2Settings, g.h2RandomizationProfile)
		}
	} else {
		agent.H2Settings = nil
	}

	if g.fingerprintProfile == FingerprintProfileMaximum {
		agent.ClientHelloSpec = ChromeLatestSpec()
		agent.ClientHelloID = utls.ClientHelloID{}
	} else {
		agent.ClientHelloID = versionProf.TLS.HelloID
		agent.ClientHelloSpec = nil
	}

	return agent, nil
}

func (g *Generator) ReleaseAgent(a *Agent) {
	if a == nil {
		return
	}

	a.UserAgent = ""
	for k := range a.Headers {
		delete(a.Headers, k)
	}

	a.HeaderOrder = nil
	a.ClientHelloSpec = nil
	a.ClientHelloID = utls.ClientHelloID{}
	a.H2Settings = nil
	g.agentPool.Put(a)
}

func (g *Generator) resolveBrowser() (Browser, error) {
	var potentialBrowsers []Browser
	hasRandom := false

	for _, b := range g.browsers {
		if b == BrowserRandom {
			hasRandom = true
			break
		}
	}
	if hasRandom {
		potentialBrowsers = allRealBrowsers
	} else {
		potentialBrowsers = g.browsers
	}

	if len(potentialBrowsers) == 0 {
		return "", fmt.Errorf("legitagent: no browsers configured for generation")
	}

	return fastrand.Choice(potentialBrowsers), nil
}

func (g *Generator) resolvePlatformAndOS(browser Browser) (Platform, OperatingSystem, error) {
	type combo struct {
		platform Platform
		os       OperatingSystem
	}

	validCombos := make([]combo, 0, len(allRealPlatforms)*len(allRealOS))

	userPlatforms := g.platforms
	if len(userPlatforms) == 1 && userPlatforms[0] == PlatformRandom {
		userPlatforms = allRealPlatforms
	}

	userOSes := g.os
	if len(userOSes) == 1 && userOSes[0] == OSRandom {
		userOSes = allRealOS
	}

	for _, p := range userPlatforms {
		for _, o := range userOSes {
			var concreteOSes []OperatingSystem

			if o == OSMac {
				concreteOSes = append(concreteOSes, macArchitectures...)
			} else {
				concreteOSes = append(concreteOSes, o)
			}

			for _, concreteOS := range concreteOSes {
				osProfile, osProfileExists := osProfiles[concreteOS]
				if !osProfileExists {
					continue
				}

				if (p == PlatformMobile) != osProfile.IsMobile {
					continue
				}

				isValidForBrowser := false
				switch browser {
				case BrowserSafari:
					if (p == PlatformMobile && concreteOS == OSiOS) || (p == PlatformDesktop && osProfile.Name == "macOS") {
						isValidForBrowser = true
					}
				case BrowserFirefox:
					if (p == PlatformMobile && concreteOS == OSAndroid) || (p == PlatformDesktop && concreteOS != OSiOS && concreteOS != OSAndroid) {
						isValidForBrowser = true
					}
				default:
					isValidForBrowser = true
				}

				if isValidForBrowser {
					validCombos = append(validCombos, combo{p, concreteOS})
				}
			}
		}
	}

	if len(validCombos) == 0 {
		return "", "", fmt.Errorf("no compatible platform/OS combination found for browser %s with the current settings", browser)
	}

	chosenCombo := fastrand.Choice(validCombos)

	return chosenCombo.platform, chosenCombo.os, nil
}

func (g *Generator) buildHeaders(browser browserProfile, os osProfile, platform platformProfile, version int, fullVersion string, versionProf versionProfile, sorter HeaderSorter) (http.Header, []string) {
	headerMap := make(map[string]string, 16)

	var acceptTemplate [][]AcceptHeaderPart
	if g.requestType == RequestTypeXHR {
		acceptTemplate = versionProf.AcceptHeaderPatternsXHR
	} else {
		acceptTemplate = versionProf.AcceptHeaderPatterns
	}

	languageTemplate := fastrand.Choice(g.languageProfiles)

	if g.acceptEnabled {
		headerMap["accept"] = buildAcceptHeader(fastrand.Choice(acceptTemplate))
	}
	if g.acceptEncodingEnabled {
		headerMap["accept-encoding"] = generateAcceptEncoding()
	}

	headerMap["accept-language"] = buildAcceptHeader(languageTemplate)

	if browser.ChromiumBased {
		headerMap["sec-ch-ua"] = buildSecChUa(browser.Brand, strconv.Itoa(version), false, true)
		headerMap["sec-ch-ua-mobile"] = platform.MobileHint
		headerMap["sec-ch-ua-platform"] = fmt.Sprintf(`"%s"`, os.Name)

		if g.fullFingerprint {
			headerMap["sec-ch-ua-full-version-list"] = buildSecChUa(browser.Brand, fullVersion, true, true)
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
	}

	if g.fingerprintProfile == FingerprintProfileExtreme {
		for k := range headerMap {
			if strings.HasPrefix(k, "sec-") && fastrand.Bool() {
				delete(headerMap, k)
			}
		}
	}

	if g.requestType == RequestTypeNavigate && browser.Brand == "Brave" {
		headerMap["sec-gpc"] = "1"
	}

	switch g.requestType {
	case RequestTypeNavigate:
		headerMap["sec-fetch-dest"] = "document"
		headerMap["sec-fetch-mode"] = "navigate"
		headerMap["sec-fetch-site"] = "none"
		headerMap["sec-fetch-user"] = "?1"
		headerMap["upgrade-insecure-requests"] = "1"
	case RequestTypeSubresource:
		headerMap["sec-fetch-dest"] = fastrand.Choice(subresourceDests)
		headerMap["sec-fetch-mode"] = "no-cors"
		headerMap["sec-fetch-site"] = "same-origin"
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

	sorter(keys)
	orderedKeys := append([]string{":method", ":authority", ":scheme", ":path"}, keys...)

	for _, k := range keys {
		header.Set(k, headerMap[k])
	}

	return header, orderedKeys
}

func buildSecChUa(brand, version string, isFull, randomize bool) string {
	var greaseBrand string
	if randomize {
		greaseBrand = fastrand.Choice(greaseBrands)
	} else {
		greaseBrand = `"Not/A)Brand";v="8"`
	}

	parts := strings.SplitN(greaseBrand, `;v=`, 2)
	greaseKey := parts[0]
	var greaseVersion string
	if isFull {
		greaseVersion = `"99.0.0.0"`
	} else {
		greaseVersion = parts[1]
	}

	v := version
	if !isFull {
		if maj := strings.Split(version, "."); len(maj) > 0 {
			v = maj[0]
		}
	}

	brands := make([]string, 0, 3)
	brands = append(
		brands,
		fmt.Sprintf(`"Chromium";v="%s"`, v),
		fmt.Sprintf(`"%s";v="%s"`, brand, v),
		fmt.Sprintf(`%s;v=%s`, greaseKey, greaseVersion),
	)

	if randomize {
		fastrand.Shuffle(len(brands), func(i, j int) { brands[i], brands[j] = brands[j], brands[i] })
	}

	return strings.Join(brands, ", ")
}

func generateAcceptEncoding() string {
	encodings := []string{"gzip", "deflate", "br"}

	fastrand.Shuffle(len(encodings), func(i, j int) {
		encodings[i], encodings[j] = encodings[j], encodings[i]
	})

	if fastrand.IntN(2) == 1 {
		return strings.Join(encodings, ", ") + ", zstd"
	}

	return strings.Join(encodings, ", ")
}

func buildAcceptHeader(parts []AcceptHeaderPart) string {
	sb := builderPool.Get().(*strings.Builder)
	defer func() {
		sb.Reset()
		builderPool.Put(sb)
	}()

	currentQ := 1.0
	for i, part := range parts {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(part.Value)

		if len(part.Extras) > 0 {
			for _, extra := range part.Extras {
				sb.WriteByte(';')
				sb.WriteString(extra)
			}
		}

		if part.Q > 0 {
			currentQ -= fastrand.Float64() * 0.1
			if currentQ < 0.1 {
				currentQ = 0.1
			}
			sb.WriteString(";q=")
			sb.WriteString(strconv.FormatFloat(currentQ, 'f', 1, 64))
		}
	}
	return sb.String()
}

func getVersionKeys(m map[int]versionProfile) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
