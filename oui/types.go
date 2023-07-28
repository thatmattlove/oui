package oui

type VendorDef struct {
	Prefix   string
	Length   int
	Org      string
	Registry string
}

func (v *VendorDef) PrefixString() string {
	if v == nil {
		return "<nil>"
	}
	return v.Prefix
}
