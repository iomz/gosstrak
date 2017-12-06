package filter

type FilterMap map[string]string

func (fm FilterMap) keys() []string {
	ks := []string{}
	for k, _ := range fm {
		ks = append(ks, k)
	}
	return ks
}
