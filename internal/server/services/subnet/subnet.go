package subnet

type TrustedSubnet struct {
	trustedSubnet string
}

var Instance TrustedSubnet

func New(trustedSubnet string) *TrustedSubnet {
	return &TrustedSubnet{trustedSubnet: trustedSubnet}
}
func (s *TrustedSubnet) Get() string {
	return s.trustedSubnet
}
