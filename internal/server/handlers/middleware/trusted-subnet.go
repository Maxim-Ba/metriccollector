package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/Maxim-Ba/metriccollector/internal/server/services/subnet"
)

func TrustedSubnetMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trustedSubnet := subnet.Instance.Get()
		if trustedSubnet == "" {
			next.ServeHTTP(w, r)
			return
		}

		ipStr := strings.TrimSpace(r.Header.Get("X-Real-IP"))
		if ipStr == "" {
			http.Error(w, "X-Real-IP header is required", http.StatusForbidden)
			return
		}

		ip := net.ParseIP(ipStr)
		if ip == nil {
			http.Error(w, "Invalid IP format in X-Real-IP", http.StatusBadRequest)
			return
		}

		_, cidrNet, err := net.ParseCIDR(trustedSubnet)
		if err != nil {
			http.Error(w, "Invalid trusted_subnet CIDR format", http.StatusInternalServerError)
			return
		}

		if !cidrNet.Contains(ip) {
			http.Error(w, "IP not in trusted subnet", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)

	})
}
