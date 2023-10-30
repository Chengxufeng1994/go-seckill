package entity

type Repository interface {
	GetActivityList() ([]*Activity, error)
	CreateActivity(*Activity) error

	GetProductList() ([]*Product, error)
	CreateProduct(*Product) error
}
