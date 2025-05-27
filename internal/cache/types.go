package cache

import "fmt"

type Mode int

const (
	JSON Mode = iota
	SQLite
)

func (m *Mode) String() string {
	switch *m {
	case JSON:
		return "JSON"
	case SQLite:
		return "SQLite"
	default:
		return "SQLite"
	}
}

func (m *Mode) Set(s string) error {
	switch s {
	case "JSON":
		*m = JSON
	case "SQLite":
		*m = SQLite
	default:
		return fmt.Errorf("invalid mode %q", s)
	}
	return nil
}
