package main

import "fmt"

func divisors(x int64) []int64 {
   var m []int64
   for i := int64(1); i*i <= x; i++ {
       if x%i == 0 {
           m = append(m, i)
       }
   }
   return m
}

func main() {
   var a, b int64
   if _, err := fmt.Scan(&a, &b); err != nil {
       return
   }
   c := a + b
   am := divisors(a)
   bm := divisors(b)
   cm := divisors(c)
   for i := len(cm) - 1; i >= 0; i-- {
       d := cm[i]
       dd := c / d
       for _, da := range am {
           if d >= da && dd >= a/da {
               fmt.Println((d + dd) * 2)
               return
           }
       }
       for _, db := range bm {
           if d >= db && dd >= b/db {
               fmt.Println((d + dd) * 2)
               return
           }
       }
   }
}
