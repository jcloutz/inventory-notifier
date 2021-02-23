package inventory_notifier

import (
	"errors"
	"regexp"
	"sync"
)

type ProductMatcher struct {
	pattern  *regexp.Regexp
	name     string
	maxPrice float64
}

type MatcherContainer struct {
	matchers []ProductMatcher
	mutex    sync.RWMutex
}

func NewMatchContainer() *MatcherContainer {
	return &MatcherContainer{}
}

func (m *MatcherContainer) Add(productConfig *ProductConfig) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.matchers = append(m.matchers, ProductMatcher{
		pattern:  regexp.MustCompile("(?i)" + productConfig.Name),
		name:     productConfig.Name,
		maxPrice: productConfig.MaxPrice,
	})
}

func (m *MatcherContainer) Find(name string) (*ProductConfig, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, matcher := range m.matchers {
		if found := matcher.pattern.MatchString(name); found {
			return &ProductConfig{
				Name:     matcher.name,
				MaxPrice: matcher.maxPrice,
			}, nil
		}
	}

	return nil, errors.New("no pattern found for product")
}
