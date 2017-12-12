package filter

// FilterMap contains filter as key and notify string as value
type FilterMap map[string]string

func (fm FilterMap) keys() []string {
	ks := []string{}
	for k := range fm {
		ks = append(ks, k)
	}
	return ks
}
