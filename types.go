package main

type VendorDef struct {
	Prefix   string
	Length   int
	Org      string
	Registry string
}

func (v *VendorDef) PrefixString() string {
	if v == nil {
		return _nilStr
	}
	return v.Prefix
}
