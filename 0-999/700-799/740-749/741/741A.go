package main

import (
   "fmt"
)

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func lcm(a, b int64) int64 {
   return a / gcd(a, b) * b
}

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   crush := make([]int, n+1)
   indeg := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Scan(&crush[i])
       indeg[crush[i]]++
   }
   // Must be permutation: every node has exactly one incoming
   for i := 1; i <= n; i++ {
       if indeg[i] != 1 {
           fmt.Println(-1)
           return
       }
   }
   visited := make([]bool, n+1)
   var res int64 = 1
   for i := 1; i <= n; i++ {
       if visited[i] {
           continue
       }
       // traverse the cycle starting from i
       u := i
       var cnt int64 = 0
       for !visited[u] {
           visited[u] = true
           cnt++
           u = crush[u]
       }
       // cnt is cycle length
       var need int64
       if cnt%2 == 0 {
           need = cnt / 2
       } else {
           need = cnt
       }
       res = lcm(res, need)
   }
   fmt.Println(res)
}
