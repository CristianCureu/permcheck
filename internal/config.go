package internal

type Config struct {
	InsecureTag   string
	SecureTag     string
	NumWorkers    int
	Benchmark     bool
	InsecureOnly  bool
	FixMode       bool
	IsTerminal    bool
	ColorsEnabled bool
}
