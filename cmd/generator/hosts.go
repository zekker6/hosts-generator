package generator

import (
	"fmt"
)

func GenerateStrings(hosts []string, ip string, lineEnding string) string {
	result := ""

	for _, host := range hosts {
		result += fmt.Sprintf("%s\t%s%s", ip, host, lineEnding)
	}

	return result
}
