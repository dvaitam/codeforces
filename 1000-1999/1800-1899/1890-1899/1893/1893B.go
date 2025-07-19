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
   for ; t > 0; t-- {
       var n, m int
       fmt.Fscan(reader, &n, &m)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       b := make([]int, m)
       for i := 0; i < m; i++ {
           fmt.Fscan(reader, &b[i])
       }
       sort.Sort(sort.Reverse(sort.IntSlice(b)))
       p := 0
       ans := make([]int, 0, n+m)
       for _, bi := range b {
           for p < n && a[p] >= bi {
               ans = append(ans, a[p])
               p++
           }
           ans = append(ans, bi)
       }
       for p < n {
           ans = append(ans, a[p])
           p++
       }
       for i, v := range ans {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
       writer.WriteByte('\n')
   }
}
