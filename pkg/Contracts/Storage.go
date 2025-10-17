package Contracts

type Storage interface {
	GetByPrimary(p []byte) (r []byte, err error)
	Insert(key []byte, data interface{}) (affRows int, err error)
}
