package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   cost := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &cost[i])
   }
   var m int
   fmt.Fscan(reader, &m)
   g := make([][]int, n)
   gr := make([][]int, n)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       u--
       v--
       g[u] = append(g[u], v)
       gr[v] = append(gr[v], u)
   }
   visited := make([]bool, n)
   order := make([]int, 0, n)

   var dfs1 func(int)
   dfs1 = func(u int) {
       visited[u] = true
       for _, v := range g[u] {
           if !visited[v] {
               dfs1(v)
           }
       }
       order = append(order, u)
   }
   for i := 0; i < n; i++ {
       if !visited[i] {
           dfs1(i)
       }
   }
   for i := range visited {
       visited[i] = false
   }

   const mod = int64(1000000007)
   var totalCost int64
   var ways int64 = 1

   var dfs2 func(int, *int64, *int64)
   dfs2 = func(u int, compMin *int64, compCnt *int64) {
       visited[u] = true
       if cost[u] < *compMin {
           *compMin = cost[u]
           *compCnt = 1
       } else if cost[u] == *compMin {
           *compCnt++
       }
       for _, v := range gr[u] {
           if !visited[v] {
               dfs2(v, compMin, compCnt)
           }
       }
   }
   for i := len(order) - 1; i >= 0; i-- {
       u := order[i]
       if !visited[u] {
           compMin := int64(math.MaxInt64)
           compCnt := int64(0)
           dfs2(u, &compMin, &compCnt)
           totalCost += compMin
           ways = ways * compCnt % mod
       }
   }
   fmt.Println(totalCost, ways)
