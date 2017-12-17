package filtering

type HuffmanTree struct {
	Notify       string
	Filter       *FilterObject
	MatchNext    *HuffmanTree
	MismatchNext *HuffmanTree
}
