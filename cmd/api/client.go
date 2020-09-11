package api

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
	TkFrontendsUrl = "%s/providers/%s/frontends"
)

type apiResponse map[string]struct {
	Routes map[string]struct {
		Rule string `json:"rule"`
	} `json:"routes"`
}

func GetHosts(apiUrl, provider string) ([]string, error) {
	url := fmt.Sprintf(TkFrontendsUrl, apiUrl, provider)

	body, err := get(url)

	var parsedBody apiResponse

	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		return []string{}, errors.Wrapf(err, "API response parse failed, body: %+v", body)
	}

	hosts := extractHosts(extractRules(parsedBody))

	return hosts, nil
}

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, errors.Wrap(err, "API request failed")
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func extractRules(response apiResponse) []string {
	rules := make([]string, 0)

	for _, v := range response {
		for _, route := range v.Routes {
			rules = append(rules, route.Rule)
		}
	}

	return rules
}

func extractHosts(rules []string) []string {
	hosts := make([]string, 0)

	re := regexp.MustCompile(`(?m)Host:(.*?)([ &;]|$)`)

	for _, v := range rules {
		newHost := re.FindStringSubmatch(v)
		hosts = append(hosts, strings.Replace(newHost[1], "Host:", "", -1))
	}

	sort.Strings(hosts)

	return hosts
}
