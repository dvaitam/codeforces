package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, m int
   adj [][]int
   degree []int
   vis [][2]bool
   path []int
   // tarjan variables
   dfn, low []int
   inStack []bool
   stack []int
   tm, scc, mx int
   compSize []int
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   fmt.Fscan(reader, &n, &m)
   adj = make([][]int, n+1)
   degree = make([]int, n+1)
   for i := 1; i <= n; i++ {
       var c int
       fmt.Fscan(reader, &c)
       degree[i] = c
       adj[i] = make([]int, c)
       for j := 0; j < c; j++ {
           fmt.Fscan(reader, &adj[i][j])
       }
   }
   var start int
   fmt.Fscan(reader, &start)
   vis = make([][2]bool, n+1)
   path = []int{start}
   dfs(start, 0)
   // no winning path found, check for draw or lose
   dfn = make([]int, n+1)
   low = make([]int, n+1)
   inStack = make([]bool, n+1)
   compSize = make([]int, n+1)
   tarjan(start)
   if mx > 1 {
       fmt.Println("Draw")
   } else {
       fmt.Println("Lose")
   }
}

func dfs(u, mk int) {
   for _, v := range adj[u] {
       if vis[v][mk^1] {
           continue
       }
       vis[v][mk^1] = true
       path = append(path, v)
       if mk == 0 && degree[v] == 0 {
           fmt.Println("Win")
           // print path
           for i, x := range path {
               if i > 0 {
                   fmt.Print(" ")
               }
               fmt.Print(x)
           }
           fmt.Println()
           os.Exit(0)
       }
       dfs(v, mk^1)
       path = path[:len(path)-1]
   }
}

func tarjan(u int) {
   tm++
   dfn[u] = tm
   low[u] = tm
   inStack[u] = true
   stack = append(stack, u)
   for _, v := range adj[u] {
       if dfn[v] == 0 {
           tarjan(v)
           if low[v] < low[u] {
               low[u] = low[v]
           }
       } else if inStack[v] {
           if dfn[v] < low[u] {
               low[u] = dfn[v]
           }
       }
   }
   if low[u] == dfn[u] {
       scc++
       cnt := 0
       for {
           w := stack[len(stack)-1]
           stack = stack[:len(stack)-1]
           inStack[w] = false
           cnt++
           if w == u {
               break
           }
       }
       compSize[scc] = cnt
       if cnt > mx {
           mx = cnt
       }
   }
}
