package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type item struct {
   a int
   t int
   s float64
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, T int
   if _, err := fmt.Fscan(reader, &n, &T); err != nil {
       return
   }
   items := make([]item, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &items[i].a)
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &items[i].t)
       items[i].t -= T
   }
   sort.Slice(items, func(i, j int) bool {
       return items[i].t < items[j].t
   })
   // if all times are on one side, no valid interval
   if items[0].t > 0 || items[n-1].t < 0 {
       fmt.Printf("0.000000")
       return
   }
   var tp float64
   for i := range items {
       items[i].s = float64(items[i].a) * float64(items[i].t)
       tp += items[i].s
   }
   // ensure total slope non-negative
   if tp < 0 {
       for i := range items {
           items[i].t = -items[i].t
           items[i].s = -items[i].s
       }
       // reverse slice
       for i := 0; i < n/2; i++ {
           items[i], items[n-1-i] = items[n-1-i], items[i]
       }
   }
   tp = 0
   var ans float64
   for i := range items {
       if tp+items[i].s > 0 {
           ans += -tp / float64(items[i].t)
           break
       }
       tp += items[i].s
       ans += float64(items[i].a)
   }
   fmt.Printf("%.6f", ans)
}
