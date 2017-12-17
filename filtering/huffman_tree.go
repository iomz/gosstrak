package filtering

type HuffmanTree struct {
	Notify       string
	Filter       *Filter
	MatchNext    *HuffmanTree
	MismatchNext *HuffmanTree
}
