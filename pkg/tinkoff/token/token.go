package token

type TinkoffToken struct {
	Token   string
	Sandbox bool
}

func (t *TinkoffToken) GetToken() string {
	return t.Token
}

func (t *TinkoffToken) IsSandbox() bool {
	return t.Sandbox
}
