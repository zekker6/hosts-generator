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

func GetHosts(apiUrl, provider string) ([]string, error) {
	url := fmt.Sprintf(TkFrontendsUrl, apiUrl, provider)

	resp, err := http.Get(url)
	if err != nil {
		return []string{}, errors.Wrap(err, "API request failed")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var parsedBody map[string]struct {
		Routes map[string]struct{ Rule string `json:"rule"` } `json:"routes"`
	}

	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		return []string{}, errors.Wrapf(err, "API response parse failed, body: %+v", body)
	}

	hosts := make([]string, 0)

	re := regexp.MustCompile(`(?m)Host:(.*?)([ &;]|$)`)

	for _, v := range parsedBody {
		for _, route := range v.Routes {
			newHost := re.FindAllString(route.Rule, -1)
			hosts = append(hosts, strings.Replace(newHost[0], "Host:", "", -1))
		}
	}

	sort.Strings(hosts)

	return hosts, nil
}
