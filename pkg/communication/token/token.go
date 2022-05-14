package token

type Token interface {
	GetToken() string
	IsSandbox() bool
}
