package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, h int
   if _, err := fmt.Fscan(reader, &n, &m, &h); err != nil {
       return
   }
   w := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &w[i])
   }
   graph := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       if (w[u]+1)%h == w[v] {
           graph[u] = append(graph[u], v)
       }
       if (w[v]+1)%h == w[u] {
           graph[v] = append(graph[v], u)
       }
   }
   // Tarjan's algorithm
   dfn := make([]int, n+1)
   low := make([]int, n+1)
   onStack := make([]bool, n+1)
   stack := make([]int, 0, n)
   var index, sccCount int
   sccId := make([]int, n+1)
   size := make([]int, n+1)
   var dfs func(u int)
   dfs = func(u int) {
       index++
       dfn[u] = index
       low[u] = index
       stack = append(stack, u)
       onStack[u] = true
       for _, v := range graph[u] {
           if dfn[v] == 0 {
               dfs(v)
               low[u] = min(low[u], low[v])
           } else if onStack[v] {
               low[u] = min(low[u], dfn[v])
           }
       }
       if low[u] == dfn[u] {
           sccCount++
           for {
               sz := stack[len(stack)-1]
               stack = stack[:len(stack)-1]
               onStack[sz] = false
               sccId[sz] = sccCount
               size[sccCount]++
               if sz == u {
                   break
               }
           }
       }
   }
   for i := 1; i <= n; i++ {
       if dfn[i] == 0 {
           dfs(i)
       }
   }
   outDeg := make([]int, sccCount+1)
   for u := 1; u <= n; u++ {
       for _, v := range graph[u] {
           if sccId[u] != sccId[v] {
               outDeg[sccId[u]]++
           }
       }
   }
   // find smallest sink SCC
   bestSize := n + 1
   bestId := 0
   for i := 1; i <= sccCount; i++ {
       if outDeg[i] == 0 && size[i] < bestSize {
           bestSize = size[i]
           bestId = i
       }
   }
   fmt.Fprintln(writer, bestSize)
   for i := 1; i <= n; i++ {
       if sccId[i] == bestId {
           fmt.Fprintln(writer, i)
       }
   }
}
