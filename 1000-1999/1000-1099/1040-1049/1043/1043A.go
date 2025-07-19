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
   sum := 0
   maxa := 0
   for i := 0; i < n; i++ {
       var a int
       fmt.Fscan(reader, &a)
       sum += a
       if a > maxa {
           maxa = a
       }
   }
   // minimal k such that n*k > 2*sum: k = floor(2*sum/n)+1
   k0 := (2*sum)/n + 1
   if maxa > k0 {
       k0 = maxa
   }
   // output answer
   fmt.Println(k0)
}
