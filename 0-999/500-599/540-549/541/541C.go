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
   f := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&f[i])
       f[i]--
   }
   // compute indegrees
   indegree := make([]int, n)
   for i := 0; i < n; i++ {
       indegree[f[i]]++
   }
   // find cycle nodes by removing non-cycle
   inCycle := make([]bool, n)
   for i := range inCycle {
       inCycle[i] = true
   }
   queue := make([]int, 0)
   for i := 0; i < n; i++ {
       if indegree[i] == 0 {
           queue = append(queue, i)
       }
   }
   for idx := 0; idx < len(queue); idx++ {
       u := queue[idx]
       inCycle[u] = false
       v := f[u]
       indegree[v]--
       if indegree[v] == 0 {
           queue = append(queue, v)
       }
   }
   // collect cycle lengths
   visited := make([]bool, n)
   var L int64 = 1
   for i := 0; i < n; i++ {
       if inCycle[i] && !visited[i] {
           // traverse cycle
           v := i
           cnt := int64(0)
           for {
               visited[v] = true
               cnt++
               v = f[v]
               if v == i {
                   break
               }
           }
           L = lcm(L, cnt)
       }
   }
   // compute depths to cycle
   depth := make([]int, n)
   var maxDepth int
   var calcDepth func(int) int
   calcDepth = func(u int) int {
       if inCycle[u] {
           return 0
       }
       if depth[u] != 0 {
           return depth[u]
       }
       depth[u] = calcDepth(f[u]) + 1
       return depth[u]
   }
   for i := 0; i < n; i++ {
       d := calcDepth(i)
       if d > maxDepth {
           maxDepth = d
       }
   }
   // minimal k: multiple of L and >= maxDepth
   var k int64
   if maxDepth == 0 {
       k = L
   } else {
       t := (int64(maxDepth) + L - 1) / L
       if t < 1 {
           t = 1
       }
       k = t * L
   }
   fmt.Println(k)
