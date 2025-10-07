package legitagent

import (
	"strconv"
	"strings"
)

type Browser string

const (
	BrowserRandom  Browser = "random"
	BrowserChrome  Browser = "chrome"
	BrowserOpera   Browser = "opera"
	BrowserEdge    Browser = "edge"
	BrowserBrave   Browser = "brave"
	BrowserFirefox Browser = "firefox"
	BrowserSafari  Browser = "safari"
)

type Platform string

const (
	PlatformRandom  Platform = "random"
	PlatformDesktop Platform = "desktop"
	PlatformMobile  Platform = "mobile"
)

type OperatingSystem string

const (
	OSRandom          OperatingSystem = "random"
	OSWindows         OperatingSystem = "windows"
	OSWindows11       OperatingSystem = "windows11"
	OSLinux           OperatingSystem = "linux"
	OSMac             OperatingSystem = "mac"
	OSAndroid         OperatingSystem = "android"
	OSiOS             OperatingSystem = "ios"
	OSChromeOS        OperatingSystem = "chromeos"
	osMacIntel        OperatingSystem = "mac_intel"
	osMacAppleSilicon OperatingSystem = "mac_apple_silicon"
	osUbuntu          OperatingSystem = "ubuntu"
	osFedora          OperatingSystem = "fedora"
)

type RequestType string

const (
	RequestTypeNavigate    RequestType = "navigate"
	RequestTypeSubresource RequestType = "subresource"
	RequestTypeXHR         RequestType = "xhr"
)

type FingerprintProfile int

const (
	FingerprintProfileNormal FingerprintProfile = iota

	FingerprintProfileMaximum
)

type Option func(*Generator)

func WithBrowsers(b ...Browser) Option {
	return func(g *Generator) {
		if len(b) > 0 {
			g.browsers = b
		}
	}
}

func WithPlatforms(p ...Platform) Option {
	return func(g *Generator) {
		if len(p) > 0 {
			g.platforms = p
		}
	}
}

func WithOS(os ...OperatingSystem) Option {
	return func(g *Generator) {
		if len(os) > 0 {
			g.os = os
		}
	}
}

func WithVersionRange(min, max int) Option {
	return func(g *Generator) {
		if min > 0 {
			g.minVersion = min
		}
		if max > 0 && max >= min {
			g.maxVersion = max
		}
	}
}

func parseLanguageHeader(header string) []AcceptHeaderPart {
	partsStr := strings.Split(header, ",")
	parts := make([]AcceptHeaderPart, 0, len(partsStr))
	for _, partStr := range partsStr {
		partStr = strings.TrimSpace(partStr)
		if strings.Contains(partStr, ";q=") {
			split := strings.SplitN(partStr, ";q=", 2)
			q, err := strconv.ParseFloat(split[1], 64)
			if err == nil {
				parts = append(parts, AcceptHeaderPart{Value: split[0], Q: q})
			}
		} else {
			parts = append(parts, AcceptHeaderPart{Value: partStr})
		}
	}
	return parts
}

func WithLanguages(langs ...string) Option {
	return func(g *Generator) {
		if len(langs) > 0 {
			profiles := make([][]AcceptHeaderPart, 0, len(langs))
			for _, langStr := range langs {
				profiles = append(profiles, parseLanguageHeader(langStr))
			}
			g.languageProfiles = profiles
		}
	}
}

func WithRequestType(rt RequestType) Option {
	return func(g *Generator) { g.requestType = rt }
}

func WithHeaderSorter(sorter HeaderSorter) Option {
	return func(g *Generator) {
		if sorter != nil {
			g.headerSorter = sorter
		}
	}
}

func WithFullFingerprint(full bool) Option {
	return func(g *Generator) {
		g.fullFingerprint = full
	}
}

func WithH2Only(h2 bool) Option {
	return func(g *Generator) {
		g.h2Only = h2
	}
}

func WithFingerprintProfile(profile FingerprintProfile) Option {
	return func(g *Generator) {
		g.fingerprintProfile = profile
	}
}
