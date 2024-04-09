package doppler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Doppler struct {
	Project     string
	Environment string
	BasicURL    string
	Secret      string
	httpClient  *http.Client
}

func NewDoppler(project, secretKey, environment string) (*Doppler, error) {
	if project == "" {
		return nil, errors.New("project name is mandatory to instance Doppler")
	}
	if environment == "" {
		environment = "production"
	}

	return &Doppler{
		Project:     project,
		Environment: environment,
		BasicURL:    "https://api.doppler.com/v3/configs/config",
		Secret:      base64.StdEncoding.EncodeToString([]byte(secretKey + ":")),
		httpClient:  &http.Client{},
	}, nil
}

func (d *Doppler) SetHTTPClient(client *http.Client) {
	d.httpClient = client
}

func (d *Doppler) GetOne(secretKey string) (string, error) {
	secrets, err := d.GetSecrets()
	if err != nil {
		return "", err
	}
	value, exists := secrets[secretKey]
	if !exists {
		return "", fmt.Errorf("secret key %s not found", secretKey)
	}
	return value, nil
}

func (d *Doppler) GetSecrets() (map[string]string, error) {
	if d.httpClient == nil {
		return nil, errors.New("the http client is nil, please set it with an HTTP client")
	}

	params := url.Values{}
	params.Add("project", d.Project)
	params.Add("config", d.Environment)
	url := fmt.Sprintf("%s/secrets?%s", d.BasicURL, params.Encode())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Basic "+d.Secret)
	req.Header.Add("Accepts", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch secrets, status code: %d", resp.StatusCode)
	}

	var secretsResponse struct {
		Secrets map[string]struct {
			Computed string `json:"computed"`
		} `json:"secrets"`
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &secretsResponse); err != nil {
		return nil, err
	}

	secrets := make(map[string]string)
	for k, v := range secretsResponse.Secrets {
		secrets[k] = v.Computed
	}

	return secrets, nil
}
