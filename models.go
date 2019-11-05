package main

import "github.com/vmihailenco/msgpack"

type hostingInstance struct {
	ID            int64
	Domain        string
	CertificateID int64
}

type certificateProvision struct {
	ID            int64
	CertificateID int64
	Revision      int
	Certificate   []byte
	PrivateKey    []byte
}

type hostingBackend struct {
	InstanceID int64
	Properties []byte
}

func (h *hostingBackend) getProps() (*hostingBackendProps, error) {
	var props *hostingBackendProps
	err := msgpack.Unmarshal(h.Properties, &props)
	if err != nil {
		return nil, err
	}
	return props, nil
}

type hostingBackendProps struct {
	BucketName      string `msgpack:"b"`
	FilePrefix      string `msgpack:"p"`
	RedirectToIdnex bool   `msgpack:"i"`
}
