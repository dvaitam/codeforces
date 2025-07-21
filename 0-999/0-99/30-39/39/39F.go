package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int64
   var m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   d := make([]int64, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &d[i])
   }
   p := make([]int64, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(reader, &p[i])
   }
   cnt := make([]int, m)
   // initialize min to a large value
   min := k + 1
   for i := 0; i < m; i++ {
       c := 0
       di := d[i]
       for j := 0; j < k; j++ {
           if p[j]%di == 0 {
               c++
           }
       }
       cnt[i] = c
       if c < min {
           min = c
       }
   }
   var res []int
   for i := 0; i < m; i++ {
       if cnt[i] == min {
           res = append(res, i+1)
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, len(res))
   for i, v := range res {
       if i > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}
