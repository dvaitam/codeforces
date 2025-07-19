package main

import (
   "fmt"
   "math/rand"
)

func main() {
   var a [4]int
   // read four integers
   for i := 0; i < 4; i++ {
       if _, err := fmt.Scan(&a[i]); err != nil {
           return
       }
   }
   // transform until all are 1
   for {
       if a[0] == 1 && a[1] == 1 && a[2] == 1 && a[3] == 1 {
           break
       }
       flag := false
       // divide adjacent evens
       for i := 0; i < 4; i++ {
           if a[i]%2 == 0 && a[(i+1)%4]%2 == 0 {
               a[i] /= 2
               a[(i+1)%4] /= 2
               fmt.Printf("/%d\n", i+1)
               flag = true
           }
       }
       if flag {
           continue
       }
       // increment a random adjacent pair
       x := rand.Intn(4)
       a[x]++
       a[(x+1)%4]++
       fmt.Printf("+%d\n", x+1)
   }
}
