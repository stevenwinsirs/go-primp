package primp

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

// LoadCACerts 从环境变量 PRIMP_CA_BUNDLE 加载 CA 证书
// 如果环境变量不存在，则使用系统证书
func LoadCACerts(caCertPath string) (*x509.CertPool, error) {
	if caCertPath == "" {
		caCertPath = os.Getenv("PRIMP_CA_BUNDLE")
		if caCertPath == "" {
			caCertPath = os.Getenv("CA_CERT_FILE")
		}
	}

	if caCertPath != "" {
		// 从文件加载自定义 CA 证书
		caCert, err := os.ReadFile(caCertPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA cert file: %w", err)
		}

		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			return nil, fmt.Errorf("failed to append CA certs")
		}

		return caCertPool, nil
	}

	// 使用系统证书作为备选
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to load system cert pool: %w", err)
	}

	return caCertPool, nil
}

// CreateTLSConfig 使用指定的 CA 证书创建 TLS 配置
func CreateTLSConfig(verify bool) (*tls.Config, error) {
	if !verify {
		return &tls.Config{
			InsecureSkipVerify: true,
		}, nil
	}

	caCertPool, err := LoadCACerts("")
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		RootCAs: caCertPool,
	}, nil
}
