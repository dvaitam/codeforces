package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "time"
)

const MAXN = 210

// Edge for adjacency list
type Edge struct {
   v, next int
   use      bool
}

var (
   edge   [MAXN * 4]Edge
   low    [MAXN]int
   dfn    [MAXN]int
   vis    [MAXN]int
   pre    [MAXN]int
   head   [MAXN]int
   uArr   [MAXN]int
   vArr   [MAXN]int
   p      []int
   viss   []int
   t      []int
   ip, sol, Count, n, m, ans int
   flag   bool
)

func initVars() {
   for i := 0; i < MAXN; i++ {
       head[i] = -1
       vis[i] = 0
   }
   Count, sol, ip = 0, 0, 0
}

// add undirected edge
func addedge(u, v int) {
   edge[ip].v = v
   edge[ip].use = false
   edge[ip].next = head[u]
   head[u] = ip
   ip++
}

func tarjan(u int) {
   vis[u] = 1
   dfn[u] = Count
   low[u] = Count
   Count++
   for i := head[u]; i != -1; i = edge[i].next {
       if !edge[i].use {
           edge[i].use = true
           edge[i^1].use = true
           v := edge[i].v
           if vis[v] == 0 {
               pre[v] = u
               tarjan(v)
               if low[v] < low[u] {
                   low[u] = low[v]
               }
               if dfn[u] < low[v] {
                   sol++
                   flag = false
               }
           } else if vis[v] == 1 {
               if dfn[v] < low[u] {
                   low[u] = dfn[v]
               }
           }
       }
   }
   vis[u] = 2
}

func main() {
   rand.Seed(time.Now().UnixNano())
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &n, &m)
   p = make([]int, m+1)
   viss = make([]int, m+1)
   t = make([]int, m+1)
   for i := 1; i <= m; i++ {
       fmt.Fscan(reader, &uArr[i], &vArr[i])
       p[i] = i
   }
   ans = -1
   // randomized greedy removal
   for jo := 1; jo <= 2000; jo++ {
       // random shuffle p
       for i := 0; i < 100; i++ {
           x := rand.Intn(m) + 1
           y := rand.Intn(m) + 1
           p[x], p[y] = p[y], p[x]
       }
       for i := 1; i <= m; i++ {
           viss[i] = 0
       }
       for i := 1; i <= m; i++ {
           // attempt keep edge i
           viss[i] = 1
           initVars()
           // build graph with edges where viss[j]==0
           for j := 1; j <= m; j++ {
               if viss[j] == 0 {
                   u := uArr[p[j]]
                   v := vArr[p[j]]
                   addedge(u, v)
                   addedge(v, u)
               }
           }
           flag = true
           tarjan(1)
           // check connectivity
           for j := 1; j <= n; j++ {
               if vis[j] == 0 {
                   flag = false
                   break
               }
           }
           if !flag {
               viss[i] = 0
           }
       }
       // count kept edges
       s := 0
       for i := 1; i <= m; i++ {
           if viss[i] == 1 {
               s++
           }
       }
       if ans < s {
           ans = s
           y := 0
           for i := 1; i <= m; i++ {
               if viss[i] == 0 {
                   y++
                   t[y] = p[i]
               }
           }
       }
   }
   // print removed edges count and list
   removeCount := m - ans
   fmt.Fprintln(writer, removeCount)
   for i := 1; i <= removeCount; i++ {
       idx := t[i]
       fmt.Fprintf(writer, "%d %d\n", uArr[idx], vArr[idx])
   }
}
