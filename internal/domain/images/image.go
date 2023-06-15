package images

type ImageRepo interface {
	DeleteImage(id uint) error
}
