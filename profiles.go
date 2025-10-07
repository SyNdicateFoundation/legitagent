package legitagent

import (
	"fmt"
	"strings"

	"github.com/SyNdicateFoundation/fastrand"
	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
)

type AcceptHeaderPart struct {
	Value  string
	Q      float64
	Extras []string
}

type UAComponentGenerator func(browserProfile browserProfile, osProfile osProfile, versionProfile versionProfile, fullVersion string) string
type BrowserFamily string

const (
	Chromium BrowserFamily = "Chromium"
	Gecko    BrowserFamily = "Gecko"
	WebKit   BrowserFamily = "WebKit"
)

type tlsProfile struct {
	HelloID    utls.ClientHelloID
	ClientSpec func() *utls.ClientHelloSpec
}

type versionProfile struct {
	BuildNumber             int
	AcceptHeaderPatterns    [][]AcceptHeaderPart
	AcceptHeaderPatternsXHR [][]AcceptHeaderPart
	TLS                     tlsProfile
	GeckoRevision           string
	WebKitVersion           string
	MobileVersion           string
	SafariVersion           string
	SupportsH2              bool
}

type browserProfile struct {
	Brand         string
	Family        BrowserFamily
	UASuffix      string
	Versions      map[int]versionProfile
	ChromiumBased bool
	H2Settings    func() map[http2.SettingID]uint32
}

type osProfile struct {
	Name          string
	PlatformToken string
	Version       string
	Arch          string
	BitnessHint   string
	IsMobile      bool
}

type platformProfile struct {
	MobileHint          string
	ComponentGenerators map[BrowserFamily][]UAComponentGenerator
}

var (
	acceptHeaderPatternsChrome = [][]AcceptHeaderPart{
		{
			{Value: "text/html"},
			{Value: "application/xhtml+xml"},
			{Value: "application/xml", Q: 0.9},
			{Value: "image/avif"},
			{Value: "image/webp"},
			{Value: "image/apng"},
			{Value: "*/*", Q: 0.8},
			{Value: "application/signed-exchange", Q: 0.7, Extras: []string{"v=b3"}},
		},
		{
			{Value: "text/html"},
			{Value: "application/xhtml+xml"},
			{Value: "application/xml", Q: 0.9},
			{Value: "image/avif"},
			{Value: "image/webp"},
			{Value: "image/apng"},
			{Value: "*/*", Q: 0.8},
		},
	}
	acceptHeaderPatternsFirefox = [][]AcceptHeaderPart{
		{
			{Value: "text/html"},
			{Value: "application/xhtml+xml"},
			{Value: "application/xml", Q: 0.9},
			{Value: "image/avif"},
			{Value: "image/webp"},
			{Value: "*/*", Q: 0.8},
		},
	}
	acceptHeaderPatternsSafari = [][]AcceptHeaderPart{
		{
			{Value: "text/html"},
			{Value: "application/xhtml+xml"},
			{Value: "application/xml", Q: 0.9},
			{Value: "*/*", Q: 0.8},
		},
	}
	acceptHeaderPatternsXHR = [][]AcceptHeaderPart{
		{{Value: "*/*"}},
	}

	tlsProfileChrome120  = tlsProfile{HelloID: utls.HelloChrome_120}
	tlsProfileFirefox120 = tlsProfile{HelloID: utls.HelloFirefox_120}
	tlsProfileSafari16   = tlsProfile{HelloID: utls.HelloSafari_16_0}

	chromeVersions = map[int]versionProfile{
		114: {BuildNumber: 5735, AcceptHeaderPatterns: acceptHeaderPatternsChrome, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileChrome120, SupportsH2: true},
		116: {BuildNumber: 5845, AcceptHeaderPatterns: acceptHeaderPatternsChrome, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileChrome120, SupportsH2: true},
		118: {BuildNumber: 5993, AcceptHeaderPatterns: acceptHeaderPatternsChrome, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileChrome120, SupportsH2: true},
		120: {BuildNumber: 6099, AcceptHeaderPatterns: acceptHeaderPatternsChrome, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileChrome120, SupportsH2: true},
		124: {BuildNumber: 6367, AcceptHeaderPatterns: acceptHeaderPatternsChrome, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileChrome120, SupportsH2: true},
		128: {BuildNumber: 6636, AcceptHeaderPatterns: acceptHeaderPatternsChrome, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileChrome120, SupportsH2: true},
		130: {BuildNumber: 6735, AcceptHeaderPatterns: acceptHeaderPatternsChrome, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileChrome120, SupportsH2: true},
		133: {BuildNumber: 6912, AcceptHeaderPatterns: acceptHeaderPatternsChrome, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileChrome120, SupportsH2: true},
		140: {BuildNumber: 7255, AcceptHeaderPatterns: acceptHeaderPatternsChrome, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileChrome120, SupportsH2: true},
	}
	edgeVersions = map[int]versionProfile{
		114: {BuildNumber: 1823, AcceptHeaderPatterns: acceptHeaderPatternsChrome, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileChrome120, SupportsH2: true},
		116: {BuildNumber: 1938, AcceptHeaderPatterns: acceptHeaderPatternsChrome, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileChrome120, SupportsH2: true},
		118: {BuildNumber: 2088, AcceptHeaderPatterns: acceptHeaderPatternsChrome, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileChrome120, SupportsH2: true},
		120: {BuildNumber: 2210, AcceptHeaderPatterns: acceptHeaderPatternsChrome, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileChrome120, SupportsH2: true},
		124: {BuildNumber: 2478, AcceptHeaderPatterns: acceptHeaderPatternsChrome, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileChrome120, SupportsH2: true},
		128: {BuildNumber: 2739, AcceptHeaderPatterns: acceptHeaderPatternsChrome, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileChrome120, SupportsH2: true},
		133: {BuildNumber: 2988, AcceptHeaderPatterns: acceptHeaderPatternsChrome, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileChrome120, SupportsH2: true},
		140: {BuildNumber: 3265, AcceptHeaderPatterns: acceptHeaderPatternsChrome, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileChrome120, SupportsH2: true},
	}
	braveVersions = chromeVersions

	browserProfiles = map[Browser]browserProfile{
		BrowserChrome: {Brand: "Google Chrome", Family: Chromium, UASuffix: "", ChromiumBased: true, Versions: chromeVersions, H2Settings: GetChromiumH2Settings},
		BrowserOpera:  {Brand: "Opera", Family: Chromium, UASuffix: "OPR/%s", ChromiumBased: true, Versions: chromeVersions, H2Settings: GetChromiumH2Settings},
		BrowserEdge:   {Brand: "Microsoft Edge", Family: Chromium, UASuffix: "Edg/%s", ChromiumBased: true, Versions: edgeVersions, H2Settings: GetChromiumH2Settings},
		BrowserBrave:  {Brand: "Brave", Family: Chromium, UASuffix: "", ChromiumBased: true, Versions: braveVersions, H2Settings: GetChromiumH2Settings},
		BrowserFirefox: {Brand: "Firefox", Family: Gecko, ChromiumBased: false, Versions: map[int]versionProfile{
			115: {GeckoRevision: "115.0", AcceptHeaderPatterns: acceptHeaderPatternsFirefox, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileFirefox120, SupportsH2: true},
			120: {GeckoRevision: "120.0", AcceptHeaderPatterns: acceptHeaderPatternsFirefox, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileFirefox120, SupportsH2: true},
			127: {GeckoRevision: "127.0", AcceptHeaderPatterns: acceptHeaderPatternsFirefox, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileFirefox120, SupportsH2: true},
			128: {GeckoRevision: "128.0", AcceptHeaderPatterns: acceptHeaderPatternsFirefox, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileFirefox120, SupportsH2: true},
		}, H2Settings: GetGeckoH2Settings},
		BrowserSafari: {Brand: "Safari", Family: WebKit, ChromiumBased: false, Versions: map[int]versionProfile{
			16: {WebKitVersion: "605.1.15", MobileVersion: "20F66", SafariVersion: "16.5", AcceptHeaderPatterns: acceptHeaderPatternsSafari, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileSafari16, SupportsH2: true},
			17: {WebKitVersion: "605.1.15", MobileVersion: "15E148", SafariVersion: "17.5", AcceptHeaderPatterns: acceptHeaderPatternsSafari, AcceptHeaderPatternsXHR: acceptHeaderPatternsXHR, TLS: tlsProfileSafari16, SupportsH2: true},
		}, H2Settings: GetWebKitH2Settings},
	}

	osProfiles = map[OperatingSystem]osProfile{
		OSWindows:         {Name: "Windows", PlatformToken: "Windows NT 10.0; Win64; x64", Version: "10.0.0", Arch: "x86", BitnessHint: "64", IsMobile: false},
		OSWindows11:       {Name: "Windows", PlatformToken: "Windows NT 10.0; Win64; x64", Version: "15.0.0", Arch: "x86", BitnessHint: "64", IsMobile: false},
		osMacIntel:        {Name: "macOS", PlatformToken: "Macintosh; Intel Mac OS X 10_15_7", Version: "14.5.0", Arch: "x86", BitnessHint: "64", IsMobile: false},
		osMacAppleSilicon: {Name: "macOS", PlatformToken: "Macintosh; ARM Mac OS X 10_15_7", Version: "14.5.0", Arch: "arm", BitnessHint: "64", IsMobile: false},
		OSLinux:           {Name: "Linux", PlatformToken: "X11; Linux x86_64", Version: "", Arch: "x86", BitnessHint: "64", IsMobile: false},
		osUbuntu:          {Name: "Linux", PlatformToken: "X11; Ubuntu; Linux x86_64", Version: "", Arch: "x86", BitnessHint: "64", IsMobile: false},
		osFedora:          {Name: "Linux", PlatformToken: "X11; Fedora; Linux x86_64", Version: "", Arch: "x86", BitnessHint: "64", IsMobile: false},
		OSAndroid:         {Name: "Android", PlatformToken: "Linux; Android 14; {device_model}", Version: "14.0.0", Arch: "arm", BitnessHint: "64", IsMobile: true},
		OSiOS:             {Name: "iOS", PlatformToken: "iPhone; CPU iPhone OS 17_5_1 like Mac OS X", Version: "17.5.1", IsMobile: true},
		OSChromeOS:        {Name: "Chrome OS", PlatformToken: "X11; CrOS x86_64 14541.0.0", Version: "14541.0.0", Arch: "x86", BitnessHint: "64", IsMobile: false},
	}

	platformProfiles = map[Platform]platformProfile{
		PlatformDesktop: {MobileHint: "?0", ComponentGenerators: map[BrowserFamily][]UAComponentGenerator{
			Chromium: {MozillaGenerator, OSGenerator, WebKitGenerator, KHTMLGenerator, ChromeGenerator, SafariGenerator, BrowserSuffixGenerator},
			Gecko:    {MozillaGenerator, FirefoxOSGenerator, GeckoTrailGenerator, FirefoxVersionGenerator},
			WebKit:   {MozillaGenerator, OSGenerator, SafariWebKitGenerator, KHTMLGenerator, SafariVersionGenerator, SafariBrowserVersionGenerator},
		}},
		PlatformMobile: {MobileHint: "?1", ComponentGenerators: map[BrowserFamily][]UAComponentGenerator{
			Chromium: {MozillaGenerator, OSGenerator, WebKitGenerator, KHTMLGenerator, ChromeGenerator, MobileSafariGenerator, BrowserSuffixGenerator},
			Gecko:    {MozillaGenerator, FirefoxOSGenerator, GeckoTrailGenerator, FirefoxVersionGenerator},
			WebKit:   {MozillaGenerator, OSGenerator, SafariWebKitGenerator, KHTMLGenerator, SafariVersionGenerator, SafariMobileTokenGenerator, SafariBrowserVersionGenerator},
		}},
	}
)

var greaseBrands = []string{`"Not/A)Brand";v="8"`, `"Not;A Brand";v="99"`, `"Not(A:Brand";v="24"`, `"Chromium";v="99"`}
var androidDevices = []string{"Pixel 7", "Pixel 8 Pro", "SM-S928B", "SM-G991U", "SM-F936U", "2201116SG", "V2109", "SM-A525F", "Pixel 6a", "SM-A536U", "Galaxy S23 Ultra"}
var subresourceDests = []string{"style", "script", "image", "font", "empty"}

func MozillaGenerator(_ browserProfile, _ osProfile, _ versionProfile, _ string) string {
	return "Mozilla/5.0"
}
func KHTMLGenerator(_ browserProfile, _ osProfile, _ versionProfile, _ string) string {
	return "(KHTML, like Gecko)"
}
func ChromeGenerator(_ browserProfile, _ osProfile, _ versionProfile, fv string) string {
	return "Chrome/" + fv
}
func SafariGenerator(_ browserProfile, _ osProfile, _ versionProfile, _ string) string {
	return "Safari/537.36"
}
func MobileSafariGenerator(_ browserProfile, _ osProfile, _ versionProfile, _ string) string {
	return "Mobile Safari/537.36"
}
func GeckoTrailGenerator(_ browserProfile, _ osProfile, _ versionProfile, _ string) string {
	return "Gecko/20100101"
}

func OSGenerator(_ browserProfile, op osProfile, _ versionProfile, _ string) string {
	token := op.PlatformToken
	if op.Name == "Android" {
		device := fastrand.Choice(androidDevices)
		token = strings.Replace(token, "{device_model}", device, 1)
	}
	return fmt.Sprintf("(%s)", token)
}

func FirefoxOSGenerator(_ browserProfile, op osProfile, vp versionProfile, _ string) string {
	token := op.PlatformToken
	if op.Name == "Android" {
		device := fastrand.Choice(androidDevices)
		token = strings.Replace(token, "{device_model}", device, 1)
	}
	return fmt.Sprintf("(%s; rv:%s)", token, vp.GeckoRevision)
}

func FirefoxVersionGenerator(_ browserProfile, _ osProfile, vp versionProfile, _ string) string {
	return fmt.Sprintf("Firefox/%s", vp.GeckoRevision)
}

func WebKitGenerator(_ browserProfile, _ osProfile, _ versionProfile, _ string) string {
	return "AppleWebKit/537.36"
}

func SafariWebKitGenerator(_ browserProfile, _ osProfile, vp versionProfile, _ string) string {
	return fmt.Sprintf("AppleWebKit/%s", vp.WebKitVersion)
}

func SafariVersionGenerator(_ browserProfile, _ osProfile, vp versionProfile, _ string) string {
	return fmt.Sprintf("Version/%s", vp.SafariVersion)
}

func SafariMobileTokenGenerator(_ browserProfile, _ osProfile, vp versionProfile, _ string) string {
	return fmt.Sprintf("Mobile/%s", vp.MobileVersion)
}

func SafariBrowserVersionGenerator(_ browserProfile, op osProfile, vp versionProfile, _ string) string {
	if op.IsMobile {
		return "Safari/604.1"
	}
	return fmt.Sprintf("Safari/%s", vp.WebKitVersion)
}

func BrowserSuffixGenerator(bp browserProfile, _ osProfile, _ versionProfile, fv string) string {
	if bp.UASuffix == "" {
		return ""
	}
	majorVersion := strings.Split(fv, ".")[0]
	return fmt.Sprintf(bp.UASuffix, majorVersion)
}
