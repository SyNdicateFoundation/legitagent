package legitagent

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"testing"
	"time"

	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
)

type CloudflareTestData struct {
	UserAgent   string
	HTTPVersion string
	JA3Hash     string
	HumanScore  string
}

func extractString(html []byte, pattern string, groupIndex int) string {
	re := regexp.MustCompile(pattern)
	matches := re.FindSubmatch(html)
	if len(matches) > groupIndex {
		return string(bytes.TrimSpace(matches[groupIndex]))
	}
	return ""
}

func parseCloudflareTestData(html []byte) *CloudflareTestData {
	data := &CloudflareTestData{}
	data.UserAgent = extractString(html, `Your\s+user-agent\s+is\s+<span[^>]*>(.*?)</span>`, 1)
	data.HTTPVersion = extractString(html, `You are using <span[^>]*>(.*?)</span>`, 1)
	data.HumanScore = extractString(html, `and\s+you\s+are\s+<span[^>]*>\s*(\d+)%\s*human\s*</span>`, 1)
	data.JA3Hash = extractString(html, `The\s+JA3\s+hash\s+is\s+<span[^>]*>\s+([a-f0-9]+)\s+</span>`, 1)
	return data
}

func performCloudflareTest(agent *Agent, addr *url.URL) (*CloudflareTestData, error) {
	if agent.ClientHelloSpec == nil && agent.ClientHelloID == (utls.ClientHelloID{}) {
		return nil, errors.New("agent has no TLS fingerprint, skipping test")
	}

	dialTLSContext := func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
		rawConn, err := net.DialTimeout(network, addr, 10*time.Second)
		if err != nil {
			return nil, fmt.Errorf("net.DialTimeout failed: %w", err)
		}

		uTLSConfig := &utls.Config{ServerName: cfg.ServerName, InsecureSkipVerify: true, NextProtos: []string{"h2", "http/1.1"}}
		var uconn *utls.UConn

		if agent.ClientHelloSpec != nil {
			uconn = utls.UClient(rawConn, uTLSConfig, utls.HelloCustom)
			if err := uconn.ApplyPreset(agent.ClientHelloSpec); err != nil {
				return nil, fmt.Errorf("uconn.ApplyPreset failed: %w", err)
			}
		} else {
			uconn = utls.UClient(rawConn, uTLSConfig, agent.ClientHelloID)
		}

		if err := uconn.HandshakeContext(ctx); err != nil {
			return nil, fmt.Errorf("uconn.HandshakeContext failed: %w", err)
		}

		return uconn, nil
	}

	h2Transport := &http2.Transport{
		DialTLSContext:  dialTLSContext,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: h2Transport, Timeout: 15 * time.Second}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, addr.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext failed: %w", err)
	}

	req.Header = agent.Headers
	req.Header.Set("User-Agent", agent.UserAgent)
	req.Header.Del("Accept-Encoding")
	req.Host = addr.Hostname()

	resp, err := client.Do(req)
	if err != nil {
		if ctxErr := context.Cause(req.Context()); ctxErr != nil && errors.Is(ctxErr, context.DeadlineExceeded) {
			return nil, fmt.Errorf("client.Do failed (context deadline exceeded): %w", err)
		}
		return nil, fmt.Errorf("client.Do failed: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll failed: %w", err)
	}

	return parseCloudflareTestData(body), nil
}

func TestCloudflareStealth(t *testing.T) {
	parse, err := url.ParseRequestURI("https://cloudflare.manfredi.io/test/")
	if err != nil {
		t.Fatal(err)
	}
	g := NewGenerator(WithBrowsers(BrowserRandom), WithOS(OSRandom), WithFullFingerprint(true))
	for i := 0; i < 5; i++ {
		t.Run(fmt.Sprintf("Agent_%d", i+1), func(t *testing.T) {
			t.Parallel()
			agent, err := g.Generate()
			if err != nil {
				t.Fatalf("Failed to generate agent: %s", err)
			}
			defer g.ReleaseAgent(agent)

			result, err := performCloudflareTest(agent, parse)
			if err != nil {
				t.Fatalf("Failed to test agent %s: %s", agent.UserAgent, err)
			}

			t.Logf("%s: HumanScore: %s", agent.UserAgent, result.HumanScore)

			if result.JA3Hash == "" {
				t.Error("JA3 hash was empty")
			}
		})
	}
}
