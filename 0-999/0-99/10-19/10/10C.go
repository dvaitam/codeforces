package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int64
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   var tot int64
   for i := int64(1); i <= n; i++ {
       tot -= n / i
   }
   var c [9]int64
   for i := int64(0); i < 9; i++ {
       c[i] = n/9
       if n%9 >= i {
           c[i]++
       }
   }
   c[0]--
   for i := int64(0); i < 9; i++ {
       for j := int64(0); j < 9; j++ {
           for k := int64(0); k < 9; k++ {
               if (i*j)%9 == k {
                   tot += c[i] * c[j] * c[k]
               }
           }
       }
   }
   fmt.Println(tot)
}
