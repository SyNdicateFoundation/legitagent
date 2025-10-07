package legitagent

import (
	"github.com/SyNdicateFoundation/fastrand"
	utls "github.com/refraction-networking/utls"
)

func shuffleExtensions(extensions []utls.TLSExtension) []utls.TLSExtension {
	shuffled := make([]utls.TLSExtension, len(extensions))
	copy(shuffled, extensions)
	fastrand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	final := make([]utls.TLSExtension, 0, len(shuffled)+3)
	final = append(final, &utls.UtlsGREASEExtension{})
	final = append(final, shuffled...)
	final = append(final, &utls.UtlsGREASEExtension{})
	final = append(final, &utls.UtlsPaddingExtension{GetPaddingLen: utls.BoringPaddingStyle})

	return final
}

func ChromeLatestSpec() *utls.ClientHelloSpec {
	extensions := []utls.TLSExtension{
		&utls.SNIExtension{},
		&utls.ExtendedMasterSecretExtension{},
		&utls.RenegotiationInfoExtension{Renegotiation: utls.RenegotiateOnceAsClient},
		&utls.SupportedCurvesExtension{Curves: []utls.CurveID{
			utls.GREASE_PLACEHOLDER, utls.X25519, utls.CurveP256, utls.CurveP384,
		}},
		&utls.SupportedPointsExtension{SupportedPoints: []byte{0}},
		&utls.SessionTicketExtension{},
		&utls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
		&utls.StatusRequestExtension{},
		&utls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []utls.SignatureScheme{
			utls.ECDSAWithP256AndSHA256, utls.PSSWithSHA256, utls.PKCS1WithSHA256,
			utls.ECDSAWithP384AndSHA384, utls.PSSWithSHA384, utls.PKCS1WithSHA384,
			utls.PSSWithSHA512, utls.PKCS1WithSHA512,
		}},
		&utls.SCTExtension{},
		&utls.KeyShareExtension{KeyShares: []utls.KeyShare{
			{Group: utls.CurveID(utls.GREASE_PLACEHOLDER), Data: []byte{0}},
			{Group: utls.X25519},
		}},
		&utls.PSKKeyExchangeModesExtension{Modes: []uint8{utls.PskModeDHE}},
		&utls.SupportedVersionsExtension{Versions: []uint16{
			utls.GREASE_PLACEHOLDER, utls.VersionTLS13, utls.VersionTLS12,
		}},
		&utls.UtlsCompressCertExtension{Algorithms: []utls.CertCompressionAlgo{utls.CertCompressionBrotli}},
	}

	cipherSuites := []uint16{
		utls.GREASE_PLACEHOLDER,
		utls.TLS_AES_128_GCM_SHA256,
		utls.TLS_AES_256_GCM_SHA384,
		utls.TLS_CHACHA20_POLY1305_SHA256,
		utls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		utls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		utls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		utls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		utls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
		utls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		utls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		utls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		utls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		utls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		utls.TLS_RSA_WITH_AES_128_CBC_SHA,
		utls.TLS_RSA_WITH_AES_256_CBC_SHA,
	}

	shuffledCiphers := make([]uint16, len(cipherSuites))
	copy(shuffledCiphers, cipherSuites)
	fastrand.Shuffle(len(shuffledCiphers), func(i, j int) {
		shuffledCiphers[i], shuffledCiphers[j] = shuffledCiphers[j], shuffledCiphers[i]
	})

	return &utls.ClientHelloSpec{
		CipherSuites:       shuffledCiphers,
		CompressionMethods: []byte{0x00},
		Extensions:         shuffleExtensions(extensions),
		GetSessionID:       nil,
	}
}
