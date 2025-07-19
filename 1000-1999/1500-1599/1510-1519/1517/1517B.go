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
   fmt.Fscan(reader, &t)
   for tc := 0; tc < t; tc++ {
       var n, m int
       fmt.Fscan(reader, &n, &m)
       // read matrix
       v := make([][]int, n)
       flat := make([]int, 0, n*m)
       for i := 0; i < n; i++ {
           v[i] = make([]int, m)
           for j := 0; j < m; j++ {
               fmt.Fscan(reader, &v[i][j])
               flat = append(flat, v[i][j])
           }
       }
       // select m smallest values
       sort.Ints(flat)
       sel := flat[:m]
       // sort each row descending
       for i := 0; i < n; i++ {
           sort.Slice(v[i], func(a, b int) bool { return v[i][a] > v[i][b] })
       }
       // assign selected values in decreasing order
       k := 0
       for idx := m - 1; idx >= 0; idx-- {
           x := sel[idx]
           for i := 0; i < n; i++ {
               // search in row i from k
               for j := k; j < m; j++ {
                   if v[i][j] == x {
                       v[i][j], v[i][k] = v[i][k], v[i][j]
                       k++
                       goto Next
                   }
               }
           }
       Next:
       }
       // output result
       for i := 0; i < n; i++ {
           for j := 0; j < m; j++ {
               if j > 0 {
                   writer.WriteByte(' ')
               }
               fmt.Fprint(writer, v[i][j])
           }
           writer.WriteByte('\n')
       }
   }
}
