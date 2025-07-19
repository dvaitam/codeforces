package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const oo = 0x0fffffff

var (
   n, m, limit int
   a            []Edge
   f            []int
   u, num       []int
   vis, vv      []bool
   d, pos       []int
   s, eCnt      int
   first        []int
   g            []GEdge
   tot          int
   rdr          = bufio.NewReader(os.Stdin)
)

// Edge represents an input edge
type Edge struct {
   u, v, c, q int
   tip        bool
}

// GEdge represents graph edge for MST forest
type GEdge struct {
   u, v, c, q int
   tip        bool
   next, b    int
}

func readInt() int {
   var x int
   fmt.Fscan(rdr, &x)
   return x
}

func getf(x int) int {
   if f[x] != x {
       f[x] = getf(f[x])
   }
   return f[x]
}

func add(u0, v0, c0, q0 int) {
   tot++
   g[tot] = GEdge{u: u0, v: v0, c: c0, q: q0, b: tot + 1, next: first[u0]}
   first[u0] = tot
   tot++
   g[tot] = GEdge{u: v0, v: u0, c: c0, q: q0, b: tot - 1, next: first[v0]}
   first[v0] = tot
}

func mst() {
   // MST excluding node 1 edges
   for i := 1; i <= m; i++ {
       if a[i].u == 1 || a[i].v == 1 {
           continue
       }
       f1 := getf(a[i].u)
       f2 := getf(a[i].v)
       if f1 == f2 {
           continue
       }
       s += a[i].c
       a[i].tip = true
       f[f2] = f1
   }
   // Add up to limit edges from node 1
   for i := 1; i <= m && eCnt < limit; i++ {
       if a[i].u == 1 {
           f1 := getf(a[i].v)
           if f1 != 1 {
               s += a[i].c
               a[i].tip = true
               vis[a[i].v] = true
               f[f1] = 1
               eCnt++
           }
       } else if a[i].v == 1 {
           f1 := getf(a[i].u)
           if f1 != 1 {
               s += a[i].c
               a[i].tip = true
               vis[a[i].u] = true
               f[f1] = 1
               eCnt++
           }
       }
   }
   // Build forest g excluding node 1 edges
   for i := 1; i <= m; i++ {
       if a[i].tip && a[i].u != 1 && a[i].v != 1 {
           add(a[i].u, a[i].v, a[i].c, a[i].q)
       }
   }
}

func dfs1(x, curMax, last int) {
   vv[x] = true
   d[x] = curMax
   pos[x] = last
   for i := first[x]; i != 0; i = g[i].next {
       ge := g[i]
       if vv[ge.v] || ge.tip {
           continue
       }
       nextMax := curMax
       lastIdx := last
       if ge.c > curMax {
           nextMax = ge.c
           lastIdx = i
       }
       dfs1(ge.v, nextMax, lastIdx)
   }
}

func work() {
   // initial dfs from vis vertices
   for i := range vv {
       vv[i] = false
   }
   for i := 1; i <= n; i++ {
       if vis[i] {
           dfs1(i, 0, 0)
       }
   }
   for eCnt < limit {
       // find replacement
       k := 0
       mins := oo
       for i := 2; i <= n; i++ {
           if vis[i] || u[i] == 0 {
               continue
           }
           delta := u[i] - d[i]
           if delta <= mins {
               mins = delta
               k = i
           }
       }
       if k == 0 {
           break
       }
       s += mins
       // include the new edge
       pp := pos[k]
       g[pp].tip = true
       g[g[pp].b].tip = true
       vis[k] = true
       // recompute dfs
       for i := range vv {
           vv[i] = false
       }
       for i := 1; i <= n; i++ {
           if vis[i] {
               dfs1(i, 0, 0)
           }
       }
       eCnt++
   }
}

func output1() bool {
   // check connectivity
   for i := 1; i <= n; i++ {
       if getf(i) != 1 {
           return true
       }
   }
   return false
}

func main() {
   n = readInt()
   m = readInt()
   limit = readInt()
   a = make([]Edge, m+1)
   u = make([]int, n+1)
   num = make([]int, n+1)
   vis = make([]bool, n+1)
   vv = make([]bool, n+1)
   d = make([]int, n+1)
   pos = make([]int, n+1)
   f = make([]int, n+1)
   first = make([]int, n+1)
   g = make([]GEdge, 2*m+2)
   for i := 1; i <= m; i++ {
       a[i].u = readInt()
       a[i].v = readInt()
       a[i].c = readInt()
       a[i].q = i
       if a[i].u == 1 {
           u[a[i].v] = a[i].c
           num[a[i].v] = i
       }
       if a[i].v == 1 {
           u[a[i].u] = a[i].c
           num[a[i].u] = i
       }
   }
   // sort edges by cost
   sort.Slice(a[1:], func(i, j int) bool {
       return a[i+1].c < a[j+1].c
   })
   // init uf
   for i := 1; i <= n; i++ {
       f[i] = i
   }
   mst()
   if output1() {
       fmt.Println(-1)
       return
   }
   work()
   if eCnt != limit {
       fmt.Println(-1)
       return
   }
   // print result
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, n-1)
   // edges from node 1
   for i := 2; i <= n; i++ {
       if vis[i] {
           fmt.Fprintf(w, "%d ", num[i])
       }
   }
   // other edges
   for i := 1; i <= tot; i++ {
       if g[i].tip {
           continue
       }
       fmt.Fprintf(w, "%d ", g[i].q)
       g[i].tip = true
       g[g[i].b].tip = true
   }
}
