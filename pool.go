package gomempool

import (
	"fmt"
	"runtime/debug"
	"sync"
)

/*
内存池
*/

var PoolByte1m = NewPool(1024 * 1024)
var PoolByte512k = NewPool(512 * 1024)
var PoolByte256k = NewPool(256 * 1024)
var PoolByte128k = NewPool(128 * 1024)
var PoolByte64k = NewPool(64 * 1024)
var PoolByte32k = NewPool(32 * 1024)
var PoolByte16k = NewPool(16 * 1024)
var PoolByte10k = NewPool(10 * 1024)
var PoolByte8k = NewPool(8 * 1024)
var PoolByte4k = NewPool(4 * 1024)
var PoolByte2k = NewPool(2 * 1024)
var PoolByte1k = NewPool(1 * 1024)
var PoolByte512 = NewPool(512)
var PoolByte256 = NewPool(256)
var PoolByte128 = NewPool(128)
var PoolByte64 = NewPool(64)
var PoolByte32 = NewPool(32)

// 这个封装的目的是防止返回的切片尺寸不对
type Pool struct {
	pool *sync.Pool
	size int
}

type CapErr struct {
	Buf   []byte
	Cap   int
	Stack []byte
}

// 新建一个池
func NewPool(size int) *Pool {
	p := Pool{}
	p.size = size
	p.pool = &sync.Pool{New: func() interface{} {
		return make([]byte, size)
	}}

	return &p
}

func (p *Pool) Get() []byte {
	// sync.pool 是线程安全的
	return p.pool.Get().([]byte)
}

// 回填池
//可以放心尺寸的问题，切片尺寸有问题会被放弃
func (p *Pool) Put(buf []byte) {
	if cap(buf) == p.size {
		p.pool.Put(buf[:p.size])
	}
}

// 自动选择合适的池
func Get(l int) []byte {
	if l > 1024*1024 {
		return make([]byte, l)
	} else if l > 512*1024 {
		return PoolByte1m.Get()[:l]
	} else if l > 256*1024 {
		return PoolByte512k.Get()[:l]
	} else if l > 128*1024 {
		return PoolByte256k.Get()[:l]
	} else if l > 64*1024 {
		return PoolByte128k.Get()[:l]
	} else if l > 32*1024 {
		return PoolByte64k.Get()[:l]
	} else if l > 16*1024 {
		return PoolByte32k.Get()[:l]
	} else if l > 10*1024 {
		return PoolByte16k.Get()[:l]
	} else if l > 8*1024 {
		return PoolByte10k.Get()[:l]
	} else if l > 4*1024 {
		return PoolByte8k.Get()[:l]
	} else if l > 2*1024 {
		return PoolByte4k.Get()[:l]
	} else if l > 1*1024 {
		return PoolByte2k.Get()[:l]
	} else if l > 512 {
		return PoolByte1k.Get()[:l]
	} else if l > 256 {
		return PoolByte512.Get()[:l]
	} else if l > 128 {
		return PoolByte256.Get()[:l]
	} else if l > 64 {
		return PoolByte128.Get()[:l]
	} else if l > 32 {
		return PoolByte64.Get()[:l]
	} else {
		return PoolByte32.Get()[:l]
	}
}

// 回填池
//可以放心尺寸的问题，切片尺寸有问题会被放弃
func Put(buf []byte) error {
	l := cap(buf)
	if l > 1024*1024 {
		return nil
	}
	switch l {
	case 1024 * 1024:
		PoolByte1m.Put(buf)
	case 512 * 1024:
		PoolByte512k.Put(buf)
	case 256 * 1024:
		PoolByte256k.Put(buf)
	case 128 * 1024:
		PoolByte128k.Put(buf)
	case 64 * 1024:
		PoolByte64k.Put(buf)
	case 32 * 1024:
		PoolByte32k.Put(buf)
	case 16 * 1024:
		PoolByte16k.Put(buf)
	case 10 * 1024:
		PoolByte10k.Put(buf)
	case 8 * 1024:
		PoolByte8k.Put(buf)
	case 4 * 1024:
		PoolByte4k.Put(buf)
	case 2 * 1024:
		PoolByte2k.Put(buf)
	case 1 * 1024:
		PoolByte1k.Put(buf)
	case 512:
		PoolByte512.Put(buf)
	case 256:
		PoolByte256.Put(buf)
	case 128:
		PoolByte128.Put(buf)
	case 64:
		PoolByte64.Put(buf)
	case 32:
		PoolByte32.Put(buf)
	case 0:
		return nil
	default:
		err := &CapErr{
			Buf:   buf,
			Cap:   cap(buf),
			Stack: debug.Stack(),
		}
		return err
	}
	return nil
}

func (e *CapErr) Error() string {
	return fmt.Sprintf("gomempool.Put，cap(buf)=%v", e.Cap)
}
