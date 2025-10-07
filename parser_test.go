package legitagent

import (
	"errors"
	"strings"
	"testing"

	utls "github.com/refraction-networking/utls"
)

func TestFromUserAgentString(t *testing.T) {
	t.Run("Successful Chrome Parse", func(t *testing.T) {
		ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36"
		agent, err := FromUserAgentString(ua, RequestTypeNavigate)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if agent.UserAgent != ua {
			t.Errorf("Expected UserAgent to be identical, got %s", agent.UserAgent)
		}
		if agent.ClientHelloID != utls.HelloChrome_133 {
			t.Error("Expected ClientHelloID for Chrome 138 to be HelloChrome_133")
		}
		if !strings.Contains(agent.Headers.Get("sec-ch-ua"), `"Google Chrome";v="138"`) {
			t.Errorf("sec-ch-ua header is incorrect: %s", agent.Headers.Get("sec-ch-ua"))
		}
	})

	t.Run("Closest Version Match", func(t *testing.T) {
		ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"
		agent, err := FromUserAgentString(ua, RequestTypeNavigate)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if agent.ClientHelloID != utls.HelloChrome_120 {
			t.Error("Expected ClientHelloID for Chrome 125 to be HelloChrome_120")
		}
	})

	t.Run("Successful Firefox Parse", func(t *testing.T) {
		ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0"
		agent, err := FromUserAgentString(ua, RequestTypeNavigate)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if agent.Headers.Get("sec-ch-ua") != "" {
			t.Error("Firefox should not have sec-ch-ua headers")
		}
		if agent.ClientHelloID != utls.HelloChrome_120 {
			t.Errorf("Incorrect TLS profile for Firefox, expected closest Chrome profile (120)")
		}
	})

	t.Run("Successful Safari Mobile Parse", func(t *testing.T) {
		ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 17_5_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1"

		agent, err := FromUserAgentString(ua, RequestTypeNavigate)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if agent.UserAgent != ua {
			t.Errorf("Expected UserAgent to be identical, got %s", agent.UserAgent)
		}

		if agent.ClientHelloID != utls.HelloChrome_120 {
			t.Errorf("Incorrect TLS profile for Safari, expected oldest stable Chrome profile")
		}
	})

	t.Run("Unsupported Browser", func(t *testing.T) {
		ua := "curl/7.64.1"
		_, err := FromUserAgentString(ua, RequestTypeNavigate)
		if !errors.Is(err, ErrUnsupportedBrowser) {
			t.Errorf("Expected ErrUnsupportedBrowser, got %v", err)
		}
	})

	t.Run("Unsupported Version", func(t *testing.T) {
		ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.0.0 Safari/537.36"
		_, err := FromUserAgentString(ua, RequestTypeNavigate)
		if !errors.Is(err, ErrUnsupportedVersion) {
			t.Errorf("Expected ErrUnsupportedVersion, got %v", err)
		}
	})
}
