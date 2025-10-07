package legitagent

import (
	"strings"
	"testing"
)

func TestBrowserSpecificGeneration(t *testing.T) {
	t.Run("Chrome", func(t *testing.T) {
		g := NewGenerator(
			WithBrowsers(BrowserChrome),
			WithVersionRange(140, 140),
			WithOS(OSWindows11),
			WithPlatforms(PlatformDesktop),
		)

		agent, err := g.Generate()
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}
		if !strings.Contains(agent.UserAgent, "Chrome/140.0.7255") {
			t.Errorf("Expected Chrome 140 UA, got: %s", agent.UserAgent)
		}
	})

	t.Run("Firefox", func(t *testing.T) {
		g := NewGenerator(
			WithBrowsers(BrowserFirefox),
			WithVersionRange(128, 128),
			WithOS(OSLinux),
			WithPlatforms(PlatformDesktop),
		)

		agent, err := g.Generate()
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}

		if !strings.Contains(agent.UserAgent, "Firefox/128.0") || !strings.Contains(agent.UserAgent, "Gecko/") {
			t.Errorf("Expected Firefox 128 UA, got: %s", agent.UserAgent)
		}
	})

	t.Run("FirefoxMobile", func(t *testing.T) {
		g := NewGenerator(
			WithBrowsers(BrowserFirefox),
			WithVersionRange(127, 127),
			WithPlatforms(PlatformMobile),
		)

		agent, err := g.Generate()
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}
		if !strings.Contains(agent.UserAgent, "Firefox/127.0") || !strings.Contains(agent.UserAgent, "Android") {
			t.Errorf("Expected Firefox Mobile UA on Android, got: %s", agent.UserAgent)
		}
		if strings.Contains(agent.UserAgent, "{device_model}") {
			t.Errorf("Device model placeholder was not replaced in Firefox Mobile UA: %s", agent.UserAgent)
		}
	})

	t.Run("SafariMobile", func(t *testing.T) {
		g := NewGenerator(
			WithBrowsers(BrowserSafari),
			WithPlatforms(PlatformMobile),
		)

		agent, err := g.Generate()
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}

		if !strings.Contains(agent.UserAgent, "iPhone") || !strings.Contains(agent.UserAgent, "Mobile/") {
			t.Errorf("Expected Safari Mobile UA on iPhone, got: %s", agent.UserAgent)
		}
	})

	t.Run("Opera", func(t *testing.T) {
		g := NewGenerator(
			WithBrowsers(BrowserOpera),
			WithVersionRange(128, 128),
			WithOS(OSMac),
			WithPlatforms(PlatformDesktop),
		)

		agent, err := g.Generate()
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}

		if !strings.Contains(agent.UserAgent, "OPR/128") {
			t.Errorf("Expected Opera 128 UA, got: %s", agent.UserAgent)
		}
	})
}

func TestFullFingerprintOption(t *testing.T) {
	t.Run("Minimal (Default)", func(t *testing.T) {
		g := NewGenerator(WithBrowsers(BrowserChrome))
		agent, err := g.Generate()
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}

		if agent.Headers.Get("sec-ch-ua-arch") != "" {
			t.Error("Minimal fingerprint should not have sec-ch-ua-arch")
		}
	})

	t.Run("Full", func(t *testing.T) {
		g := NewGenerator(
			WithBrowsers(BrowserChrome),
			WithFullFingerprint(true),
			WithOS(OSWindows11),
		)

		agent, err := g.Generate()
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}

		if agent.Headers.Get("sec-ch-ua-arch") == "" {
			t.Error("Full fingerprint should have sec-ch-ua-arch")
		}
	})
}

func TestRandomOptions(t *testing.T) {
	t.Run("OSRandom", func(t *testing.T) {
		g := NewGenerator(
			WithOS(OSRandom),
			WithBrowsers(BrowserChrome),
		)

		platforms := make(map[string]bool)
		for i := 0; i < 50; i++ {
			agent, err := g.Generate()
			if err != nil {
				t.Fatalf("Generate failed: %v", err)
			}
			platforms[agent.Headers.Get("sec-ch-ua-platform")] = true
		}

		if len(platforms) < 3 {
			t.Errorf("OSRandom did not produce a variety of OS platforms, got: %v", platforms)
		}
	})

	t.Run("BrowserRandom", func(t *testing.T) {
		g := NewGenerator(WithBrowsers(BrowserRandom))

		browsers := make(map[string]bool)

		for i := 0; i < 50; i++ {
			agent, err := g.Generate()
			if err != nil {
				t.Fatalf("Generate failed on iteration %d: %v", i, err)
			}

			ua := agent.UserAgent
			switch {
			case strings.Contains(ua, "Firefox/"):
				browsers["firefox"] = true
			case strings.Contains(ua, "OPR/"):
				browsers["opera"] = true
			case strings.Contains(ua, "Edg/"):
				browsers["edge"] = true
			case strings.Contains(ua, "Chrome/"):
				browsers["chrome"] = true
			case strings.Contains(ua, "Safari/") && !strings.Contains(ua, "Chrome/"):
				browsers["safari"] = true
			}
		}

		if len(browsers) < 4 {
			t.Errorf("BrowserRandom did not produce a variety of browser agents, got: %v", browsers)
		}
	})

	t.Run("PlatformRandom", func(t *testing.T) {
		g := NewGenerator(
			WithPlatforms(PlatformRandom),
			WithBrowsers(BrowserChrome),
		)

		platforms := make(map[string]bool)

		for i := 0; i < 50; i++ {
			agent, err := g.Generate()
			if err != nil {
				t.Fatalf("Generate failed: %v", err)
			}
			platforms[agent.Headers.Get("sec-ch-ua-mobile")] = true
		}

		if !platforms["?0"] || !platforms["?1"] {
			t.Errorf("PlatformRandom did not produce both Desktop ('?0') and Mobile ('?1') agents, got: %v", platforms)
		}
	})
}

func TestAcceptEncodingOption(t *testing.T) {
	t.Run("Disabled (Default)", func(t *testing.T) {
		g := NewGenerator()
		agent, err := g.Generate()
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}
		defer g.ReleaseAgent(agent)

		if val := agent.Headers.Get("accept-encoding"); val != "" {
			t.Errorf("Expected 'accept-encoding' header to be empty, but got '%s'", val)
		}
	})

	t.Run("Enabled", func(t *testing.T) {
		g := NewGenerator(WithAcceptEncoding(true))
		agent, err := g.Generate()
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}
		defer g.ReleaseAgent(agent)

		val := agent.Headers.Get("accept-encoding")
		if val == "" {
			t.Error("Expected 'accept-encoding' header to be present, but it was empty")
		}
		if !strings.Contains(val, "gzip") {
			t.Errorf("Expected 'accept-encoding' to contain 'gzip', but got '%s'", val)
		}
	})

	t.Run("Bot Agent (Unaffected)", func(t *testing.T) {
		g := NewGenerator(WithBotAgents(BotGoogle), WithAcceptEncoding(false))
		agent, err := g.Generate()
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}
		defer g.ReleaseAgent(agent)

		val := agent.Headers.Get("accept-encoding")
		if val == "" {
			t.Error("Expected bot agent 'accept-encoding' header to be present regardless of the option, but it was removed")
		}
	})
}

func TestAcceptHeaderOption(t *testing.T) {
	t.Run("Enabled (Default)", func(t *testing.T) {
		g := NewGenerator()
		agent, err := g.Generate()
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}
		defer g.ReleaseAgent(agent)

		if val := agent.Headers.Get("accept"); val == "" {
			t.Errorf("Expected 'accept' header to be present, but it was empty")
		}
	})

	t.Run("Disabled", func(t *testing.T) {
		g := NewGenerator(WithAccept(false))
		agent, err := g.Generate()
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}
		defer g.ReleaseAgent(agent)

		if val := agent.Headers.Get("accept"); val != "" {
			t.Errorf("Expected 'accept' header to be empty, but got '%s'", val)
		}
	})

	t.Run("Bot Agent (Unaffected)", func(t *testing.T) {
		g := NewGenerator(WithBotAgents(BotGoogle), WithAccept(false))
		agent, err := g.Generate()
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}
		defer g.ReleaseAgent(agent)

		val := agent.Headers.Get("accept")
		if val == "" {
			t.Error("Expected bot agent 'accept' header to be present regardless of the option, but it was removed")
		}
	})
}
