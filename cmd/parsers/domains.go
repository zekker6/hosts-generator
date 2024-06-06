package parsers

// FilterWildcardEntries Remove entries which start with a wildcard
func FilterWildcardEntries(domains []string) []string {
	var filtered []string

	for _, domain := range domains {
		if domain[0] == '*' {
			continue
		}
		filtered = append(filtered, domain)
	}

	return filtered
}
