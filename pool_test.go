package mempool

import (
	"bytes"
	"math/rand"
	"runtime"
	"runtime/debug"
	"testing"
)

var PoolSizes = []int{
	1024 * 1024,
	512 * 1024,
	256 * 1024,
	128 * 1024,
	64 * 1024,
	32 * 1024,
	16 * 1024,
	10 * 1024,
	8 * 1024,
	4 * 1024,
	2 * 1024,
	1 * 1024,
	512,
	256,
	128,
	64,
	32,
}

func TestPool(t *testing.T) {

	// 需要测试各个尺寸的池，确保不存在错误尺寸 panic 的问题
	for _, size := range PoolSizes {
		if size < 1024*1024 {
			testPoolSize(t, size+1)
		}

		testPoolSize(t, size)

		testPoolSize(t, size-1)
	}

	m := Get(1024*1024 + 1)
	if len(m) != (1024*1024 + 1) {
		t.Errorf("!=")
	}

}

func testPoolSize(t *testing.T, size int) {
	// 禁用垃圾回收
	debug.SetGCPercent(-1)
	defer func() {
		debug.SetGCPercent(100)
		runtime.GC()
	}()

	data := make([]byte, size)
	rand.Read(data)

	// ----------- 确保申请的内存尺寸正确
	m1 := Get(size)
	if s := len(m1); s != size {
		t.Errorf("%v != %v", s, size)
	}
	copy(m1, data)

	// ------- 测试池的确在起作用
	err := Put(m1)
	if err != nil {
		t.Errorf("put,%v", err)
	}

	m2 := Get(size)
	if bytes.Equal(m1, m2) == false {
		if isZero(m2) == true {
			t.Logf("size:%v m2==0", size)
		} else {
			t.Errorf("%#v\r\n!=\r\n%#v", m1, m2)
		}
	}

	// ---------- 	确保池不会重复
	m3 := Get(size)
	if bytes.Equal(m1, m3) {
		t.Errorf("%#v\r\n==\r\n%#v", m1, m3)
	}

}

func isZero(b []byte) bool {
	for _, v := range b {
		if v != 0 {
			return false
		}
	}
	return true
}
