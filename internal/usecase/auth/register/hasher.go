package register

type passHasher interface {
	Hash(password string) (string, error)
}
