package main

import (
   "bufio"
   "fmt"
   "os"
)

// Segment tree for binary values with range assign and range sum
type Seg struct {
   n    int
   sum  []int
   lazy []int // -1: none, 0 or 1 for assign
}
// build initializes the tree with array a of 0/1 values
func (s *Seg) build(o, l, r int, a []int) {
   s.lazy[o] = -1
   if l == r {
       s.sum[o] = a[l]
       return
   }
   m := (l + r) >> 1
   s.build(o*2, l, m, a)
   s.build(o*2+1, m+1, r, a)
   s.sum[o] = s.sum[o*2] + s.sum[o*2+1]
}

func NewSeg(n int) *Seg {
   size := 4 * n
   s := &Seg{n: n, sum: make([]int, size), lazy: make([]int, size)}
   for i := range s.lazy {
       s.lazy[i] = -1
   }
   return s
}

func (s *Seg) push(o, l, r int) {
   if s.lazy[o] != -1 {
       v := s.lazy[o]
       s.sum[o] = (r - l + 1) * v
       if l < r {
           s.lazy[o*2] = v
           s.lazy[o*2+1] = v
       }
       s.lazy[o] = -1
   }
}

func (s *Seg) update(o, l, r, ql, qr, v int) {
   s.push(o, l, r)
   if qr < l || r < ql {
       return
   }
   if ql <= l && r <= qr {
       s.lazy[o] = v
       s.push(o, l, r)
       return
   }
   m := (l + r) >> 1
   s.update(o*2, l, m, ql, qr, v)
   s.update(o*2+1, m+1, r, ql, qr, v)
   s.sum[o] = s.sum[o*2] + s.sum[o*2+1]
}

func (s *Seg) query(o, l, r, ql, qr int) int {
   s.push(o, l, r)
   if qr < l || r < ql {
       return 0
   }
   if ql <= l && r <= qr {
       return s.sum[o]
   }
   m := (l + r) >> 1
   return s.query(o*2, l, m, ql, qr) + s.query(o*2+1, m+1, r, ql, qr)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   fmt.Fscan(in, &n, &m)
   var s0 string
   fmt.Fscan(in, &s0)
   sBytes := []byte(s0)
   // build segment trees for each letter
   segs := make([]*Seg, 26)
   // build segment trees for each letter
   for c := 0; c < 26; c++ {
       segs[c] = NewSeg(n)
       // initial array for letter c
       a := make([]int, n)
       for i := 0; i < n; i++ {
           if sBytes[i] == byte('a'+c) {
               a[i] = 1
           }
       }
       segs[c].build(1, 0, n-1, a)
   }
   // process queries
   for qi := 0; qi < m; qi++ {
       var l, r int
       fmt.Fscan(in, &l, &r)
       l--;
       r--;
       // get counts
       cnt := [26]int{}
       oddC := 0
       oddIdx := -1
       for c := 0; c < 26; c++ {
           cnt[c] = segs[c].query(1, 0, n-1, l, r)
           if cnt[c]&1 == 1 {
               oddC++
               oddIdx = c
           }
       }
       if oddC > 1 {
           continue
       }
       // clear range
       for c := 0; c < 26; c++ {
           if cnt[c] > 0 {
               segs[c].update(1, 0, n-1, l, r, 0)
           }
       }
       left, right := l, r
       // place halves
       for c := 0; c < 26; c++ {
           h := cnt[c] / 2
           if h > 0 {
               segs[c].update(1, 0, n-1, left, left+h-1, 1)
               segs[c].update(1, 0, n-1, right-h+1, right, 1)
               left += h
               right -= h
           }
       }
       // place middle
       if oddC == 1 {
           segs[oddIdx].update(1, 0, n-1, left, left, 1)
       }
   }
   // build result
   res := make([]byte, n)
   for i := 0; i < n; i++ {
       for c := 0; c < 26; c++ {
           if segs[c].query(1, 0, n-1, i, i) == 1 {
               res[i] = byte('a' + c)
               break
           }
       }
   }
   out.Write(res)
}
