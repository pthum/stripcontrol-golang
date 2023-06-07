package database

type DBReader[T any] interface {
	GetAll() ([]T, error)
	Get(id string) (*T, error)
	Close()
}

type DBWriter[T any] interface {
	Save(input *T) (err error)
	Update(dbObject T, input T) (err error)
	Create(input *T) (err error)
	Delete(input *T) (err error)
	Close()
}

//go:generate mockery --name=DBHandler --with-expecter=true
type DBHandler[T any] interface {
	DBReader[T]
	DBWriter[T]
}
