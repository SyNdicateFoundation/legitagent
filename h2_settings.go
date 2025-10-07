package legitagent

import "golang.org/x/net/http2"

func GetChromiumH2Settings() map[http2.SettingID]uint32 {
	return map[http2.SettingID]uint32{
		http2.SettingHeaderTableSize:      65536,
		http2.SettingMaxConcurrentStreams: 1000,
		http2.SettingInitialWindowSize:    6291456,
		http2.SettingMaxHeaderListSize:    262144,
	}
}

func GetGeckoH2Settings() map[http2.SettingID]uint32 {
	return map[http2.SettingID]uint32{
		http2.SettingHeaderTableSize:      65536,
		http2.SettingMaxConcurrentStreams: 1000,
		http2.SettingInitialWindowSize:    131072,
		http2.SettingMaxHeaderListSize:    262144,
	}
}

func GetWebKitH2Settings() map[http2.SettingID]uint32 {
	return map[http2.SettingID]uint32{
		http2.SettingHeaderTableSize:      4096,
		http2.SettingMaxConcurrentStreams: 100,
		http2.SettingInitialWindowSize:    2097152,
		http2.SettingMaxHeaderListSize:    16384,
	}
}
