package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

var (
   n, src, snk, visID int
   a     []int
   adj   [][]int
   flow  [][]int
   vis   []int
   ans   [][]int
   tmp   []int
)

func isPrimeSum(x, y int) bool {
   s := x + y
   if s <= 2 || s%2 == 0 {
       return false
   }
   lim := int(math.Sqrt(float64(s)))
   for i := 3; i <= lim; i += 2 {
       if s%i == 0 {
           return false
       }
   }
   return true
}

func dfs(u, f int) int {
   if f == 0 || u == snk {
       return f
   }
   if vis[u] == visID {
       return 0
   }
   vis[u] = visID
   for _, v := range adj[u] {
       if cap := flow[u][v]; cap > 0 {
           x := dfs(v, min(f, cap))
           if x > 0 {
               flow[u][v] -= x
               flow[v][u] += x
               return x
           }
       }
   }
   return 0
}

func maxFlow() int {
   total := 0
   for {
       visID++
       pushed := dfs(src, n)
       if pushed == 0 {
           break
       }
       total += pushed
   }
   return total
}

func get(u int) {
   tmp = append(tmp, u)
   vis[u] = visID
   for _, v := range adj[u] {
       if vis[v] == visID || v == src || v == snk {
           continue
       }
       if a[u]%2 == flow[v][u] {
           get(v)
       }
   }
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   fmt.Fscan(in, &n)
   a = make([]int, n+2)
   adj = make([][]int, n+2)
   flow = make([][]int, n+2)
   vis = make([]int, n+2)
   src = 0
   snk = n + 1
   odd, even := 0, 0
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
       if a[i]&1 == 1 {
           odd++
       } else {
           even++
       }
   }
   if odd != even {
       fmt.Fprintln(out, "Impossible")
       return
   }
   for i := 0; i <= n+1; i++ {
       flow[i] = make([]int, n+2)
   }
   // build graph
   for i := 1; i <= n; i++ {
       for j := i + 1; j <= n; j++ {
           if isPrimeSum(a[i], a[j]) {
               o, e := i, j
               if a[o]&1 == 0 {
                   o, e = j, i
               }
               adj[o] = append(adj[o], e)
               adj[e] = append(adj[e], o)
               flow[o][e] = 1
           }
       }
   }
   // src and snk edges
   for i := 1; i <= n; i++ {
       if a[i]&1 == 1 {
           adj[src] = append(adj[src], i)
           flow[src][i] = 2
       } else {
           adj[i] = append(adj[i], snk)
           flow[i][snk] = 2
       }
   }
   if maxFlow() != n {
       fmt.Fprintln(out, "Impossible")
       return
   }
   // extract cycles
   visID++
   for i := 1; i <= n; i++ {
       if vis[i] != visID {
           tmp = nil
           get(i)
           ans = append(ans, append([]int(nil), tmp...))
       }
   }
   // output
   fmt.Fprintln(out, len(ans))
   for _, comp := range ans {
       fmt.Fprint(out, len(comp))
       for _, x := range comp {
           fmt.Fprint(out, " ", x)
       }
       fmt.Fprintln(out)
   }
}
