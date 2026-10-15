package config

import (
	"net/http"

	"github.com/prometheus/client_golang/api"
)

// authTransport es un wrapper para añadir Basic Auth a las peticiones
type authTransport struct {
	username string
	password string
	rt       http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.username, t.password)
	return t.rt.RoundTrip(req)
}

func GetPrometheusClient() (api.Client, error) {
	address := "http://212.56.40.50:9090"
	username := "admin"
	password := "$Yc0m-2023"

	return api.NewClient(api.Config{
		Address: address,
		// Inyectamos el transporte personalizado con las credenciales
		RoundTripper: &authTransport{
			username: username,
			password: password,
			rt:       http.DefaultTransport,
		},
	})
}
