package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

const (
	TkFrontendsUrl = "/providers/docker/frontends"
)

func GetHosts(apiUrl string) ([]string, error) {
	resp, err := http.Get(apiUrl + TkFrontendsUrl)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var parsedBody map[string]struct {
		Routes map[string]struct{ Rule string `json:"rule"` } `json:"routes"`
	}

	err = json.Unmarshal([]byte(body), &parsedBody)
	if err != nil {
		return []string{}, err
	}

	hosts := make([]string, 0)

	for _, v := range parsedBody {
		for _, route := range v.Routes {
			hosts = append(hosts, strings.Replace(route.Rule, "Host:", "", -1))
		}
	}


	sort.Strings(hosts)

	return hosts, nil
}
