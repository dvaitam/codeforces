package main

import (
   "fmt"
)

func main() {
   var w, m, k uint64
   if _, err := fmt.Scan(&w, &m, &k); err != nil {
       return
   }
   var ans uint64
   rem := w
   curr := m
   for rem > 0 {
       // compute digit count of curr
       tmp := curr
       d := 0
       for tmp > 0 {
           d++
           tmp /= 10
       }
       // upper bound of same digit count
       pow10 := uint64(1)
       for i := 0; i < d; i++ {
           pow10 *= 10
       }
       end := pow10 - 1
       // number of numbers in this block
       count := end - curr + 1
       costPer := uint64(d) * k
       if rem < costPer {
           break
       }
       // max we can take in this block
       maxTake := rem / costPer
       var take uint64
       if maxTake < count {
           take = maxTake
       } else {
           take = count
       }
       ans += take
       rem -= take * costPer
       if take < count {
           break
       }
       curr = end + 1
   }
   fmt.Println(ans)
}
