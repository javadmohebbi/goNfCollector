package reputation

import (
	"errors"

	"github.com/goNfCollector/debugger"
)

type ReputationResponse struct {
	Result  map[string]interface{} `json:"result"`
	Current uint                   `json:"current"`
}

type ReputationType uint

const (
	TYPE_IPSum ReputationType = iota
)

func (r ReputationType) String() string {
	return [...]string{
		"IPSum",
	}[r]
}

type Reputation interface {
	Get(string) ReputationResponse

	GetType() string
}

//  new Reputation
func New(reput interface{}, d *debugger.Debugger) (*Reputation, error) {

	switch reput.(type) {
	case IPSum:
		repu := reput.(IPSum)
		r := Reputation(&repu)

		return &r, nil
	}

	return nil, errors.New("No valid Reputations have found!")
}
