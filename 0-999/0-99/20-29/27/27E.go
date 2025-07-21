package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   primes := []uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53}
   // best initialized to a large value (<=1e18)
   const maxVal = uint64(1e18)
   best := maxVal

   var dfs func(pos, lastExp, remDiv int, cur uint64)
   dfs = func(pos, lastExp, remDiv int, cur uint64) {
       if remDiv == 1 {
           if cur < best {
               best = cur
           }
           return
       }
       if pos >= len(primes) {
           return
       }
       p := primes[pos]
       val := cur
       for exp := 1; exp <= lastExp; exp++ {
           // multiply by p and check overflow or exceeding best
           if val > best/p {
               break
           }
           val *= p
           if remDiv%(exp+1) != 0 {
               continue
           }
           dfs(pos+1, exp, remDiv/(exp+1), val)
       }
   }

   dfs(0, 64, n, 1)
   fmt.Println(best)
}
