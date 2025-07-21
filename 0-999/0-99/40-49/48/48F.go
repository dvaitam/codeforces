package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m int
   var W int64
   fmt.Fscan(in, &n, &m, &W)
   w := make([]int64, m)
   c := make([]int64, m)
   a := make([]int64, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &w[i])
   }
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &c[i])
   }
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &a[i])
   }
   type comp struct {
       price float64
       w     int64
   }
   comps := make([]comp, m)
   var answer float64
   // Process each day
   for day := 0; day < n; day++ {
       // build price list
       for i := 0; i < m; i++ {
           price := (float64(c[i]) - float64(a[i])*float64(day)) / float64(w[i])
           comps[i] = comp{price: price, w: w[i]}
       }
       sort.Slice(comps, func(i, j int) bool {
           return comps[i].price < comps[j].price
       })
       rem := W
       for i := 0; i < m && rem > 0; i++ {
           if comps[i].w <= rem {
               answer += float64(comps[i].w) * comps[i].price
               rem -= comps[i].w
           } else {
               answer += float64(rem) * comps[i].price
               rem = 0
           }
       }
   }
   // Print with decimal point
   fmt.Fprintf(out, "%.9f", answer)
}
