package traefik_v2

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

const (
	TkRoutersUrl = "%s/api/http/routers"
)

type TraefikV2Client struct {
	apiUrl string
}

func NewTraefikV2Client(apiUrl string) *TraefikV2Client {
	return &TraefikV2Client{apiUrl: apiUrl}
}

type apiResponse []struct {
	Rule string `json:"rule"`
}

func (t *TraefikV2Client) Get() ([]string, error) {
	url := fmt.Sprintf(TkRoutersUrl, t.apiUrl)

	body, err := t.request(url)

	var parsedBody apiResponse

	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		return []string{}, errors.Wrapf(err, "API response parse failed, body: %+v", body)
	}

	hosts := t.extractHosts(t.extractRules(parsedBody))

	return hosts, nil
}

func (t *TraefikV2Client) request(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, errors.Wrap(err, "API request failed")
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (t *TraefikV2Client) extractRules(response apiResponse) []string {
	rules := make([]string, 0)

	for _, v := range response {
		rules = append(rules, v.Rule)
	}

	return rules
}

func (t *TraefikV2Client) extractHosts(rules []string) []string {
	hosts := make([]string, 0)

	re := regexp.MustCompile("(?m)Host\\(`(.*?)`\\)([ &;]|$)")

	for _, v := range rules {
		if !strings.Contains(v, "Host(`") {
			continue
		}

		newHost := re.FindStringSubmatch(v)
		hosts = append(hosts, strings.Replace(newHost[1], "Host:", "", -1))
	}

	sort.Strings(hosts)

	return hosts
}
