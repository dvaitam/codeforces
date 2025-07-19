package main

import (
   "fmt"
)

func main() {
   // read input from stdin

   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   w := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&w[i])
   }
   var p int
   fmt.Scan(&p)
   const T = 550
   a := make([]int, p)
   b := make([]int, p)
   res := make([]int64, p)
   c := make([][]int, T+1)
   for i := 0; i < p; i++ {
       var ai, bi int
       fmt.Scan(&ai, &bi)
       ai--
       a[i] = ai
       b[i] = bi
       if bi <= T {
           c[bi] = append(c[bi], i)
       } else {
           var cur int64
           for j := ai; j < n; j += bi {
               cur += w[j]
           }
           res[i] = cur
       }
   }
   for bi := 1; bi <= T; bi++ {
       if len(c[bi]) == 0 {
           continue
       }
       s := make([]int64, n)
       for j := 0; j < n; j++ {
           if j < bi {
               s[j] = w[j]
           } else {
               s[j] = w[j] + s[j-bi]
           }
       }
       for _, idx := range c[bi] {
           ai := a[idx]
           k := (n-1 - ai) / bi
           e := ai + k*bi
           if ai >= bi {
               res[idx] = s[e] - s[ai-bi]
           } else {
               res[idx] = s[e]
           }
       }
   }
   for i := 0; i < p; i++ {
       fmt.Println(res[i])
   }
}
