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
   var vb, vs int
   if _, err := fmt.Fscan(in, &n, &vb, &vs); err != nil {
       return
   }
   xs := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &xs[i])
   }
   var xu, yu int
   fmt.Fscan(in, &xu, &yu)
   bestIdx := 2
   // initialize with stop 2
   eps := 1e-9
   x := xs[1]
   bestTime := float64(x)/float64(vb) + math.Hypot(float64(xu-x), float64(yu))/float64(vs)
   bestDist := math.Hypot(float64(xu-x), float64(yu))
   for i := 2; i < n; i++ {
       xi := xs[i]
       tBus := float64(xi) / float64(vb)
       dRun := math.Hypot(float64(xu-xi), float64(yu))
       tTotal := tBus + dRun/float64(vs)
       if tTotal+eps < bestTime {
           bestTime = tTotal
           bestDist = dRun
           bestIdx = i + 1
       } else if math.Abs(tTotal-bestTime) <= eps && dRun < bestDist {
           bestDist = dRun
           bestIdx = i + 1
       }
   }
   fmt.Println(bestIdx)
}
