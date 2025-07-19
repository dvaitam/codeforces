package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   x := make([]int64, n)
   y := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &x[i], &y[i])
   }
   t := make([]int64, n)
   for i := 0; i < n; i++ {
       dx := x[i] - x[0]
       if dx < 0 {
           dx = -dx
       }
       dy := y[i] - y[0]
       if dy < 0 {
           dy = -dy
       }
       t[i] = dx*dx + dy*dy
   }
   for {
       o := 0
       for i := 0; i < n; i++ {
           if t[i]&1 == 1 {
               o++
           }
       }
       if o > 0 {
           fmt.Println(o)
           first := true
           for i := 0; i < n; i++ {
               if t[i]&1 == 1 {
                   if !first {
                       fmt.Print(" ")
                   }
                   fmt.Print(i+1)
                   first = false
               }
           }
           fmt.Println()
           return
       }
       for i := 0; i < n; i++ {
           t[i] >>= 1
       }
   }
}
