package main

import "fmt"

func main() {
   var n, k int
   if _, err := fmt.Scan(&n, &k); err != nil {
       return
   }
   a := 1
   b := n
   k--
   md := 0
   for i := 1; i <= n; i++ {
       if k == 0 {
           if md == 0 {
               if i%2 == 0 {
                   md = 1
               } else {
                   md = 2
               }
           }
           if md == 1 {
               fmt.Printf("%d ", b)
               b--
           }
           if md == 2 {
               fmt.Printf("%d ", a)
               a++
           }
       } else {
           if i%2 == 0 {
               fmt.Printf("%d ", a)
               a++
           } else {
               fmt.Printf("%d ", b)
               b--
           }
           if i > 1 {
               k--
           }
       }
   }
   fmt.Println()
}
