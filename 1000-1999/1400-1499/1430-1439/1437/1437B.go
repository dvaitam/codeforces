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
       var s string
       fmt.Fscan(reader, &n)
       fmt.Fscan(reader, &s)
       // Two patterns: p1 starts with '0', p2 starts with '1'
       cnt1 := 0
       cnt2 := 0
       in1 := false
       in2 := false
       for i := 0; i < n; i++ {
           // expected characters
           var e1, e2 byte
           if i%2 == 0 {
               e1 = '0'
               e2 = '1'
           } else {
               e1 = '1'
               e2 = '0'
           }
           // pattern1
           if s[i] != e1 {
               if !in1 {
                   cnt1++
                   in1 = true
               }
           } else {
               in1 = false
           }
           // pattern2
           if s[i] != e2 {
               if !in2 {
                   cnt2++
                   in2 = true
               }
           } else {
               in2 = false
           }
       }
       // answer is min operations
       if cnt1 < cnt2 {
           fmt.Fprintln(writer, cnt1)
       } else {
           fmt.Fprintln(writer, cnt2)
       }
   }
}
