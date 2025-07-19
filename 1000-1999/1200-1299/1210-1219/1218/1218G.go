package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, m    int
   ed       [][3]int
   g        [][]int
   t        string
   h, par, col []int
   used     []bool
   out      = bufio.NewWriter(os.Stdout)
)

func printAns() {
   for i := 0; i < m; i++ {
       x := ed[i][2] % 3
       for x <= 0 {
           x += 3
       }
       fmt.Fprintf(out, "%d %d %d\n", ed[i][0], ed[i][1], x)
   }
   out.Flush()
   os.Exit(0)
}

func getOther(id, v int) int {
   return ed[id][0] ^ ed[id][1] ^ v
}

func dfsSolve(v int) {
   for _, id := range g[v] {
       u := getOther(id, v)
       if used[u] {
           continue
       }
       used[u] = true
       dfsSolve(u)
       ed[id][2] = col[u]
       col[u] -= ed[id][2]
       col[v] -= ed[id][2]
       if col[v] < 0 {
           col[v] += 3
       }
   }
}

func solveOdd(id0 int) {
   v, u := ed[id0][0], ed[id0][1]
   if h[v] < h[u] {
       v, u = u, v
   }
   type pii struct{first, second int}
   var cyc []pii
   for v != u {
       cyc = append(cyc, pii{v, par[v]})
       v = getOther(par[v], v)
   }
   cyc = append(cyc, pii{u, id0})
   // init cols from t
   for i := 0; i < n; i++ {
       col[i] = int(t[i] - 'X')
   }
   for _, p := range cyc {
       used[p.first] = true
   }
   for _, p := range cyc {
       dfsSolve(p.first)
   }
   sum := 0
   for _, p := range cyc {
       sum += col[p.first]
   }
   if sum&1 != 0 {
       sum += 3
   }
   sum /= 2
   for i := 2; i < len(cyc); i += 2 {
       sum -= col[cyc[i].first]
   }
   sum %= 3
   if sum < 0 {
       sum += 3
   }
   ed[cyc[0].second][2] = sum
   for i := 1; i < len(cyc); i++ {
       x := col[cyc[i].first] - ed[cyc[i-1].second][2]
       if x < 0 {
           x += 3
       }
       ed[cyc[i].second][2] = x
   }
   printAns()
}

func solveEven(id0 int) {
   v, u := ed[id0][0], ed[id0][1]
   if h[v] < h[u] {
       v, u = u, v
   }
   type pii struct{first, second int}
   var cyc []pii
   for v != u {
       cyc = append(cyc, pii{v, par[v]})
       v = getOther(par[v], v)
   }
   cyc = append(cyc, pii{u, id0})
   for i := 0; i < n; i++ {
       col[i] = h[i] & 1
   }
   if col[cyc[0].first] == 1 {
       // rotate left by 1
       cyc = append(cyc[1:], cyc[0])
   }
   for _, p := range cyc {
       used[p.first] = true
   }
   for _, p := range cyc {
       dfsSolve(p.first)
   }
   bal := 0
   for i, p := range cyc {
       x := col[p.first]
       if i%2 == 0 {
           bal += x
       } else {
           bal -= x
       }
   }
   bal %= 3
   if bal < 0 {
       bal += 3
   }
   if bal == 1 {
       v := cyc[0].first
       col[v]--
       if col[v] < 0 {
           col[v] += 3
       }
   } else if bal == 2 {
       for i := 0; i < 3; i += 2 {
           v := cyc[i].first
           col[v]--
           if col[v] < 0 {
               col[v] += 3
           }
       }
   }
   for i := 1; i < len(cyc); i++ {
       x := col[cyc[i].first] - ed[cyc[i-1].second][2]
       if x < 0 {
           x += 3
       }
       ed[cyc[i].second][2] = x
   }
   printAns()
}

func dfs1(v int) {
   for _, id := range g[v] {
       u := getOther(id, v)
       if h[u] == -1 {
           h[u] = h[v] + 1
           par[u] = id
           dfs1(u)
       } else if (h[v]&1) == (h[u]&1) {
           solveOdd(id)
       }
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   fmt.Fscan(in, &n, &m)
   fmt.Fscan(in, &t)
   ed = make([][3]int, m)
   g = make([][]int, n)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &ed[i][0], &ed[i][1])
       g[ed[i][0]] = append(g[ed[i][0]], i)
       g[ed[i][1]] = append(g[ed[i][1]], i)
   }
   h = make([]int, n)
   par = make([]int, n)
   col = make([]int, n)
   used = make([]bool, n)
   for i := 0; i < n; i++ {
       h[i] = -1
   }
   for i := 0; i < n; i++ {
       if len(g[i]) == 1 {
           for j := 0; j < n; j++ {
               col[j] = int(t[j] - 'X')
           }
           used = make([]bool, n)
           dfsSolve(i)
           printAns()
       }
   }
   h[0] = 0
   dfs1(0)
   for v := 0; v < n; v++ {
       for _, id := range g[v] {
           if par[v] == id {
               continue
           }
           u := getOther(id, v)
           if h[u] < h[v] {
               solveEven(id)
           }
       }
   }
   panic("no solution")
}
