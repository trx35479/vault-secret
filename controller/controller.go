package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/vault-secret/apis"
)

const ServiceAccountTokenPath = "/var/run/secrets/data/token"

type Headers struct {
	//	VaultToken string
	Url string
	//	Role       string
	//	Namespace  string
	Path string
}

//type VaultToken interface {
//	GetAuthPath(string) string
//	GetDataPath(string) string
//}

// GetAuthPath onstruct a url string for auth endpoit
func (s Headers) GetAuthPath(p string) string {
	return fmt.Sprintf("%s/v1/auth/%s/%s", s.Url, s.Path, p)
}

// GetDataPath construct a url string for data endpoint
func (s Headers) GetDataPath(p string) string {
	return fmt.Sprintf("%s/v1/%s/%s", s.Url, s.Path, p)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Get the token mounted to volume
	// Pod should have an access to service account and token usually mounted to
	// /var/run/secrets/kubernetes.io/serviceaccount/token
	lines, err := apis.GetJwt("token")
	if err != nil {
		log.Fatal(err)
	}

	// GetJwt function returns a []string
	// Expect that Container will only mount the Token of Service Account
	var jwt string
	for _, line := range lines {
		jwt = line
	}

	apiUrl := Headers{os.Getenv("VAULT_ADDR"), os.Getenv("VAULT_PATH")}
	p := apiUrl.GetAuthPath("login")

	// Get the token
	token, err := apis.GetClientToken(jwt, os.Getenv("VAULT_ROLE"), p, os.Getenv("VAULT_NAMESPACE"))
	if err != nil {
		log.Fatal(err)
	}

	// Initialise another header method and use the return token of the function
	dataUrl := Headers{os.Getenv("VAULT_ADDR"), os.Getenv("VAULT_PATH")}
	d := dataUrl.GetDataPath("secrets")

	// Get Payload
	payload, err := apis.GetData(token, d, os.Getenv("NAMESPACE"))
	if err != nil {
		log.Fatal(err)
	}

	m := payload.Data.(map[string]interface{})

	data, err := json.Marshal(m)

	io.WriteString(w, fmt.Sprint(bytes.NewBuffer(data)))
}

func Logging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("HTTP Request: %s %s  %s %s\n", r.RemoteAddr, r.Method, r.URL, r.UserAgent())
		h.ServeHTTP(w, r)
	})
}

func LogFile(logfile string) {
	if logfile != "" {
		lf, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)

		if err != nil {
			log.Fatal("Logfile: os.OpenFile:", err)
		}
		log.SetOutput(lf)
	}
}
