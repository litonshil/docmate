package llm

import (
	"fmt"
)

type Factory struct {
	providers map[string]Provider
}

func NewFactory() *Factory {
	f := &Factory{
		providers: make(map[string]Provider),
	}
	f.Register(NewGeminiProvider(""))

	return f
}

func (f *Factory) Register(p Provider) {
	f.providers[p.GetName()] = p
}

func (f *Factory) GetProvider(name string) (Provider, error) {
	p, ok := f.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", name)
	}

	return p, nil
}
