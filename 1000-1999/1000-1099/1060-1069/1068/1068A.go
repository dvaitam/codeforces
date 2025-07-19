package main

import (
   "fmt"
)

func main() {
   var n, m, k, l int64
   if _, err := fmt.Scan(&n, &m, &k, &l); err != nil {
       return
   }
   // If total capacity less than one package, or initial fill exceeds capacity, or all slots reserved
   if n < m || l > n || k >= n {
       fmt.Println(-1)
       return
   }
   // Number of students that need packages: initial l plus extra k
   need := l + k
   // Compute minimum number of packages of size m
   packs := need / m
   if need%m != 0 {
       packs++
   }
   // Check total packages do not exceed n/m total possible
   if packs*m <= n {
       fmt.Println(packs)
   } else {
       fmt.Println(-1)
   }
}
