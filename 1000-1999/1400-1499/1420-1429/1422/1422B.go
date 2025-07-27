package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func abs(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
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
       a := make([][]int64, n)
       for i := 0; i < n; i++ {
           a[i] = make([]int64, m)
           for j := 0; j < m; j++ {
               fmt.Fscan(reader, &a[i][j])
           }
       }
       var res int64
       // process groups of symmetric positions
       for i := 0; i <= (n-1)/2; i++ {
           for j := 0; j <= (m-1)/2; j++ {
               i2 := n - 1 - i
               j2 := m - 1 - j
               vals := make([]int64, 0, 4)
               vals = append(vals, a[i][j])
               if i2 != i {
                   vals = append(vals, a[i2][j])
               }
               if j2 != j {
                   vals = append(vals, a[i][j2])
               }
               if i2 != i && j2 != j {
                   vals = append(vals, a[i2][j2])
               }
               sort.Slice(vals, func(x, y int) bool { return vals[x] < vals[y] })
               median := vals[len(vals)/2]
               for _, v := range vals {
                   res += abs(v - median)
               }
           }
       }
       fmt.Fprintln(writer, res)
   }
}
