package login

//go:generate mockgen -source=pass_comparor.go -destination=mocks/pass_comparor_mock.go -package=mocks
type passComparor interface {
	Compare(hash, password string) error
}
