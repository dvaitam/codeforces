package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   b := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &b[i])
   }
   var joy int64
   for i := 0; i < n; i++ {
       ai := a[i]
       bi := b[i]
       // determine feasible x range: x+y=bi, 1<=x,y<=ai
       L := bi - ai
       if L < 1 {
           L = 1
       }
       R := ai
       if R > bi-1 {
           R = bi - 1
       }
       if L > R {
           joy--
       } else {
           x0 := bi / 2
           var x int
           if x0 < L {
               x = L
           } else if x0 > R {
               x = R
           } else {
               x = x0
           }
           joy += int64(x) * int64(bi-x)
       }
   }
   fmt.Println(joy)
}
