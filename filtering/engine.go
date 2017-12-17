package filtering

type Engine interface {
	AnalyzeLocality(id []byte, prefix string, lm *LocalityMap)
	Dump() string
	MarshalBinary() ([]byte, error)
	Search(id []byte) []string
	UnmarshalBinary(data []byte) error
}
