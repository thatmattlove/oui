package main

import (
	"fmt"
	"net/url"
)

const (
	REGISTRY_OUI   string = "MA-L"
	REGISTRY_OUI36 string = "MA-S"
	REGISTRY_OUI28 string = "MA-M"
	REGISTRY_IAB   string = "IAB"
	REGISTRY_CID   string = "CID"
)

type Registry struct {
	Name          string
	BaseURL       string
	FilePrefix    string
	FileExtension string
}

func (reg *Registry) URL() *url.URL {
	s := fmt.Sprintf("%s/%s.%s", reg.BaseURL, reg.FilePrefix, reg.FileExtension)
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}

func (reg *Registry) FileName() string {
	return fmt.Sprintf("%s.%s", reg.FilePrefix, reg.FileExtension)
}

func (reg *Registry) TempFilePattern() string {
	return fmt.Sprintf("*-%s.%s", reg.FilePrefix, reg.FileExtension)
}

func Registries() []*Registry {
	return []*Registry{
		{Name: REGISTRY_OUI, BaseURL: "https://standards-oui.ieee.org/oui", FilePrefix: "oui", FileExtension: "csv"},
		{Name: REGISTRY_CID, BaseURL: "https://standards-oui.ieee.org/cid", FilePrefix: "cid", FileExtension: "csv"},
		{Name: REGISTRY_IAB, BaseURL: "https://standards-oui.ieee.org/iab", FilePrefix: "iab", FileExtension: "csv"},
		{Name: REGISTRY_OUI28, BaseURL: "https://standards-oui.ieee.org/oui28", FilePrefix: "mam", FileExtension: "csv"},
		{Name: REGISTRY_OUI36, BaseURL: "https://standards-oui.ieee.org/oui36", FilePrefix: "oui36", FileExtension: "csv"},
	}
}
