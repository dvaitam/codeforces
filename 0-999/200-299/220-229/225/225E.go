package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   const mod = 1000000007
   count := 0
   var z int
   for z = 1; ; z++ {
       found := false
       // try all x from 1 to 2*z
       for x := 1; x <= 2*z; x++ {
           f := x / 2
           num := z - f
           if num <= 0 {
               // further x will only increase f
               break
           }
           denom := x + 1
           if num%denom == 0 {
               y := num / denom
               if y > 0 {
                   found = true
                   break
               }
           }
       }
       if !found {
           count++
           if count == n {
               // output z modulo mod
               ans := int64(z) % mod
               fmt.Println(ans)
               return
           }
       }
   }
}
