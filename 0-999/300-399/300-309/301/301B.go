package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var d int64
   if _, err := fmt.Fscan(in, &n, &d); err != nil {
       return
   }
   a := make([]int64, n+1)
   for i := 2; i <= n-1; i++ {
       fmt.Fscan(in, &a[i])
   }
   x := make([]int64, n+1)
   y := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &x[i], &y[i])
   }
   // cost to go from i to j: d * (|xi-xj| + |yi-yj|)
   // maximum needed initial fuel is cost from 1 to n
   dist := abs(x[1]-x[n]) + abs(y[1]-y[n])
   high := d * dist
   low := int64(0)
   // binary search minimal initial fuel
   for low < high {
       mid := (low + high) / 2
       if canReach(n, d, a, x, y, mid) {
           high = mid
       } else {
           low = mid + 1
       }
   }
   // low is minimal fuel units
   fmt.Println(low)
}

func canReach(n int, d int64, a, x, y []int64, initFuel int64) bool {
   const INF = math.MaxInt64 / 4
   best := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       best[i] = -INF
   }
   best[1] = initFuel
   // relax edges up to n-1 times
   for it := 0; it < n-1; it++ {
       updated := false
       for i := 1; i <= n; i++ {
           if best[i] < 0 {
               continue
           }
           bi := best[i]
           for j := 1; j <= n; j++ {
               if i == j {
                   continue
               }
               cost := d * (abs(x[i]-x[j]) + abs(y[i]-y[j]))
               if bi < cost {
                   continue
               }
               nf := bi - cost + a[j]
               if nf > best[j] {
                   best[j] = nf
                   updated = true
               }
           }
       }
       if !updated {
           break
       }
   }
   if best[n] >= 0 {
       return true
   }
   // check for positive cycle reachable => infinite fuel
   for i := 1; i <= n; i++ {
       if best[i] < 0 {
           continue
       }
       bi := best[i]
       for j := 1; j <= n; j++ {
           if i == j {
               continue
           }
           cost := d * (abs(x[i]-x[j]) + abs(y[i]-y[j]))
           if bi < cost {
               continue
           }
           nf := bi - cost + a[j]
           if nf > best[j] {
               // positive cycle detected, fuel unbounded
               return true
           }
       }
   }
   return false
}

func abs(a int64) int64 {
   if a < 0 {
       return -a
   }
   return a
}
