package legitagent

import (
	"reflect"
	"testing"

	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
)

func TestFingerprintProfileMaximum(t *testing.T) {
	g := NewGenerator(WithFingerprintProfile(FingerprintProfileMaximum))

	agent1, err := g.Generate()
	if err != nil {
		t.Fatalf("Failed to generate first agent: %v", err)
	}

	defer g.ReleaseAgent(agent1)

	agent2, err := g.Generate()
	if err != nil {
		t.Fatalf("Failed to generate second agent: %v", err)
	}

	defer g.ReleaseAgent(agent2)

	if agent1.ClientHelloSpec == nil || agent1.ClientHelloID != (utls.ClientHelloID{}) {
		t.Error("Expected a dynamic ClientHelloSpec, but got a static ClientHelloID")
	}

	if agent2.ClientHelloSpec == nil {
		t.Error("Expected the second agent to also have a dynamic ClientHelloSpec")
	}

	if reflect.DeepEqual(agent1.HeaderOrder, agent2.HeaderOrder) {
		t.Logf("Warning: Header orders were identical. This is possible but less likely with shuffling.")
		t.Logf("Agent 1: %v", agent1.HeaderOrder)
		t.Logf("Agent 2: %v", agent2.HeaderOrder)
	}

	t.Log("FingerprintProfileMaximum correctly generated dynamic JA3 and shuffled headers.")
}

func TestFingerprintProfileNormal(t *testing.T) {
	g := NewGenerator(WithFingerprintProfile(FingerprintProfileNormal))

	agent, err := g.Generate()
	if err != nil {
		t.Fatalf("Failed to generate agent: %v", err)
	}

	defer g.ReleaseAgent(agent)

	if agent.ClientHelloSpec != nil || agent.ClientHelloID == (utls.ClientHelloID{}) {
		t.Errorf("Expected a static ClientHelloID, but got a dynamic ClientHelloSpec or nil ID. Spec: %v, ID: %v", agent.ClientHelloSpec, agent.ClientHelloID)
	}

	t.Log("FingerprintProfileNormal correctly generated a static JA3 profile.")
}

func TestBrowserSpecificH2Settings(t *testing.T) {
	testCases := map[Browser]map[http2.SettingID]uint32{
		BrowserChrome:  {http2.SettingHeaderTableSize: 65536},
		BrowserFirefox: {http2.SettingInitialWindowSize: 131072},
		BrowserSafari:  {http2.SettingMaxHeaderListSize: 16384},
	}

	for browser, expectedSettings := range testCases {
		t.Run(string(browser), func(t *testing.T) {
			g := NewGenerator(WithBrowsers(browser))

			agent, err := g.Generate()
			if err != nil {
				t.Fatalf("Failed to generate agent for %s: %v", browser, err)
			}

			defer g.ReleaseAgent(agent)

			for settingID, expectedValue := range expectedSettings {
				actualValue, ok := agent.H2Settings[settingID]
				if !ok {
					t.Errorf("Expected H2 setting %s not found", settingID)
				}
				if actualValue != expectedValue {
					t.Errorf("Incorrect H2 setting for %s. Expected %s=%d, got %d", browser, settingID, expectedValue, actualValue)
				}
			}
		})
	}
}
