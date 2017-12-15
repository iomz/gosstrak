package filter

// Map contains filter as key and notify string as value
type Map map[string]string

func (fm Map) keys() []string {
	ks := []string{}
	for k := range fm {
		ks = append(ks, k)
	}
	return ks
}
