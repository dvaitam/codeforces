package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for t > 0 {
       t--
       var n int
       var m int64
       fmt.Fscan(reader, &n, &m)
       items := make([]struct{ weight int64; index int }, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &items[i].weight)
           items[i].index = i + 1
       }
       sort.Slice(items, func(i, j int) bool {
           return items[i].weight > items[j].weight
       })
       var sum int64
       var res []int
       half := (m + 1) / 2
       for _, item := range items {
           if sum+item.weight <= m {
               sum += item.weight
               res = append(res, item.index)
           }
       }
       if sum < half {
           fmt.Fprintln(writer, -1)
       } else {
           fmt.Fprintln(writer, len(res))
           // reverse to match order
           for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
               res[i], res[j] = res[j], res[i]
           }
           for i, idx := range res {
               if i > 0 {
                   fmt.Fprint(writer, " ")
               }
               fmt.Fprint(writer, idx)
           }
           fmt.Fprintln(writer)
       }
   }
}
