package inventory_notifier

type ProductNotification struct {
	Name      string
	Url       string
	Site      string
	SalePrice float64
	MaxPrice  float64
}

type Notifier interface {
	Notify(product *ProductNotification)
}

type Notifiers struct {
	Items []Notifier
}

func (n *Notifiers) Notify(product *ProductNotification) {
	for _, notifier := range n.Items {
		notifier.Notify(product)
	}
}

func (n *Notifiers) Add(notifier Notifier) {
	n.Items = append(n.Items, notifier)
}
