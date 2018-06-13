package tdt

// PartitionTableKey is used for PartitionTables
type PartitionTableKey int

// PartitionTable is used to get the related values for each coding scheme
type PartitionTable map[int]map[PartitionTableKey]int

// Key values for PartitionTables
const (
	PValue PartitionTableKey = iota
	CPBits
	IRBits
	IRDigits
	EBits
	EDigits
	ATBits
	ATDigits
	IARBits
	IARDigits
)

// GIAI96PartitionTable is PT for GIAI
var GIAI96PartitionTable = PartitionTable{
	12: {PValue: 0, CPBits: 40, IARBits: 42, IARDigits: 13},
	11: {PValue: 1, CPBits: 37, IARBits: 45, IARDigits: 14},
	10: {PValue: 2, CPBits: 34, IARBits: 48, IARDigits: 15},
	9:  {PValue: 3, CPBits: 30, IARBits: 52, IARDigits: 16},
	8:  {PValue: 4, CPBits: 27, IARBits: 55, IARDigits: 17},
	7:  {PValue: 5, CPBits: 24, IARBits: 58, IARDigits: 18},
	6:  {PValue: 6, CPBits: 20, IARBits: 62, IARDigits: 19},
}

// GRAI96PartitionTable is PT for GRAI
var GRAI96PartitionTable = PartitionTable{
	12: {PValue: 0, CPBits: 40, ATBits: 4, ATDigits: 0},
	11: {PValue: 1, CPBits: 37, ATBits: 7, ATDigits: 1},
	10: {PValue: 2, CPBits: 34, ATBits: 10, ATDigits: 2},
	9:  {PValue: 3, CPBits: 30, ATBits: 14, ATDigits: 3},
	8:  {PValue: 4, CPBits: 27, ATBits: 17, ATDigits: 4},
	7:  {PValue: 5, CPBits: 24, ATBits: 20, ATDigits: 5},
	6:  {PValue: 6, CPBits: 20, ATBits: 24, ATDigits: 6},
}

// SGTIN96PartitionTable is PT for SGTIN
var SGTIN96PartitionTable = PartitionTable{
	12: {PValue: 0, CPBits: 40, IRBits: 4, IRDigits: 1},
	11: {PValue: 1, CPBits: 37, IRBits: 7, IRDigits: 2},
	10: {PValue: 2, CPBits: 34, IRBits: 10, IRDigits: 3},
	9:  {PValue: 3, CPBits: 30, IRBits: 14, IRDigits: 4},
	8:  {PValue: 4, CPBits: 27, IRBits: 17, IRDigits: 5},
	7:  {PValue: 5, CPBits: 24, IRBits: 20, IRDigits: 6},
	6:  {PValue: 6, CPBits: 20, IRBits: 24, IRDigits: 7},
}

// SSCC96PartitionTable is PT for SSCC
var SSCC96PartitionTable = PartitionTable{
	12: {PValue: 0, CPBits: 40, EBits: 18, EDigits: 5},
	11: {PValue: 1, CPBits: 37, EBits: 21, EDigits: 6},
	10: {PValue: 2, CPBits: 34, EBits: 24, EDigits: 7},
	9:  {PValue: 3, CPBits: 30, EBits: 28, EDigits: 8},
	8:  {PValue: 4, CPBits: 27, EBits: 31, EDigits: 9},
	7:  {PValue: 5, CPBits: 24, EBits: 34, EDigits: 10},
	6:  {PValue: 6, CPBits: 20, EBits: 38, EDigits: 11},
}
