package traefik

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const (
	TkFrontendsUrl = "%s/providers/%s/frontends"
)

type TraefikV1Client struct {
	apiUrl   string
	provider string
}

func NewTraefikV1Client(apiUrl string, provider string) *TraefikV1Client {
	return &TraefikV1Client{apiUrl: apiUrl, provider: provider}
}

type apiResponse map[string]struct {
	Routes map[string]struct {
		Rule string `json:"rule"`
	} `json:"routes"`
}

func (t *TraefikV1Client) Get() ([]string, error) {
	url := fmt.Sprintf(TkFrontendsUrl, t.apiUrl, t.provider)

	body, err := t.request(url)

	var parsedBody apiResponse

	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		return []string{}, errors.Wrapf(err, "API response parse failed, body: %+v", body)
	}

	hosts := t.extractHosts(t.extractRules(parsedBody))

	return hosts, nil
}

func (t *TraefikV1Client) request(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, errors.Wrap(err, "API request failed")
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (t *TraefikV1Client) extractRules(response apiResponse) []string {
	rules := make([]string, 0)

	for _, v := range response {
		for _, route := range v.Routes {
			rules = append(rules, route.Rule)
		}
	}

	return rules
}

func (t *TraefikV1Client) extractHosts(rules []string) []string {
	hosts := make([]string, 0)

	re := regexp.MustCompile(`(?m)Host:(.*?)([ &;]|$)`)

	for _, v := range rules {
		newHost := re.FindStringSubmatch(v)
		hosts = append(hosts, strings.Replace(newHost[1], "Host:", "", -1))
	}

	return hosts
}
