package main

import (
   "bufio"
   "fmt"
   "os"
)

// deque holds a sliding-window deque of (value, index) pairs
type deque struct {
   vals []int32
   idxs []int32
   head int
   len  int
   cap  int
}

// init initializes the deque with given capacity
func (d *deque) init(cap int) {
   d.vals = make([]int32, cap)
   d.idxs = make([]int32, cap)
   d.head = 0
   d.len = 0
   d.cap = cap
}

// clear resets the deque
func (d *deque) clear() {
   d.head = 0
   d.len = 0
}

// push adds (v, i) to back, removing back while >= v
func (d *deque) push(v int32, i int32) {
   // pop back
   for d.len > 0 {
       back := (d.head + d.len - 1) % d.cap
       if d.vals[back] < v {
           break
       }
       d.len--
   }
   pos := (d.head + d.len) % d.cap
   d.vals[pos] = v
   d.idxs[pos] = i
   d.len++
}

// popFront removes from front while idx <= threshold
func (d *deque) popFront(threshold int32) {
   for d.len > 0 && d.idxs[d.head] <= threshold {
       d.head = (d.head + 1) % d.cap
       d.len--
   }
}

// front returns the front value
func (d *deque) front() int32 {
   return d.vals[d.head]
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, a, b int
   fmt.Fscan(in, &n, &m, &a, &b)
   var g0, x, y, z int64
   fmt.Fscan(in, &g0, &x, &y, &z)
   // prepare deques
   M2 := m - b + 1
   colD := make([]deque, M2)
   for j := 0; j < M2; j++ {
       colD[j].init(a)
   }
   // row deque
   var rowD deque
   rowD.init(b)
   // buffer for current row
   h := make([]int32, m)
   // generate
   var gPrev = g0
   var sum int64
   for i := 0; i < n; i++ {
       // generate row of h values from g sequence
       for j := 0; j < m; j++ {
           h[j] = int32(gPrev)
           gPrev = (gPrev*x + y) % z
       }
       // sliding min on row of width b
       rowD.clear()
       // rowMin for this row
       for j := 0; j < m; j++ {
           rowD.push(h[j], int32(j))
           // pop front outside window
           rowD.popFront(int32(j - b))
           if j >= b-1 {
               v := rowD.front()
               // update corresponding column deque
               cd := &colD[j-(b-1)]
               cd.push(v, int32(i))
               cd.popFront(int32(i - a))
               if i >= a-1 {
                   sum += int64(cd.front())
               }
           }
       }
   }
   // output result
   fmt.Println(sum)
}
