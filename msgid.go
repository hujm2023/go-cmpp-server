package cmppserver

import (
	"sync/atomic"
	"time"

	"github.com/bytedance/gopkg/lang/fastrand"
	"github.com/hujm2023/go-sms-protocol/cmpp"
)

var (
	msgID  uint64
	gateID = fastrand.Uint64n(1024)
)

func GenMsgID() uint64 {
	now := time.Now()
	sequenceID := atomic.AddUint64(&msgID, 1)
	return cmpp.CombineMsgID(uint64(now.Month()), uint64(now.Day()), uint64(now.Hour()), uint64(now.Minute()), uint64(now.Day()), gateID, sequenceID)
}
