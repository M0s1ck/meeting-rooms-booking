package login

type passComparor interface {
	Compare(hash, password string) error
}
