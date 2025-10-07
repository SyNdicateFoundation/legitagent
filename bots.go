package legitagent

import (
	utls "github.com/refraction-networking/utls"
)

type botProfile struct {
	UserAgent string
	HelloID   utls.ClientHelloID
	Headers   map[string]string
}

const (
	BotAhrefs         = "AhrefsBot"
	BotApple          = "AppleBot"
	BotBaidu          = "BaiduBot"
	BotBing           = "BingBot"
	BotBytespider     = "BytespiderBot"
	BotCC             = "CCBot"
	BotChatGPT        = "ChatGPTUser"
	BotClaude         = "ClaudeBot"
	BotCohere         = "CohereBot"
	BotDiffbot        = "Diffbot"
	BotDuckDuckGo     = "DuckDuckGoBot"
	BotFacebook       = "FacebookBot"
	BotGPT            = "GPTBot"
	BotGoogle         = "GoogleBot"
	BotGoogleExtended = "GoogleExtended"
	BotLinkedIn       = "LinkedInBot"
	BotMajestic       = "MajesticBot"
	BotMoz            = "MozBot"
	BotPerplexity     = "PerplexityBot"
	BotPetal          = "PetalBot"
	BotPinterest      = "PinterestBot"
	BotSemrush        = "SemrushBot"
	BotSogou          = "SogouBot"
	BotTwitter        = "TwitterBot"
	BotUptimeRobot    = "UptimeRobot"
	BotWhatsApp       = "WhatsAppBot"
	BotYahoo          = "YahooBot"
	BotYandex         = "YandexBot"
	BotYou            = "YouBot"
)

var botProfileCategories map[string][]botProfile
var allBotProfiles []botProfile

func init() {
	baseBotHeaders := map[string]string{
		"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"accept-encoding": "gzip, deflate",
	}

	botProfileCategories = map[string][]botProfile{
		// --- Google ---
		BotGoogle: {
			{UserAgent: "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)", HelloID: utls.HelloGolang, Headers: map[string]string{
				"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
				"accept-encoding": "gzip, deflate",
				"from":            "googlebot@googlebot.com",
			}},
			{UserAgent: "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Googlebot/2.1; +http://www.google.com/bot.html) Chrome/120.0.0.0 Safari/537.36", HelloID: utls.HelloChrome_120, Headers: map[string]string{
				"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
				"accept-encoding":           "gzip, deflate, br",
				"accept-language":           "en-US,en;q=0.9",
				"upgrade-insecure-requests": "1",
				"sec-fetch-site":            "none",
				"sec-fetch-mode":            "navigate",
				"sec-fetch-user":            "?1",
				"sec-fetch-dest":            "document",
				"sec-ch-ua":                 `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`,
				"sec-ch-ua-mobile":          "?0",
				"sec-ch-ua-platform":        `"Linux"`,
				"from":                      "googlebot@googlebot.com",
			}},
			{UserAgent: "Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)", HelloID: utls.HelloChrome_120, Headers: map[string]string{
				"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
				"accept-encoding":           "gzip, deflate, br",
				"accept-language":           "en-US,en;q=0.9",
				"upgrade-insecure-requests": "1",
				"sec-fetch-site":            "none",
				"sec-fetch-mode":            "navigate",
				"sec-fetch-user":            "?1",
				"sec-fetch-dest":            "document",
				"sec-ch-ua":                 `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`,
				"sec-ch-ua-mobile":          "?1",
				"sec-ch-ua-platform":        `"Android"`,
				"from":                      "googlebot@googlebot.com",
			}},
			{UserAgent: "Googlebot-Image/1.0", HelloID: utls.HelloGolang, Headers: map[string]string{
				"accept":          "image/*",
				"accept-encoding": "gzip, deflate",
				"from":            "googlebot@googlebot.com",
			}},
			{UserAgent: "Googlebot-News", HelloID: utls.HelloGolang, Headers: map[string]string{
				"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
				"accept-encoding": "gzip, deflate",
				"from":            "googlebot@googlebot.com",
			}},
			{UserAgent: "Googlebot-Video/1.0", HelloID: utls.HelloGolang, Headers: map[string]string{
				"accept":          "video/*",
				"accept-encoding": "gzip, deflate",
				"from":            "googlebot@googlebot.com",
			}},
			{UserAgent: "Mediapartners-Google", HelloID: utls.HelloGolang, Headers: map[string]string{
				"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
				"accept-encoding": "gzip, deflate",
				"from":            "googlebot@googlebot.com",
			}},
			{
				UserAgent: "AdsBot-Google (+http://www.google.com/adsbot.html)",
				HelloID:   utls.HelloGolang,
				Headers: map[string]string{
					"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
					"accept-encoding": "gzip, deflate",
					"from":            "googlebot@googlebot.com",
				},
			},
			{
				UserAgent: "FeedFetcher-Google; (+http://www.google.com/feedfetcher.html)",
				HelloID:   utls.HelloGolang,
				Headers:   map[string]string{"accept": "application/atom+xml,application/rss+xml,application/xml;q=0.9,*/*;q=0.8", "accept-encoding": "gzip, deflate"},
			},
		},
		BotGoogleExtended: {
			{UserAgent: "Google-Extended", HelloID: utls.HelloGolang, Headers: baseBotHeaders},
		},

		// --- Microsoft ---
		BotBing: {
			{UserAgent: "Mozilla/5.0 (compatible; Bingbot/2.0; +http://www.bing.com/bingbot.htm)", HelloID: utls.HelloGolang, Headers: map[string]string{
				"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
				"accept-encoding": "gzip, deflate",
				"accept-language": "en-US,en;q=0.9",
			}},
			{UserAgent: "Mozilla/5.0 (compatible; Bingbot/2.0; +http://www.bing.com/bingbot.htm)", HelloID: utls.HelloEdge_106, Headers: map[string]string{
				"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
				"accept-encoding":           "gzip, deflate, br",
				"accept-language":           "en-US,en;q=0.9",
				"upgrade-insecure-requests": "1",
				"sec-fetch-site":            "none",
				"sec-fetch-mode":            "navigate",
				"sec-fetch-user":            "?1",
				"sec-fetch-dest":            "document",
				"sec-ch-ua":                 `"Microsoft Edge";v="106", "Chromium";v="106", "Not;A=Brand";v="99"`,
				"sec-ch-ua-mobile":          "?0",
				"sec-ch-ua-platform":        `"Linux"`,
			}},
			{UserAgent: "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534+ (KHTML, like Gecko) BingPreview/1.0b", HelloID: utls.HelloEdge_106, Headers: map[string]string{
				"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
				"accept-encoding":           "gzip, deflate, br",
				"accept-language":           "en-US,en;q=0.9",
				"upgrade-insecure-requests": "1",
				"sec-fetch-site":            "none",
				"sec-fetch-mode":            "navigate",
				"sec-fetch-user":            "?1",
				"sec-fetch-dest":            "document",
				"sec-ch-ua":                 `"Microsoft Edge";v="106", "Chromium";v="106", "Not;A=Brand";v="99"`,
				"sec-ch-ua-mobile":          "?0",
				"sec-ch-ua-platform":        `"Windows"`,
			}},
		},

		// --- Other Major Search Engines ---
		BotDuckDuckGo: {
			{UserAgent: "Mozilla/5.0 (compatible; DuckDuckBot/1.0; +http://duckduckgo.com/duckduckbot.html)", HelloID: utls.HelloGolang, Headers: baseBotHeaders},
		},
		BotBaidu: {
			{UserAgent: "Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)", HelloID: utls.HelloGolang, Headers: map[string]string{
				"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
				"accept-encoding": "gzip, deflate",
				"accept-language": "zh-CN,zh;q=0.8,en;q=0.6",
			}},
		},
		BotYandex: {
			{UserAgent: "Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)", HelloID: utls.HelloGolang, Headers: baseBotHeaders},
			{UserAgent: "Mozilla/5.0 (compatible; YandexImages/3.0; +http://yandex.com/bots)", HelloID: utls.HelloGolang, Headers: map[string]string{
				"accept":          "image/*,*/*;q=0.8",
				"accept-encoding": "gzip, deflate",
			}},
		},
		BotYahoo: {
			{UserAgent: "Mozilla/5.0 (compatible; Yahoo! Slurp; http://help.yahoo.com/help/us/ysearch/slurp)", HelloID: utls.HelloGolang, Headers: baseBotHeaders},
		},
		BotSogou: {
			{UserAgent: "Sogou web spider/4.0(+http://www.sogou.com/docs/help/webmasters.htm#07)", HelloID: utls.HelloGolang, Headers: map[string]string{
				"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
				"accept-encoding": "gzip, deflate",
				"accept-language": "zh-CN,zh;q=0.8",
			}},
		},

		// --- SEO/Marketing Bots ---
		BotAhrefs: {
			{UserAgent: "Mozilla/5.0 (compatible; AhrefsBot/7.0; +http://ahrefs.com/robot/)", HelloID: utls.HelloGolang, Headers: baseBotHeaders},
		},
		BotSemrush: {
			{UserAgent: "SemrushBot/7~bl; +http://www.semrush.com/bot.html", HelloID: utls.HelloGolang, Headers: baseBotHeaders},
		},
		BotMajestic: {
			{UserAgent: "Mozilla/5.0 (compatible; MJ12bot/v1.4.8; http://mj12bot.com/)", HelloID: utls.HelloGolang, Headers: baseBotHeaders},
		},
		BotMoz: {
			{UserAgent: "Mozilla/5.0 (compatible; DotBot/1.1; http://www.opensiteexplorer.org/dotbot, help@moz.com)", HelloID: utls.HelloGolang, Headers: baseBotHeaders},
		},

		// --- AI/LLM Bots ---
		BotGPT: {
			{UserAgent: "GPTBot/1.0 (+http://openai.com/gptbot)", HelloID: utls.HelloGolang, Headers: baseBotHeaders},
		},
		BotChatGPT: {
			{UserAgent: "ChatGPT-User", HelloID: utls.HelloGolang, Headers: baseBotHeaders},
		},
		BotClaude: {
			{UserAgent: "ClaudeBot/1.0 (+claudebot@anthropic.com)", HelloID: utls.HelloGolang, Headers: baseBotHeaders},
		},
		BotCohere: {
			{UserAgent: "cohere-ai/1.0 (+https://cohere.com/bot)", HelloID: utls.HelloGolang, Headers: baseBotHeaders},
		},
		BotPerplexity: {
			{UserAgent: "PerplexityBot/1.0 (+https://about.perplexity.ai/docs/perplexitybot)", HelloID: utls.HelloGolang, Headers: baseBotHeaders},
		},
		BotYou: {
			{UserAgent: "Mozilla/5.0 (compatible; YouBot/1.0; +http://about.you.com/youbot)", HelloID: utls.HelloGolang, Headers: baseBotHeaders},
		},
		BotDiffbot: {
			{UserAgent: "Diffbot/1.0 (+http://www.diffbot.com/our-bot/)", HelloID: utls.HelloGolang, Headers: baseBotHeaders},
		},

		// --- Social Media Bots ---
		BotFacebook: {
			{UserAgent: "facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)", HelloID: utls.HelloGolang, Headers: map[string]string{
				"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
				"accept-encoding": "gzip, deflate",
			}},
		},
		BotTwitter: {
			{UserAgent: "Twitterbot/1.0", HelloID: utls.HelloGolang, Headers: map[string]string{
				"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
				"accept-encoding": "gzip, deflate",
			}},
		},
		BotPinterest: {
			{UserAgent: "Pinterest/0.2 (+http://www.pinterest.com/bot.html)", HelloID: utls.HelloGolang, Headers: map[string]string{
				"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/*;q=0.8,*/*;q=0.7",
				"accept-encoding": "gzip, deflate",
			}},
		},
		BotLinkedIn: {
			{UserAgent: "LinkedInBot/1.0 (compatible; Mozilla/5.0; Apache-HttpClient +http://www.linkedin.com)", HelloID: utls.HelloGolang, Headers: map[string]string{
				"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
				"accept-encoding": "gzip, deflate",
			}},
		},
		BotWhatsApp: {
			{UserAgent: "WhatsApp/2.21.18.17 A", HelloID: utls.HelloGolang, Headers: map[string]string{
				"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
				"accept-encoding": "gzip, deflate",
			}},
		},

		// --- Miscellaneous Bots ---
		BotApple: {
			{UserAgent: "Mozilla/5.0 (compatible; Applebot/1.0; +http://www.apple.com/go/applebot)", HelloID: utls.HelloGolang, Headers: baseBotHeaders},
		},
		BotUptimeRobot: {
			{UserAgent: "Mozilla/5.0 (compatible; UptimeRobot/2.0; https://www.uptimerobot.com/)", HelloID: utls.HelloGolang, Headers: baseBotHeaders},
		},
		BotPetal: {
			{UserAgent: "Mozilla/5.0 (compatible; PetalBot; +http://aspiegel.com/petalbot)", HelloID: utls.HelloGolang, Headers: baseBotHeaders},
		},
		BotBytespider: {
			{UserAgent: "Mozilla/5.0 (compatible; Bytespider; +http://www.bytespider.com/)", HelloID: utls.HelloGolang, Headers: baseBotHeaders},
		},
		BotCC: {
			{UserAgent: "CCBot/2.0 (+https://commoncrawl.org/commoncrawl/projects/bots)", HelloID: utls.HelloGolang, Headers: baseBotHeaders},
		},
	}

	for _, profiles := range botProfileCategories {
		allBotProfiles = append(allBotProfiles, profiles...)
	}
}
