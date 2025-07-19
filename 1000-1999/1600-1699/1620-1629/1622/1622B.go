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

   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       p := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &p[i])
       }
       var s string
       fmt.Fscan(reader, &s)
       // count zeros (disliked)
       z := 0
       for i := 0; i < n; i++ {
           if s[i] == '0' {
               z++
           }
       }
       x := make([]int, 0)
       y := make([]int, 0)
       for i := 0; i < n; i++ {
           if s[i] == '0' && p[i] > z {
               x = append(x, i)
           }
           if s[i] == '1' && p[i] <= z {
               y = append(y, i)
           }
       }
       // swap mismatched positions
       for i := 0; i < len(x); i++ {
           xi := x[i]
           yi := y[i]
           p[xi], p[yi] = p[yi], p[xi]
       }
       // output result
       for i, v := range p {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
       writer.WriteByte('\n')
   }
}
