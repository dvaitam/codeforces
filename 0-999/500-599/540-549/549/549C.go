package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   odd, even := 0, 0
   for i := 0; i < n; i++ {
       var a int
       fmt.Fscan(reader, &a)
       if a%2 != 0 {
           odd++
       } else {
           even++
       }
   }
   m := n - k
   // moves: m removals
   // if m is odd, Stannis moves last
   // use strategy based on counts
   var winner string
   if m%2 == 1 {
       // Stannis last
       if even <= m/2 {
           winner = "Stannis"
       } else {
           winner = "Daenerys"
       }
   } else {
       // Daenerys last
       if odd <= m/2 {
           winner = "Daenerys"
       } else {
           winner = "Stannis"
       }
   }
   fmt.Println(winner)
}
