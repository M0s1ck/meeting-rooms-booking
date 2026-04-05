package register

//go:generate mockgen -source=hasher.go -destination=mocks/hasher_mock.go -package=mocks
type passHasher interface {
	Hash(password string) (string, error)
}
