package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Interval represents a subarray with start L and end R (1-based indices).
type Interval struct {
   L, R int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   prefix := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       prefix[i] = prefix[i-1] + a[i]
   }

   // Group intervals by their sum
   M := make(map[int64][]Interval)
   for i := 1; i <= n; i++ {
       for j := i; j <= n; j++ {
           s := prefix[j] - prefix[i-1]
           M[s] = append(M[s], Interval{L: i, R: j})
       }
   }

   var best []Interval
   // For each sum, select maximum number of non-overlapping intervals
   for _, intervals := range M {
       sort.Slice(intervals, func(i, j int) bool {
           if intervals[i].R == intervals[j].R {
               return intervals[i].L < intervals[j].L
           }
           return intervals[i].R < intervals[j].R
       })
       var tmp []Interval
       lastR := 0
       for _, iv := range intervals {
           if iv.L > lastR {
               tmp = append(tmp, iv)
               lastR = iv.R
           }
       }
       if len(tmp) > len(best) {
           best = tmp
       }
   }

   // Output result
   fmt.Fprintln(writer, len(best))
   for _, iv := range best {
       fmt.Fprintln(writer, iv.L, iv.R)
   }
}
