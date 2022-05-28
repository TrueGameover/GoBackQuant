package strategy

import "errors"

type Parameter struct {
	Name  string
	value int
	Min   int
	Max   int
}

func (p *Parameter) SetValue(value int) error {
	if p.Min <= value && p.Max >= value {
		p.value = value
		return nil
	}

	return errors.New("invalid value")
}

func (p *Parameter) GetValue() int {
	return p.value
}

func NewParameter(name string, min int, max int) Parameter {
	return Parameter{
		Name:  name,
		value: 0,
		Min:   min,
		Max:   max,
	}
}
