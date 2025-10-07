# legitagent

[![Go Report Card](https://goreportcard.com/badge/github.com/SyNdicateFoundation/legitagent)](https://goreportcard.com/report/github.com/SyNdicateFoundation/legitagent)
[![GoDoc](https://godoc.org/github.com/SyNdicateFoundation/legitagent?status.svg)](https://godoc.org/github.com/SyNdicateFoundation/legitagent)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Stop being fingerprinted.** `legitagent` is an advanced Go library for generating realistic, difficult-to-fingerprint browser agents.

While many libraries can generate a `User-Agent` string, modern bot detection and web scrapers look far deeper. They analyze the entire network fingerprint:

-   **TLS ClientHello (JA3/JA4 Hash):** The initial handshake reveals dozens of data points about your client.
-   **HTTP/2 Settings:** The specific parameters of your H2 connection are a strong fingerprint.
-   **HTTP Header Ordering:** The exact order of your headers is another tell-tale sign.

`legitagent` is architecturally designed to solve this problem by generating **holistic, consistent, and dynamic browser profiles**, making your requests appear legitimately human.

## Key Features

-   **Browser-Specific Network Fingerprints:** Generates agents with matching TLS and HTTP/2 fingerprints for each browser family (Chromium, Gecko, WebKit). A Firefox agent *acts* like Firefox on the network level.
-   **Dynamic JA3 Anti-Fingerprinting:** An optional "Maximum Stealth" mode dynamically shuffles TLS extensions and cipher suites on every generation, making consistent JA3 fingerprinting impossible.
-   **Plausible Header Randomization:** Mimics real-world browser behavior by subtly shuffling HTTP headers within their standard priority groups, avoiding the static fingerprint of a fixed order.
-   **Expansive and Accurate Profiles:** Includes a wide range of common desktop and mobile browser versions (Chrome, Firefox, Safari, Edge, etc.) with plausible build numbers and User-Agent formats.
-   **High Performance:** Built with performance in mind, using `sync.Pool` for agent and string builder reuse to achieve zero-allocation header generation in hot loops.
-   **Rich Customization:** Provides a clean, option-based API to control the browser, OS, platform, version range, and anti-fingerprinting level.

## Installation

```sh
go get github.com/SyNdicateFoundation/legitagent
```

## Quick Start

Generate a random, legitimate browser agent with just a few lines of code.

```go
package main

import (
	"fmt"
	"github.com/SyNdicateFoundation/legitagent"
	"log"
)

func main() {
	// 1. Create a new generator
	g := legitagent.NewGenerator()

	// 2. Generate an agent
	agent, err := g.Generate()
	if err != nil {
		log.Fatalf("Failed to generate agent: %v", err)
	}

	// 3. IMPORTANT: Defer the release of the agent back to the pool
	defer g.ReleaseAgent(agent)

	// 4. Use the generated data
	fmt.Printf("User-Agent: %s\n", agent.UserAgent)
	fmt.Printf("Accept-Language Header: %s\n", agent.Headers.Get("accept-language"))
	fmt.Printf("TLS ClientHello ID: %v\n", agent.ClientHelloID)
}
```

## Advanced Usage & Examples

### Example 1: Generate a Specific Browser

Generate a specific version of Firefox running on Linux.

```go
g := legitagent.NewGenerator(
    legitagent.WithBrowsers(legitagent.BrowserFirefox),
    legitagent.WithVersionRange(128, 128),
    legitagent.WithOS(legitagent.OSLinux),
    legitagent.WithPlatforms(legitagent.PlatformDesktop),
)

agent, err := g.Generate()
// ...
// User-Agent will be something like:
// Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0
```

### Example 2: Generate a Mobile Safari Agent

Generate a Safari agent running on a mobile iOS device.

```go
g := legitagent.NewGenerator(
    legitagent.WithBrowsers(legitagent.BrowserSafari),
    legitagent.WithPlatforms(legitagent.PlatformMobile),
)

agent, err := g.Generate()
// ...
// User-Agent will be something like:
// Mozilla/5.0 (iPhone; CPU iPhone OS 17_5_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1
```

### Example 3: Maximum Anti-Fingerprinting (Dynamic JA3)

For maximum stealth, enable the dynamic fingerprinting profile. This will randomize the TLS fingerprint (JA3) and header order on **every single generation**.

```go
g := legitagent.NewGenerator(
    legitagent.WithFingerprintProfile(legitagent.FingerprintProfileMaximum),
)

// Each of these agents will have a different JA3 hash and header order
agent1, _ := g.Generate()
agent2, _ := g.Generate()
```

### Example 4: Using with a Custom HTTP Client

This is how you apply the full network fingerprint to a real HTTP request using `utls`.

```go
import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/SyNdicateFoundation/legitagent"
	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
)

func main() {
	g := legitagent.NewGenerator()
	agent, err := g.Generate()
	if err != nil {
		log.Fatalf("Failed to generate agent: %v", err)
	}
	defer g.ReleaseAgent(agent)

	// Create a custom DialTLSContext function that uses the agent's TLS fingerprint
	dialTLSContext := func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
		rawConn, err := net.DialTimeout(network, addr, 10*time.Second)
		if err != nil {
			return nil, err
		}

		uTLSConfig := &utls.Config{
			ServerName:         cfg.ServerName,
			InsecureSkipVerify: true,
			NextProtos:         []string{"h2", "http/1.1"},
		}

		uconn := utls.UClient(rawConn, uTLSConfig, agent.ClientHelloID)
		if err := uconn.HandshakeContext(ctx); err != nil {
			return nil, err
		}
		return uconn, nil
	}
	
	// Create the HTTP/2 transport with the custom dialer
	h2Transport := &http2.Transport{
		DialTLSContext:  dialTLSContext,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	
	// Apply the agent's browser-specific H2 settings
	for id, val := range agent.H2Settings {
		h2Transport.Settings = append(h2Transport.Settings, http2.Setting{ID: id, Val: val})
	}
	
	client := &http.Client{
		Transport: h2Transport,
		Timeout:   15 * time.Second,
	}
	
	req, _ := http.NewRequest(http.MethodGet, "https://cloudflare.com/cdn-cgi/trace", nil)
	
	// Apply the generated headers
	req.Header = agent.Headers
	req.Header.Set("User-Agent", agent.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}
```

## Detailed Options

Customize the generator using these `Option` functions:

-   `WithBrowsers(...Browser)`: Specifies which browsers to choose from (e.g., `BrowserChrome`, `BrowserFirefox`).
-   `WithPlatforms(...Platform)`: Specifies the platform (e.g., `PlatformDesktop`, `PlatformMobile`).
-   `WithOS(...OperatingSystem)`: Specifies the operating system (e.g., `OSWindows11`, `OSMac`, `OSiOS`).
-   `WithVersionRange(min, max int)`: Constrains the major version of the generated browser.
-   `WithLanguages(...string)`: Sets the `Accept-Language` profiles to use (e.g., `"fr-FR,fr;q=0.9"`).
-   `WithFullFingerprint(bool)`: Toggles the inclusion of extended `sec-ch-ua-*` headers for a more detailed fingerprint.
-   `WithH2Only(bool)`: (Default: `true`) Ensures only browsers that support HTTP/2 are generated.
-   `WithFingerprintProfile(profile FingerprintProfile)`: Sets the anti-fingerprinting level (`FingerprintProfileNormal` or `FingerprintProfileMaximum`).
-   `WithH2Randomization(profile H2RandomizationProfile)`: Controls the randomization of HTTP/2 settings to further prevent fingerprinting.
    -   `H2RandomizationProfileNone` (Default): Uses the exact, default settings for the generated browser.
    -   `H2RandomizationProfileNormal`: Applies small, realistic variations to the browser's default settings.
    -   `H2RandomizationProfileMaximum`: Uses a high-throughput profile with aggressive, randomized values, including setting `MaxConcurrentStreams` to `math.MaxUint32` minus a small random salt.

## Architectural Philosophy: Consistency is Legitimacy

The core design principle of `legitagent` is that a legitimate fingerprint must be **consistent across all network layers**. It is not enough to send a Firefox User-Agent; you must also send a Firefox TLS fingerprint *and* a Firefox HTTP/2 fingerprint.

This library solves the common pitfall where a program sends a Firefox TLS fingerprint but uses Go's default HTTP/2 settings, which look like Chrome. This mismatch is an immediate red flag for sophisticated bot detectors. `legitagent` ensures that all layers of the network stack are perfectly aligned with the chosen browser family, providing a truly legitimate and difficult-to-detect fingerprint.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)
