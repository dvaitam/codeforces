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

   var n, m int
   fmt.Fscan(reader, &n, &m)
   points := make([]bool, m+1)
   for i := 0; i < n; i++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       if l < 1 {
           l = 1
       }
       if r > m {
           r = m
       }
       for j := l; j <= r; j++ {
           points[j] = true
       }
   }

   var missing []int
   for i := 1; i <= m; i++ {
       if !points[i] {
           missing = append(missing, i)
       }
   }

   fmt.Fprintln(writer, len(missing))
   for i, v := range missing {
       if i > 0 {
           writer.WriteString(" ")
       }
       writer.WriteString(fmt.Sprint(v))
   }
   writer.WriteString("\n")
}
