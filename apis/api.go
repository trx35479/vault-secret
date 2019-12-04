package apis

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/vault-secret/models"
)

// Login function to get the needed token before data can be provided by vault
// Note that we will use the kubernetes auth on vault
func GetClientToken(Jwt, Role, Url, Namespace string) (string, error) {
	client := &http.Client{}
	requestBody, err := json.Marshal(map[string]string{
		"jwt":  Jwt,
		"role": Role,
	})
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", Url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal(err)
	}

	// Namespace is optional in this case
	if Namespace != "null" {
		// Add header to the request (needed if the vault is multi-tenant)
		req.Header.Add("X-Vault-Namespace", Namespace)
	}
	// Send the request to vault
	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var token *models.Payload
	err = json.Unmarshal([]byte(body), &token)
	if err != nil {
		log.Fatal(err)
	}

	var clientToken string

	clientToken = token.Auth.ClientToken

	return clientToken, nil
}

// GetData function to get the secret data from vault
// This shopud be executed after the login is successful
func GetData(Token, Url, Namespace string) (*models.Payload, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Add the header X-Vault-Token in the http request
	req.Header.Add("X-Vault-Token", Token)

	// This header is needed in vault enterprise
	if Namespace != "null" {
		// Add header to the request (needed if the vault is multi-tenant)
		req.Header.Add("X-Vault-Namespace", Namespace)
	}

	// Send the request to Vault
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data *models.Payload
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		log.Fatal(err)
	}

	return data, nil
}
