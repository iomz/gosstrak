package filtering

type Engine interface {
	//AnalyzeLocality()
	Build(fm Map) *Engine
	Dump() string
	MarshalBinary() ([]byte, error)
	Search(id []byte) []string
	UnmarshalBinary(data []byte) error
}
