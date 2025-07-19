package main

import (
   "bufio"
   "fmt"
   "os"
)

// N holds interval [l, r] and index i
type N struct {
   l, r int64
   i    int
}

// reverseDeque supports back operations with a flipped orientation
type reverseDeque struct {
   org int64
   fw  bool
   d   []N
}

func newReverseDeque() *reverseDeque {
   return &reverseDeque{org: 0, fw: true, d: make([]N, 0)}
}

func (rd *reverseDeque) empty() bool {
   return len(rd.d) == 0
}

func (rd *reverseDeque) clear() {
   rd.d = rd.d[:0]
}

// r2i converts stored N to absolute coordinates
func (rd *reverseDeque) r2i(x N) N {
   if rd.fw {
       x.l += rd.org
       x.r += rd.org
   } else {
       tmp := x.l
       x.l = -x.r + rd.org
       x.r = -tmp + rd.org
   }
   return x
}

// back returns the last element in absolute coordinates
func (rd *reverseDeque) back() N {
   var x N
   if rd.fw {
       x = rd.d[len(rd.d)-1]
   } else {
       x = rd.d[0]
   }
   return rd.r2i(x)
}

// popBack removes the last element
func (rd *reverseDeque) popBack() {
   if rd.fw {
       rd.d = rd.d[:len(rd.d)-1]
   } else {
       rd.d = rd.d[1:]
   }
}

// pushBack appends x (absolute coordinates) into deque
func (rd *reverseDeque) pushBack(x N) {
   rd.d = append(rd.d, x)
}

// rev applies reverse operation with shift v
func (rd *reverseDeque) rev(v int64) {
   // shift(-v)
   if rd.fw {
       rd.org += v
   } else {
       rd.org -= v
   }
   rd.fw = !rd.fw
}

func solve() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   var C int64
   if _, err := fmt.Fscan(in, &n, &C); err != nil {
       return
   }
   w := make([]int64, n-1)
   for i := 0; i < n-2; i++ {
       fmt.Fscan(in, &w[i])
   }
   w[n-2] = C
   d := newReverseDeque()
   // g and pre arrays
   g := make([]int64, n-1)
   pre := make([]int, n-1)
   // init deque
   d.pushBack(N{l: 0, r: C, i: -1})
   for i := 0; i < n-1; i++ {
       v := w[i]
       for !d.empty() {
           x := d.back()
           if x.r <= v {
               break
           }
           d.popBack()
           if x.l <= v {
               x.r = v
               d.pushBack(x)
           }
       }
       if d.empty() {
           fmt.Fprintln(out, "NO")
           return
       }
       x := d.back()
       g[i] = x.r
       pre[i] = x.i
       // reverse with v
       d.rev(v)
       if g[i] == v {
           d.clear()
           d.pushBack(N{l: 0, r: v, i: i})
       } else {
           d.pushBack(N{l: v, r: v, i: i})
       }
   }
   // compute dif and mx
   dif := make([]int64, n-1)
   mx := make([]int, n-2)
   for pos := n - 2; pos >= 0; pos = pre[pos] {
       nx := pre[pos]
       if nx >= 0 {
           mx[nx] = 1
       }
       dif[pos] = g[pos]
       for i := pos - 1; i > nx; i-- {
           dif[i] = w[i] - dif[i+1]
       }
   }
   // build h
   h := make([]int64, n)
   if n > 1 {
       h[1] = dif[0]
   }
   for i := 0; i < n-2; i++ {
       a := h[i] < h[i+1]
       if mx[i] == 1 {
           a = !a
       }
       if a {
           h[i+2] = h[i+1] + dif[i+1]
       } else {
           h[i+2] = h[i+1] - dif[i+1]
       }
   }
   // offset to make min element zero or more
   mn := h[0]
   for _, v := range h {
       if v < mn {
           mn = v
       }
   }
   off := -mn
   fmt.Fprintln(out, "YES")
   for i, v := range h {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, v+off)
   }
   fmt.Fprintln(out)
}

func main() {
   solve()
}
