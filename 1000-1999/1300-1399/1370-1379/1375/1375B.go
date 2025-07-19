package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var T int
   fmt.Fscan(reader, &T)
   for T > 0 {
       T--
       var n, m int
       fmt.Fscan(reader, &n, &m)
       a := make([][]int, n)
       for i := 0; i < n; i++ {
           a[i] = make([]int, m)
           for j := 0; j < m; j++ {
               fmt.Fscan(reader, &a[i][j])
           }
       }
       valid := true
       b := make([][]int, n)
       for i := 0; i < n; i++ {
           b[i] = make([]int, m)
           for j := 0; j < m; j++ {
               cnt := 0
               if i > 0 {
                   cnt++
               }
               if i < n-1 {
                   cnt++
               }
               if j > 0 {
                   cnt++
               }
               if j < m-1 {
                   cnt++
               }
               b[i][j] = cnt
               if a[i][j] > b[i][j] {
                   valid = false
               }
           }
       }
       if !valid {
           fmt.Fprintln(writer, "NO")
           continue
       }
       fmt.Fprintln(writer, "YES")
       for i := 0; i < n; i++ {
           for j := 0; j < m; j++ {
               fmt.Fprint(writer, b[i][j])
               if j+1 < m {
                   fmt.Fprint(writer, " ")
               }
           }
           fmt.Fprintln(writer)
       }
   }
}
