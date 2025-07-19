package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type item struct {
   w, x, id int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n, m int
       fmt.Fscan(reader, &n, &m)
       a := make([]item, m)
       for i := 0; i < m; i++ {
           var x, w int
           fmt.Fscan(reader, &x, &w)
           a[i] = item{w: w, x: x, id: i + 1}
       }
       // sort by weight ascending
       sort.Slice(a, func(i, j int) bool {
           if a[i].w != a[j].w {
               return a[i].w < a[j].w
           }
           return a[i].x < a[j].x
       })
       // take smallest 2*n by weight
       b := a[:2*n]
       // sort selected by x ascending
       sort.Slice(b, func(i, j int) bool {
           return b[i].x < b[j].x
       })
       var sum int64
       for _, it := range b {
           sum += int64(it.w)
       }
       fmt.Fprintln(writer, sum)

       left, right := 0, len(b)-1
       for left < right {
           fmt.Fprintf(writer, "%d %d\n", b[left].id, b[right].id)
           left++
           right--
       }
       // blank line after each test case
       fmt.Fprintln(writer)
   }
}
