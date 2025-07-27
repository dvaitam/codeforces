package main

import (
   "bufio"
   "fmt"
   "os"
)

// Fenwick tree for point update and prefix sum
type Fenw struct {
   n int
   t []int64
}

func NewFenw(n int) *Fenw {
   return &Fenw{n: n, t: make([]int64, n+1)}
}

func (f *Fenw) Update(i int, v int64) {
   for x := i; x <= f.n; x += x & -x {
       f.t[x] += v
   }
}

func (f *Fenw) Query(i int) int64 {
   var s int64
   for x := i; x > 0; x -= x & -x {
       s += f.t[x]
   }
   return s
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, q int
   fmt.Fscan(in, &n, &q)
   N := 1 << n
   a := make([]int64, N)
   for i := 0; i < N; i++ {
       fmt.Fscan(in, &a[i])
   }
   bit := NewFenw(N)
   for i, v := range a {
       bit.Update(i+1, v)
   }
   var M int
   type iv struct{ l, r, add int }
   ivs := make([]iv, 0, 64)
   var dfs func(l0, r0, mask, b, add int)
   dfs = func(l0, r0, mask, b, add int) {
       if l0 > r0 {
           return
       }
       if mask == 0 || b < 0 {
           ivs = append(ivs, iv{add + l0, add + r0, 0})
           return
       }
       th := 1 << b
       low := mask & (th - 1)
       if (mask>>b)&1 == 0 {
           if r0 < th {
               dfs(l0, r0, low, b-1, add)
           } else if l0 >= th {
               dfs(l0-th, r0-th, low, b-1, add+th)
           } else {
               dfs(l0, th-1, low, b-1, add)
               dfs(0, r0-th, low, b-1, add+th)
           }
       } else {
           if r0 < th {
               dfs(l0, r0, low, b-1, add+th)
           } else if l0 >= th {
               dfs(l0-th, r0-th, low, b-1, add)
           } else {
               dfs(l0, th-1, low, b-1, add+th)
               dfs(0, r0-th, low, b-1, add)
           }
       }
   }
   for qi := 0; qi < q; qi++ {
       var tp int
       fmt.Fscan(in, &tp)
       switch tp {
       case 1:
           var x int
           var k int64
           fmt.Fscan(in, &x, &k)
           pos := (x - 1) ^ M
           delta := k - a[pos]
           a[pos] = k
           bit.Update(pos+1, delta)
       case 2:
           var k int
           fmt.Fscan(in, &k)
           if k > 0 {
               M ^= (1 << k) - 1
           }
       case 3:
           var k int
           fmt.Fscan(in, &k)
           M ^= 1 << k
       case 4:
           var l, r int
           fmt.Fscan(in, &l, &r)
           ivs = ivs[:0]
           dfs(l-1, r-1, M, n-1, 0)
           var ans int64
           for _, seg := range ivs {
               L := seg.l + 1
               R := seg.r + 1
               ans += bit.Query(R) - bit.Query(L-1)
           }
           fmt.Fprintln(out, ans)
       }
   }
}
