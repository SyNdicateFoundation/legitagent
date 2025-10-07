package legitagent

import (
	"strings"
	"testing"

	utls "github.com/refraction-networking/utls"
)

func TestBotAgentGenerationDefault(t *testing.T) {
	g := NewGenerator(WithBotAgents())

	agent, err := g.Generate()
	if err != nil {
		t.Fatalf("Failed to generate bot agent: %v", err)
	}
	defer g.ReleaseAgent(agent)

	if agent.UserAgent == "" {
		t.Fatal("Generated bot agent has an empty User-Agent string")
	}

	if agent.Headers.Get("accept") == "" {
		t.Error("Generated bot agent is missing an 'accept' header")
	}

	isKnownBot := false
	for _, bot := range allBotProfiles {
		if agent.UserAgent == bot.UserAgent {
			isKnownBot = true
			break
		}
	}
	if !isKnownBot {
		t.Errorf("Generated agent's User-Agent is not in the full bot list: %s", agent.UserAgent)
	}
}

func TestBotAgentGenerationSpecific(t *testing.T) {
	g := NewGenerator(WithBotAgents(BotGoogle))

	for i := 0; i < 10; i++ {
		agent, err := g.Generate()
		if err != nil {
			t.Fatalf("Failed to generate bot agent: %v", err)
		}
		defer g.ReleaseAgent(agent)

		if !strings.Contains(agent.UserAgent, "Google") {
			t.Errorf("Expected a GoogleBot agent, but got: %s", agent.UserAgent)
		}

		if agent.ClientHelloID != utls.HelloGolang && agent.ClientHelloID != utls.HelloChrome_120 {
			t.Errorf("GoogleBot agent has unexpected HelloID: %v", agent.ClientHelloID)
		}
	}
}

func TestBotAgentGenerationMultipleSpecific(t *testing.T) {
	g := NewGenerator(WithBotAgents(BotDuckDuckGo, BotBaidu))
	foundDDG := false
	foundBaidu := false

	for i := 0; i < 50; i++ {
		agent, err := g.Generate()
		if err != nil {
			t.Fatalf("Failed to generate bot agent: %v", err)
		}

		if strings.Contains(agent.UserAgent, "DuckDuckBot") {
			foundDDG = true
		}
		if strings.Contains(agent.UserAgent, "Baiduspider") {
			foundBaidu = true
		}
		if !strings.Contains(agent.UserAgent, "DuckDuckBot") && !strings.Contains(agent.UserAgent, "Baiduspider") {
			t.Fatalf("Generated agent that was not in the specified list: %s", agent.UserAgent)
		}

		g.ReleaseAgent(agent)
	}

	if !foundDDG || !foundBaidu {
		t.Errorf("Did not generate agents from all specified categories. Found DDG: %v, Found Baidu: %v", foundDDG, foundBaidu)
	}
}

func TestBotAgentGenerationInvalid(t *testing.T) {
	g := NewGenerator(WithBotAgents("NonExistentBot"))
	_, err := g.Generate()
	if err == nil {
		t.Fatal("Expected an error for invalid bot type, but got nil")
	}
}
