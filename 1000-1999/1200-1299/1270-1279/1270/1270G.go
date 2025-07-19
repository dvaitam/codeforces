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
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := n - 1; i >= 0; i-- {
           var x int
           fmt.Fscan(reader, &x)
           a[i] = x + i
       }
       b := make([]int, n)
       s := make([]int, 0, n)
       x := 0
       for {
           if b[x] != 0 {
               start := b[x] - 1
               length := len(s) - start
               // output cycle length
               fmt.Fprint(writer, length, "\n")
               // output cycle nodes transformed
               for i := start; i < len(s); i++ {
                   fmt.Fprint(writer, n - s[i])
                   if i+1 < len(s) {
                       fmt.Fprint(writer, " ")
                   } else {
                       fmt.Fprint(writer, "\n")
                   }
               }
               break
           }
           s = append(s, x)
           b[x] = len(s)
           x = a[x]
       }
   }
}
