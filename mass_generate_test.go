package legitagent

import (
	"testing"
	"time"
)

func TestMassGenerateAndPrint(t *testing.T) {
	g := NewGenerator(WithOS(OSRandom))

	start := time.Now()

	for i := 0; i < 1000; i++ {
		agent, err := g.Generate()
		if err != nil {
			t.Fatalf("Generate failed on iteration %d: %v", i+1, err)
		}
		if i%100 == 0 {
			t.Logf("Generated %d: %s", i+1, agent.UserAgent)
		}
		g.ReleaseAgent(agent)
	}

	elapsed := time.Since(start)
	t.Logf("Generated 1000 agents in %v", elapsed)
}

func TestGenerateAndPrint(t *testing.T) {
	g := NewGenerator(WithOS(OSRandom))
	start := time.Now()

	for i := 0; i < 10; i++ {
		agent, err := g.Generate()
		if err != nil {
			t.Fatalf("Generate failed on iteration %d: %v", i+1, err)
		}
		t.Logf("Generated %d: %s, %#v", i+1, agent.UserAgent, agent.Headers)
		g.ReleaseAgent(agent)
	}

	elapsed := time.Since(start)
	t.Logf("Generated 10 agents in %v", elapsed)
}
