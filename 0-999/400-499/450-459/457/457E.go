package main

import (
   "bufio"
   "fmt"
   "os"
)

type DSU struct {
   p    []int
   diff []int64 // diff[x] = potential[x] - potential[p[x]]
}

func NewDSU(n int) *DSU {
   p := make([]int, n+1)
   diff := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       p[i] = i
       diff[i] = 0
   }
   return &DSU{p: p, diff: diff}
}

// find returns root of x and accumulated diff from x to root: potential[x] - potential[root]
func (d *DSU) find(x int) (int, int64) {
   if d.p[x] == x {
       return x, 0
   }
   r, dr := d.find(d.p[x])
   // d.diff[x] = potential[x] - potential[p[x]]; dr = potential[p[x]] - potential[r]
   d.diff[x] += dr
   d.p[x] = r
   return r, d.diff[x]
}

// unite adds constraint: potential[f] - potential[t] = delta
// returns false if conflict
func (d *DSU) unite(f, t int, delta int64) bool {
   rf, df := d.find(f)
   rt, dt := d.find(t)
   if rf == rt {
       // check df - dt == delta
       if df-dt != delta {
           return false
       }
       return true
   }
   // attach rf under rt: set d.p[rf]=rt
   // we need potential[rf] - potential[rt] = delta + dt - df
   d.p[rf] = rt
   d.diff[rf] = delta + dt - df
   return true
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   dsu := NewDSU(n)
   bad := 0
   for i := 1; i <= m; i++ {
       var f, t int
       var w, b int64
       fmt.Fscan(in, &f, &t, &w, &b)
       if bad != 0 {
           continue
       }
       delta := 2 * w * b
       if !dsu.unite(f, t, delta) {
           bad = i
       }
   }
   if bad != 0 {
       fmt.Printf("BAD %d\n", bad)
       return
   }
   // check if potentials of 1 and n are connected
   r1, d1 := dsu.find(1)
   rn, dn := dsu.find(n)
   if r1 != rn {
       fmt.Println("UNKNOWN")
   } else {
       // potential[1] - potential[n] = d1 - dn
       diff := d1 - dn
       // efficiency = (potential[1]-potential[n]) / 2, rounded to nearest int
       var ans int64
       if diff >= 0 {
           ans = (diff + 1) / 2
       } else {
           ans = (diff - 1) / 2
       }
       fmt.Println(ans)
   }
}
