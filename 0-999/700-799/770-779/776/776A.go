package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a, b string
   var n int
   fmt.Fscan(reader, &a, &b)
   fmt.Fscan(reader, &n)
   fmt.Println(a, b)
   for i := 0; i < n; i++ {
       var c, d string
       fmt.Fscan(reader, &c, &d)
       var survived, nw string
       if c == a || d == a {
           // a is replaced
           survived = b
           if c == a {
               nw = d
           } else {
               nw = c
           }
       } else {
           // b is replaced
           survived = a
           if c == b {
               nw = d
           } else {
               nw = c
           }
       }
       a = survived
       b = nw
       fmt.Println(a, b)
   }
}
