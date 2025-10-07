package legitagent

import (
	"math"
	"sort"

	"github.com/SyNdicateFoundation/fastrand"
)

type HeaderSorter func(keys []string)

var headerPriority = map[string]int{
	":authority": 0, ":method": 1, ":path": 2, ":scheme": 3, ":status": 4,
	"host": 10, "connection": 11, "upgrade": 12, "upgrade-insecure-requests": 13, "user-agent": 14,
	"sec-ch-ua": 15, "sec-ch-ua-arch": 16, "sec-ch-ua-bitness": 17, "sec-ch-ua-full-version": 18, "sec-ch-ua-full-version-list": 19, "sec-ch-ua-mobile": 20, "sec-ch-ua-model": 21, "sec-ch-ua-platform": 22, "sec-ch-ua-platform-version": 23, "sec-ch-ua-wow64": 24,
	"authorization": 30, "proxy-authorization": 31, "cookie": 32, "sec-gpc": 33, "expect": 34, "max-forwards": 35, "from": 36,
	"accept": 40, "accept-charset": 41, "accept-encoding": 42, "accept-language": 43, "te": 44,
	"if-match": 50, "if-none-match": 51, "if-modified-since": 52, "if-unmodified-since": 53, "if-range": 54,
	"range":          60,
	"sec-fetch-site": 65, "sec-fetch-mode": 66, "sec-fetch-user": 67, "sec-fetch-dest": 68,
	"referer":      70,
	"content-type": 80, "content-length": 81, "content-encoding": 82, "content-language": 83, "content-location": 84, "content-md5": 85, "content-range": 86, "transfer-encoding": 87,
	"date": 100, "location": 101, "retry-after": 102, "set-cookie": 103, "expires": 104, "pragma": 105, "cache-control": 106, "etag": 107, "last-modified": 108, "age": 109, "vary": 110, "accept-ranges": 111, "allow": 112, "server": 113, "via": 114, "warning": 115,
	"strict-transport-security": 120, "content-security-policy": 121,
	"permissions-policy": 123, "cross-origin-opener-policy": 124, "cross-origin-resource-policy": 125, "cross-origin-embedder-policy": 126, "x-frame-options": 127, "x-content-type-options": 128, "x-xss-protection": 129, "report-to": 130, "reporting-endpoints": 131,
	"www-authenticate": 140, "proxy-authenticate": 141,
	"accept-ch": 150,
	"alt-svc":   160,
	"trailer":   170, "x-ua-compatible": 171,
}

func RandomHeaderSorter(keys []string) {
	fastrand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})
}

func PriorityHeaderSorter(keys []string) {
	sort.SliceStable(keys, func(i, j int) bool {
		p1, ok1 := headerPriority[keys[i]]
		if !ok1 {
			p1 = math.MaxInt
		}
		p2, ok2 := headerPriority[keys[j]]
		if !ok2 {
			p2 = math.MaxInt
		}
		return p1 < p2
	})
}

func ShuffledPriorityHeaderSorter(keys []string) {
	PriorityHeaderSorter(keys)

	start := 0
	for i := 1; i < len(keys); i++ {
		p1, _ := headerPriority[keys[i-1]]
		p2, _ := headerPriority[keys[i]]
		if p1 != p2 {
			if i-start > 1 {
				fastrand.Shuffle(i-start, func(a, b int) {
					keys[start+a], keys[start+b] = keys[start+b], keys[start+a]
				})
			}
			start = i
		}
	}
	if len(keys)-start > 1 {
		fastrand.Shuffle(len(keys)-start, func(a, b int) {
			keys[start+a], keys[start+b] = keys[start+b], keys[start+a]
		})
	}
}
