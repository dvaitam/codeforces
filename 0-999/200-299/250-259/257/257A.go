package main

import (
   "fmt"
   "sort"
)

func main() {
   var n, m, k int
   if _, err := fmt.Scan(&n, &m, &k); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&a[i])
   }
   // Compute net gains of using each filter: gain = sockets on filter - 1 used
   gains := make([]int, n)
   for i, ai := range a {
       gains[i] = ai - 1
   }
   // Sort gains in descending order
   sort.Sort(sort.Reverse(sort.IntSlice(gains)))
   // Current available sockets
   curr := k
   if curr >= m {
       fmt.Println(0)
       return
   }
   used := 0
   for _, g := range gains {
       if g <= 0 {
           break
       }
       curr += g
       used++
       if curr >= m {
           fmt.Println(used)
           return
       }
   }
   // Impossible to reach required sockets
   fmt.Println(-1)
}
