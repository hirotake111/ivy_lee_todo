package db

type Db struct{}

func NewDb() *Db {
	return new(Db)
}
