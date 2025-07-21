package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type event struct {
   pos   float64
   delta int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var l, v1, v2 int64
   if _, err := fmt.Fscan(in, &n, &l, &v1, &v2); err != nil {
       return
   }
   a := make([]float64, n)
   for i := 0; i < n; i++ {
       var ai int64
       fmt.Fscan(in, &ai)
       a[i] = float64(ai)
   }
   m := float64(2 * l)
   R := float64(l) * float64(v2) / float64(v1+v2)
   events := make([]event, 0, 2*n)
   // build events
   for i := 0; i < n; i++ {
       ai := a[i]
       s := ai - R
       if s < 0 {
           s += m
       }
       events = append(events, event{pos: s, delta: +1})
       events = append(events, event{pos: ai, delta: -1})
   }
   sort.Slice(events, func(i, j int) bool {
       return events[i].pos < events[j].pos
   })
   // initial count: number of ai < R
   c := 0
   for i := 0; i < n; i++ {
       if a[i] < R {
           c++
       }
   }
   ans := make([]float64, n+1)
   prev := 0.0
   for _, e := range events {
       if e.pos > prev {
           ans[c] += e.pos - prev
       }
       c += e.delta
       prev = e.pos
   }
   // last segment
   if prev < m {
       ans[c] += m - prev
   }
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for i := 0; i <= n; i++ {
       // probability = length / m
       prob := ans[i] / m
       fmt.Fprintf(w, "%.10f\n", prob)
   }
}
