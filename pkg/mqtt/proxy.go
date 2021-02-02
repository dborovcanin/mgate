package mqtt

import "crypto/tls"

// Proxy represents MQTT Proxy
type ProxyIfc interface {
	Listen() error
	ListenTLS(tlsCfg *tls.Config) error
}
