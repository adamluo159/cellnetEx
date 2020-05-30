package peer

import (
	"github.com/adamluo159/cellnetEx"
	"sync"
	"sync/atomic"
)

// 完整功能的会话管理
type SessionManager interface {
	cellnetEx.SessionAccessor

	Add(cellnetEx.Session)
	Remove(cellnetEx.Session)
	Count() int

	// 设置ID开始的号
	SetIDBase(base int64)
}

type CoreSessionManager struct {
	sesById sync.Map // 使用Id关联会话

	sesIDGen int64 // 记录已经生成的会话ID流水号

	count int64 // 记录当前在使用的会话数量
}

func (self *CoreSessionManager) SetIDBase(base int64) {

	atomic.StoreInt64(&self.sesIDGen, base)
}

func (self *CoreSessionManager) Count() int {
	return int(atomic.LoadInt64(&self.count))
}

func (self *CoreSessionManager) Add(ses cellnetEx.Session) {

	id := atomic.AddInt64(&self.sesIDGen, 1)

	atomic.AddInt64(&self.count, 1)

	ses.(interface {
		SetID(int64)
	}).SetID(id)

	self.sesById.Store(id, ses)
}

func (self *CoreSessionManager) Remove(ses cellnetEx.Session) {

	self.sesById.Delete(ses.ID())

	atomic.AddInt64(&self.count, -1)
}

// 获得一个连接
func (self *CoreSessionManager) GetSession(id int64) cellnetEx.Session {
	if v, ok := self.sesById.Load(id); ok {
		return v.(cellnetEx.Session)
	}

	return nil
}

func (self *CoreSessionManager) VisitSession(callback func(cellnetEx.Session) bool) {

	self.sesById.Range(func(key, value interface{}) bool {

		return callback(value.(cellnetEx.Session))

	})
}

func (self *CoreSessionManager) CloseAllSession() {

	self.VisitSession(func(ses cellnetEx.Session) bool {

		ses.Close()

		return true
	})
}

// 活跃的会话数量
func (self *CoreSessionManager) SessionCount() int {

	v := atomic.LoadInt64(&self.count)

	return int(v)
}
