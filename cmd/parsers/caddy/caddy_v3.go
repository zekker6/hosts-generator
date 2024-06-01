package caddy

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type CaddyV3 struct {
	apiUrl string
}

type Encode struct {
	Encodings map[string]struct{} `json:"encodings"`
	Handler   string              `json:"handler"`
	Prefer    []string            `json:"prefer"`
}

type Upstream struct {
	Dial string `json:"dial"`
}

type TransportTLS struct {
	InsecureSkipVerify bool `json:"insecure_skip_verify"`
}

type Transport struct {
	Protocol string       `json:"protocol"`
	TLS      TransportTLS `json:"tls"`
}

type ReverseProxy struct {
	Handler   string     `json:"handler"`
	Upstreams []Upstream `json:"upstreams"`
	Transport *Transport `json:"transport,omitempty"`
}

type SubRoute struct {
	Handler string   `json:"handler"`
	Routes  []Route  `json:"routes"`
	Group   string   `json:"group,omitempty"`
	Handle  []Handle `json:"handle,omitempty"`
	Match   []Match  `json:"match,omitempty"`
}

type Handle struct {
	Handler   string              `json:"handler"`
	Encodings map[string]struct{} `json:"encodings,omitempty"`
	Prefer    []string            `json:"prefer,omitempty"`
	Upstreams []Upstream          `json:"upstreams,omitempty"`
	Transport *Transport          `json:"transport,omitempty"`
	Routes    []Route             `json:"routes,omitempty"`
	Group     string              `json:"group,omitempty"`
}

type Match struct {
	Host []string `json:"host"`
}

type Route struct {
	Handle   []Handle `json:"handle"`
	Match    []Match  `json:"match"`
	Terminal bool     `json:"terminal,omitempty"`
	Group    string   `json:"group,omitempty"`
}

type LoggerNames struct {
	WildcardHeZekkerCloud string `json:"*.he.zekker.cloud"`
	WildcardZekkerDevTk   string `json:"*.zekker-dev.tk"`
	WildcardZekkerCloud   string `json:"*.zekker.cloud"`
}

type Logs struct {
	LoggerNames LoggerNames `json:"logger_names"`
}

type Server struct {
	Listen []string `json:"listen"`
	Logs   Logs     `json:"logs"`
	Routes []Route  `json:"routes"`
}

type Servers struct {
	Srv0 Server `json:"srv0"`
}

type Http struct {
	GracePeriod int64             `json:"grace_period"`
	Servers     map[string]Server `json:"servers"`
}

type DnsProvider struct {
	AuthToken string `json:"auth_token"`
	Name      string `json:"name"`
}

type Dns struct {
	Provider DnsProvider `json:"provider"`
}

type Challenges struct {
	Dns Dns `json:"dns"`
}

type Issuer struct {
	Challenges Challenges `json:"challenges"`
	Email      string     `json:"email"`
	Module     string     `json:"module"`
}

type Policy struct {
	Issuers  []Issuer `json:"issuers"`
	Subjects []string `json:"subjects"`
}

type Automation struct {
	Policies []Policy `json:"policies"`
}

type Tls struct {
	Automation Automation `json:"automation"`
}

type Encoder struct {
	Format string `json:"format"`
}

type Writer struct {
	Output string `json:"output"`
}

type Log struct {
	Encoder Encoder  `json:"encoder"`
	Exclude []string `json:"exclude,omitempty"`
	Include []string `json:"include,omitempty"`
	Writer  Writer   `json:"writer"`
}

type Logging struct {
	Logs map[string]Log `json:"logs"`
}

type Apps struct {
	Http    Http    `json:"http"`
	Tls     Tls     `json:"tls"`
	Logging Logging `json:"logging"`
}

type Config struct {
	Apps Apps `json:"apps"`
}

func NewCaddyV3(apiUrl string) *CaddyV3 {
	if !strings.HasSuffix("/config/", apiUrl) {
		apiUrl += "/config/"
	}

	if !strings.HasSuffix(apiUrl, "/") {
		apiUrl += "/"
	}

	return &CaddyV3{apiUrl: apiUrl}
}

func (c *CaddyV3) Get() ([]string, error) {
	resp, err := http.DefaultClient.Get(c.apiUrl)
	if err != nil {
		return []string{}, err
	}

	defer resp.Body.Close()
	var config Config

	bodyContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	err = json.Unmarshal(bodyContent, &config)
	if err != nil {
		return []string{}, err
	}

	var result []string

	for _, s := range config.Apps.Http.Servers {
		for _, r := range s.Routes {
			result = append(result, extractFromRoute(r)...)
		}
	}

	return result, nil
}

func extractFromRoute(r Route) []string {
	result := make([]string, 0)

	for _, m := range r.Match {
		if m.Host != nil {
			result = append(result, m.Host...)
		}

		for _, hi := range r.Handle {
			result = append(result, extractFromHandle(hi)...)
		}
	}

	return result
}

func extractFromHandle(h Handle) []string {
	var result []string
	for _, r := range h.Routes {
		result = append(result, extractFromRoute(r)...)
	}

	return result
}
