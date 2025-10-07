package legitagent

import (
	"github.com/SyNdicateFoundation/fastrand"
	"golang.org/x/net/http2"
	"math"
)

func randomizeValue(base uint32, percentage float64) uint32 {
	if base == 0 {
		return 0
	}

	delta := uint32(float64(base) * percentage)
	minX := base - delta
	if base < delta {
		minX = 1
	}

	maxX := base + delta
	return fastrand.Number(minX, maxX)
}

func randomizeH2Settings(baseSettings map[http2.SettingID]uint32, profile H2RandomizationProfile) map[http2.SettingID]uint32 {
	randomized := make(map[http2.SettingID]uint32)
	for id, val := range baseSettings {
		randomized[id] = val
	}

	switch profile {
	case H2RandomizationProfileNormal:
		if val, ok := randomized[http2.SettingHeaderTableSize]; ok {
			randomized[http2.SettingHeaderTableSize] = randomizeValue(val, 0.10)
		}
		if val, ok := randomized[http2.SettingInitialWindowSize]; ok {
			randomized[http2.SettingInitialWindowSize] = randomizeValue(val, 0.15)
		}
		if val, ok := randomized[http2.SettingMaxHeaderListSize]; ok {
			randomized[http2.SettingMaxHeaderListSize] = randomizeValue(val, 0.10)
		}

	case H2RandomizationProfileMaximum:
		randomized[http2.SettingHeaderTableSize] = randomizeValue(4096, 0.20)
		randomized[http2.SettingEnablePush] = 0
		randomized[http2.SettingInitialWindowSize] = randomizeValue(65535, 0.20)
		randomized[http2.SettingMaxFrameSize] = randomizeValue(16384, 0.20)
		randomized[http2.SettingMaxConcurrentStreams] = uint32(math.MaxUint32 - fastrand.IntN(1024))
	default:
		panic("unhandled default case")
	}

	return randomized
}
