package uid

type UniqueIDService interface {
	NewID() uint64
}
