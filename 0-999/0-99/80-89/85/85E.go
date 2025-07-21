package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func abs(a int) int {
   if a < 0 {
       return -a
   }
   return a
}

func pow2(exp int) int {
   res := 1
   base := 2
   for exp > 0 {
       if exp & 1 == 1 {
           res = int((int64(res) * int64(base)) % MOD)
       }
       base = int((int64(base) * int64(base)) % MOD)
       exp >>= 1
   }
   return res
}

func check(n int, x, y []int, D int, needCount bool) (bool, int) {
   color := make([]int8, n)
   for i := range color {
       color[i] = -1
   }
   q := make([]int, n)
   comps := 0
   for i := 0; i < n; i++ {
       if color[i] != -1 {
           continue
       }
       // new component
       comps++
       head, tail := 0, 0
       color[i] = 0
       q[tail] = i; tail++
       for head < tail {
           u := q[head]; head++
           xu, yu := x[u], y[u]
           cu := color[u]
           for v := 0; v < n; v++ {
               if u == v {
                   continue
               }
               dx := xu - x[v]
               if dx < 0 {
                   dx = -dx
               }
               if dx <= D {
                   dy := yu - y[v]
                   if dy < 0 {
                       dy = -dy
                   }
                   if dx + dy <= D {
                       continue
                   }
               }
               // edge exists between u and v in conflict graph
               if color[v] == -1 {
                   color[v] = 1 - cu
                   q[tail] = v; tail++
               } else if color[v] == cu {
                   return false, 0
               }
           }
       }
   }
   if needCount {
       return true, comps
   }
   return true, 0
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   // read n
   var n int
   fmt.Fscan(in, &n)
   x := make([]int, n)
   y := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &x[i], &y[i])
   }
   // binary search D
   lo, hi := -1, 10001
   for lo+1 < hi {
       mid := (lo + hi) >> 1
       ok, _ := check(n, x, y, mid, false)
       if ok {
           hi = mid
       } else {
           lo = mid
       }
   }
   D := hi
   _, comps := check(n, x, y, D, true)
   cnt := pow2(comps)
   fmt.Fprintln(out, D)
   fmt.Fprintln(out, cnt)
}
