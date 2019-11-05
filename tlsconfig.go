package main

import (
	"crypto/tls"
	"errors"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type tlsConfig struct{}

func (cfg *tlsConfig) GetCertificate(domain string) (*tls.Certificate, error) {
	var instance hostingInstance
	db.Where(&hostingInstance{Domain: domain}).First(&instance)
	if instance.ID == 0 {
		return nil, errors.New("no such domain: " + domain)
	}

	var provision certificateProvision
	db.Where(&certificateProvision{CertificateID: instance.CertificateID}).
		Order("revision desc").
		First(&provision)
	if provision.ID == 0 {
		return nil, errors.New("unable to locate tls certificate")
	}

	cert, err := tls.X509KeyPair(provision.Certificate, provision.PrivateKey)
	return &cert, err
}

func (cfg *tlsConfig) GetCertificateFunc() func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		return cfg.GetCertificate(clientHello.ServerName)
	}
}
