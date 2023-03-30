package uid

import (
	"sync"
	"time"

	"github.com/sony/sonyflake"
)

var (
	startTime = time.Unix(1658709116, 0)
	once      sync.Once

	// singleton
	uniqueIDGenerator *UniqueIDGenerator
)

type UniqueIDGenerator struct {
	sf *sonyflake.Sonyflake
}

var _ UniqueIDService = &UniqueIDGenerator{}

func NewUniqueIDGenerator() *UniqueIDGenerator {
	once.Do(func() {
		var st sonyflake.Settings
		st.StartTime = startTime
		uniqueIDGenerator = new(UniqueIDGenerator)
		uniqueIDGenerator.sf = sonyflake.NewSonyflake(st)
		if uniqueIDGenerator.sf == nil {
			panic("sonyflake not created")
		}
	})

	return uniqueIDGenerator
}

func (u *UniqueIDGenerator) NewID() uint64 {
	id, _ := u.sf.NextID()
	return id
}
