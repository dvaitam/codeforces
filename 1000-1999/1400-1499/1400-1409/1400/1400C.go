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
   for ; t > 0; t-- {
       var s string
       var x int
       fmt.Fscan(reader, &s)
       fmt.Fscan(reader, &x)
       n := len(s)
       w := make([]byte, n)
       for i := range w {
           w[i] = '1'
       }
       for i := 0; i < n; i++ {
           if s[i] == '0' {
               if i-x >= 0 {
                   w[i-x] = '0'
               }
               if i+x < n {
                   w[i+x] = '0'
               }
           }
       }
       ok := true
       for i := 0; i < n; i++ {
           if s[i] == '1' {
               left := false
               right := false
               if i-x >= 0 && w[i-x] == '1' {
                   left = true
               }
               if i+x < n && w[i+x] == '1' {
                   right = true
               }
               if !left && !right {
                   ok = false
                   break
               }
           }
       }
       if !ok {
           fmt.Fprintln(writer, -1)
       } else {
           fmt.Fprintln(writer, string(w))
       }
   }
}
