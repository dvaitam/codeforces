package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([][]int, n)
   b := make([][]int, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(reader, &a[i][j])
       }
   }
   for i := 0; i < n; i++ {
       b[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(reader, &b[i][j])
       }
   }
   inv := make([]int, m)
   for i := 1; i < n; i++ {
       for j := 0; j < m; j++ {
           if b[i][j] < b[i-1][j] {
               inv[j]++
           }
       }
   }
   cut := make([]bool, n-1)
   var qu []int
   for j := 0; j < m; j++ {
       if inv[j] == 0 {
           qu = append(qu, j)
       }
   }
   var ord []int
   for len(qu) > 0 {
       v := qu[len(qu)-1]
       qu = qu[:len(qu)-1]
       ord = append(ord, v)
       for i := 1; i < n; i++ {
           if !cut[i-1] && b[i-1][v] < b[i][v] {
               cut[i-1] = true
               for j := 0; j < m; j++ {
                   if b[i-1][j] > b[i][j] {
                       inv[j]--
                       if inv[j] == 0 {
                           qu = append(qu, j)
                       }
                   }
               }
           }
       }
   }
   // apply sorts in reverse order
   for idx := len(ord) - 1; idx >= 0; idx-- {
       col := ord[idx]
       sort.SliceStable(a, func(i, j int) bool {
           return a[i][col] < a[j][col]
       })
   }
   // check equality
   ok := true
   for i := 0; i < n && ok; i++ {
       for j := 0; j < m; j++ {
           if a[i][j] != b[i][j] {
               ok = false
               break
           }
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   if !ok {
       fmt.Fprintln(writer, -1)
       return
   }
   // output result
   fmt.Fprintln(writer, len(ord))
   for i, v := range ord {
       if i > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, v+1)
   }
   fmt.Fprintln(writer)
}
