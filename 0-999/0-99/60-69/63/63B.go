package main

import (
   "fmt"
)

func main() {
   var n, k int
   if _, err := fmt.Scan(&n, &k); err != nil {
       return
   }
   // counts of soldiers at each rank, indices 1..k
   counts := make([]int, k+2)
   for i := 0; i < n; i++ {
       var a int
       fmt.Scan(&a)
       if a >= 1 && a <= k {
           counts[a]++
       }
   }
   sessions := 0
   // simulate training sessions until all soldiers reach rank k
   for {
       done := true
       for r := 1; r < k; r++ {
           if counts[r] > 0 {
               done = false
               break
           }
       }
       if done {
           break
       }
       sessions++
       // perform one training session
       newCounts := make([]int, k+2)
       // copy current counts
       copy(newCounts, counts)
       for r := 1; r < k; r++ {
           if counts[r] > 0 {
               newCounts[r]--
               newCounts[r+1]++
           }
       }
       counts = newCounts
   }
   fmt.Println(sessions)
}
