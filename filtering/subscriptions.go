package filtering

// Subscriptions contains filter string as key and Info as value
type Subscriptions map[string]*Info

// Info contains notificationURI and pValue for a filter
type Info struct {
	NotificationURI string
	EntropyValue    float64
}

func (sub Subscriptions) keys() []string {
	ks := []string{}
	for k := range sub {
		ks = append(ks, k)
	}
	return ks
}
