package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   a := make([][]int, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(reader, &a[i][j])
       }
   }
   r := make([]int, n)
   c := make([]int, m)
   rowSum := make([]int, n)
   colSum := make([]int, m)
   for i := 0; i < n; i++ {
       sum := 0
       for j := 0; j < m; j++ {
           sum += a[i][j]
       }
       rowSum[i] = sum
   }
   for j := 0; j < m; j++ {
       sum := 0
       for i := 0; i < n; i++ {
           sum += a[i][j]
       }
       colSum[j] = sum
   }
   changed := true
   for changed {
       changed = false
       // flip negative-sum rows
       for i := 0; i < n; i++ {
           if rowSum[i] < 0 {
               changed = true
               r[i] ^= 1
               rowSum[i] = -rowSum[i]
               for j := 0; j < m; j++ {
                   oldParity := (r[i]^1 + c[j]) & 1
                   var oldVal int
                   if oldParity == 1 {
                       oldVal = -a[i][j]
                   } else {
                       oldVal = a[i][j]
                   }
                   colSum[j] += -2 * oldVal
               }
           }
       }
       // flip negative-sum columns
       for j := 0; j < m; j++ {
           if colSum[j] < 0 {
               changed = true
               c[j] ^= 1
               colSum[j] = -colSum[j]
               for i := 0; i < n; i++ {
                   oldParity := (r[i] + (c[j]^1)) & 1
                   var oldVal int
                   if oldParity == 1 {
                       oldVal = -a[i][j]
                   } else {
                       oldVal = a[i][j]
                   }
                   rowSum[i] += -2 * oldVal
               }
           }
       }
   }
   // collect result
   var rows []int
   for i := 0; i < n; i++ {
       if r[i] == 1 {
           rows = append(rows, i+1)
       }
   }
   var cols []int
   for j := 0; j < m; j++ {
       if c[j] == 1 {
           cols = append(cols, j+1)
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // output rows
   fmt.Fprint(writer, len(rows))
   if len(rows) > 0 {
       fmt.Fprint(writer, " ")
       for i, v := range rows {
           if i > 0 {
               fmt.Fprint(writer, " ")
           }
           fmt.Fprint(writer, v)
       }
   }
   fmt.Fprintln(writer)
   // output columns
   fmt.Fprint(writer, len(cols))
   if len(cols) > 0 {
       fmt.Fprint(writer, " ")
       for i, v := range cols {
           if i > 0 {
               fmt.Fprint(writer, " ")
           }
           fmt.Fprint(writer, v)
       }
   }
   fmt.Fprintln(writer)
}
