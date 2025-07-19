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
       var s string
       fmt.Fscan(reader, &n, &s)
       bs := []byte(s)
       l := -1
       for i := 0; i < n; i++ {
           if bs[i] != '?' {
               l = i
           }
       }
       if l == -1 {
           bs[0] = 'B'
           l = 0
       }
       for i := l - 1; i >= 0; i-- {
           if bs[i] == '?' {
               if bs[i+1] == 'B' {
                   bs[i] = 'R'
               } else {
                   bs[i] = 'B'
               }
           }
       }
       for i := l + 1; i < n; i++ {
           if bs[i] == '?' {
               if bs[i-1] == 'B' {
                   bs[i] = 'R'
               } else {
                   bs[i] = 'B'
               }
           }
       }
       writer.WriteString(string(bs))
       writer.WriteByte('\n')
   }
}
