package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

// vec represents a vector with its original index and computed angle.
type vec struct {
   idx   int
   angle float64
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   v := make([]vec, n)
   for i := 0; i < n; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       v[i].idx = i + 1
       v[i].angle = math.Atan2(float64(y), float64(x))
   }

   sort.Slice(v, func(i, j int) bool {
       return v[i].angle < v[j].angle
   })

   // Initialize with the wrap-around pair
   ans1, ans2 := v[0].idx, v[n-1].idx
   minDiff := v[0].angle + 2*math.Pi - v[n-1].angle

   // Find minimum adjacent angle difference
   for i := 0; i < n-1; i++ {
       diff := v[i+1].angle - v[i].angle
       if diff < minDiff {
           minDiff = diff
           ans1 = v[i].idx
           ans2 = v[i+1].idx
       }
   }

   fmt.Fprintln(writer, ans1, ans2)
}
