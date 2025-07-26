package main

import (
   "bufio"
   "fmt"
   "os"
   "math/bits"
)

func main() {
   const MOD = 998244353
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   u := make([]int, m)
   v := make([]int, m)
   undEdges := make([][2]int, m)
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(in, &a, &b)
       a--;
       b--;
       u[i] = a; v[i] = b
       undEdges[i] = [2]int{a, b}
   }
   // no edges
   if m == 0 {
       fmt.Println(0)
       return
   }
   // check if undirected graph is acyclic (forest)
   parent := make([]int, n)
   for i := range parent { parent[i] = i }
   var ufFind func(int) int
   ufFind = func(x int) int {
       if parent[x] != x { parent[x] = ufFind(parent[x]) }
       return parent[x]
   }
   acyclicUnd := true
   for i := 0; i < m; i++ {
       a, b := undEdges[i][0], undEdges[i][1]
       ra, rb := ufFind(a), ufFind(b)
       if ra == rb {
           acyclicUnd = false
       } else {
           parent[ra] = rb
       }
   }
   // if forest: all orientations acyclic => sum costs = m * 2^{m-1}
   if acyclicUnd {
       pow := 1
       for i := 0; i < m-1; i++ {
           pow = pow * 2 % MOD
       }
       res := int64(m) * int64(pow) % MOD
       fmt.Println(res)
       return
   }
   // brute for small m
   if m <= 18 {
       masksz := 1 << m
       total := int64(0)
       for mask := 0; mask < masksz; mask++ {
           cost := bits.OnesCount(uint(mask))
           // build graph
           g := make([][]int, n)
           for i := 0; i < m; i++ {
               if mask>>i&1 == 0 {
                   g[u[i]] = append(g[u[i]], v[i])
               } else {
                   g[v[i]] = append(g[v[i]], u[i])
               }
           }
           // check acyclic
           vis := make([]int, n)
           ok := true
           var dfs func(int)
           dfs = func(x int) {
               vis[x] = 1
               for _, w := range g[x] {
                   if vis[w] == 1 {
                       ok = false
                       return
                   }
                   if vis[w] == 0 {
                       dfs(w)
                       if !ok { return }
                   }
               }
               vis[x] = 2
           }
           for i := 0; i < n && ok; i++ {
               if vis[i] == 0 {
                   dfs(i)
               }
           }
           if ok {
               total = (total + int64(cost)) % MOD
           }
       }
       fmt.Println(total)
       return
   }
   // DP over independent sets for general case
   fullMask := (1 << n) - 1
   neighbor := make([]int, n)
   adjRev := make([]int, n)
   for i := 0; i < m; i++ {
       a, b := u[i], v[i]
       // undirected
       neighbor[a] |= 1 << b
       neighbor[b] |= 1 << a
       // rev edges: orig b->a
       adjRev[b] |= 1 << a
   }
   sz := 1 << n
   independent := make([]bool, sz)
   independent[0] = true
   for mask := 1; mask < sz; mask++ {
       lb := mask & -mask
       v0 := bits.TrailingZeros(uint(lb))
       prev := mask ^ lb
       independent[mask] = independent[prev] && ((neighbor[v0] & prev) == 0)
   }
   dpWays := make([]int64, sz)
   dpCost := make([]int64, sz)
   dpWays[0] = 1
   for mask := 1; mask < sz; mask++ {
       var wsum, csum int64
       // enumerate independent non-empty subsets D of mask
       for sub := mask; sub > 0; sub = (sub - 1) & mask {
           if !independent[sub] {
               continue
           }
           prev := mask ^ sub
           w := dpWays[prev]
           if w == 0 {
               continue
           }
           // cost flips: edges orig prev->sub must flip
           cf := 0
           d := sub
           for d > 0 {
               lb2 := d & -d
               i0 := bits.TrailingZeros(uint(lb2))
               cf += bits.OnesCount(uint(adjRev[i0] & prev))
               d ^= lb2
           }
           wsum = (wsum + w) % MOD
           csum = (csum + dpCost[prev] + w*int64(cf)) % MOD
       }
       dpWays[mask] = wsum
       dpCost[mask] = csum
   }
   fmt.Println(dpCost[fullMask])
}
