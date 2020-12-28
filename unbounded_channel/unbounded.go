package unbounded

import "sync"

type Unbounded struct {
	// 通用接口, 在类型断言上存在性能开销
	c chan interface{}
	sync.Mutex
	backlog []interface{}
}

func NewUnbounded() *Unbounded {
	return &Unbounded{
		c: make(chan interface{}, 1),
	}
}

/*
 *@func: add t to the unbounded buffer
 */
func (b *Unbounded) Put(t interface{}) {
	b.Lock()

	if len(b.backlog) == 0 {
		select {
		case b.c <- t:
			b.Unlock()
			return
		default:
		}
	}

	// 不阻塞的缓存
	b.backlog = append(b.backlog, t)
	b.Unlock()
}

/*
 *@func: send earliest buffered data, read channel by Get()
 */
func (b *Unbounded) Load() {
	b.Lock()

	if len(b.backlog) > 0 {
		select {
		case b.c <- b.backlog[0]:
			b.backlog[0] = nil
			b.backlog = b.backlog[1:]
		default:

		}
	}

	b.Unlock()
}

/*
 *@module: return a read channel on which values added to the buffer, prev call Load()
 */
func (b *Unbounded) Get() <-chan interface{} {
	return b.c
}
