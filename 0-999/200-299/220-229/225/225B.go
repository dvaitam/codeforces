package main

import (
   "fmt"
)

func main() {
   var s int64
   var kk int64
   if _, err := fmt.Scan(&s, &kk); err != nil {
       return
   }
   k := int(kk)
   // Build k-bonacci sequence until >= s
   f := make([]int64, 0, 1005)
   f = append(f, 0)
   f = append(f, 1)
   for {
       i := len(f)
       var sum int64
       // sum f[i-1] .. f[max(0,i-k)]
       for j := 1; j <= k; j++ {
           if j > i {
               break
           }
           sum += f[i-j]
       }
       f = append(f, sum)
       if sum >= s {
           break
       }
   }
   // Greedy representation
   var v []int64
   remaining := s
   for idx := len(f) - 1; idx >= 0 && remaining > 0; idx-- {
       if remaining >= f[idx] {
           remaining -= f[idx]
           v = append(v, f[idx])
       }
   }
   // include zero as in original solution
   v = append(v, 0)
   // output
   fmt.Println(len(v))
   for _, val := range v {
       fmt.Print(val, " ")
   }
   fmt.Println()
}
