package main

import (
	"bufio"
	"fmt"
	"os"
)

type Deque struct {
	data []int
	head int
	size int
}

func NewDeque() *Deque {
	return &Deque{data: make([]int, 1)}
}

func (d *Deque) ensureCapacity() {
	if d.size < len(d.data) {
		return
	}
	newCap := len(d.data) * 2
	if newCap == 0 {
		newCap = 1
	}
	newData := make([]int, newCap)
	for i := 0; i < d.size; i++ {
		newData[i] = d.data[(d.head+i)%len(d.data)]
	}
	d.data = newData
	d.head = 0
}

func (d *Deque) PushFront(val int) {
	d.ensureCapacity()
	d.head = (d.head - 1 + len(d.data)) % len(d.data)
	d.data[d.head] = val
	d.size++
}

func (d *Deque) PushBack(val int) {
	d.ensureCapacity()
	idx := (d.head + d.size) % len(d.data)
	d.data[idx] = val
	d.size++
}

func (d *Deque) PopFront() int {
	val := d.data[d.head]
	d.head = (d.head + 1) % len(d.data)
	d.size--
	return val
}

func (d *Deque) PopBack() int {
	idx := (d.head + d.size - 1) % len(d.data)
	val := d.data[idx]
	d.size--
	return val
}

func (d *Deque) Len() int {
	return d.size
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var q int
		fmt.Fscan(in, &q)
		dq := NewDeque()
		var length int64
		var sum int64
		var riz int64
		rev := false

		for i := 0; i < q; i++ {
			var s int
			fmt.Fscan(in, &s)
			if s == 3 {
				var k int
				fmt.Fscan(in, &k)
				if !rev {
					dq.PushBack(k)
				} else {
					dq.PushFront(k)
				}
				length++
				sum += int64(k)
				riz += int64(k) * length
			} else if s == 1 {
				if length == 0 {
					fmt.Fprintln(out, 0)
					continue
				}
				var last int
				if !rev {
					last = dq.PopBack()
					dq.PushFront(last)
				} else {
					last = dq.PopFront()
					dq.PushBack(last)
				}
				riz += sum - int64(last)*length
			} else if s == 2 {
				riz = (length+1)*sum - riz
				rev = !rev
			}
			fmt.Fprintln(out, riz)
		}
	}
}
