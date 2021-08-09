package parsers


type Parser interface {
	Get() ([]string, error)
}
