package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, M int
   adj [][]edge
   used []bool
   subSize []int
   p10, invP10 []int
   inv10 int
   ans int64
)

type edge struct{
   to, d int
}

func modInv(a, m int) int {
   b, u, v := m, 1, 0
   for b != 0 {
       t := a / b
       a, b = b, a - t*b
       u, v = v, u - t*v
   }
   u %= m
   if u < 0 {
       u += m
   }
   return u
}

func addMod(x, y int) int {
   r := x + y
   if r >= M {
       r -= M
   }
   return r
}
func mulMod(x, y int) int {
   return int((int64(x) * int64(y)) % int64(M))
}

// compute subtree sizes
func dfsSize(u, p int) {
   subSize[u] = 1
   for _, e := range adj[u] {
       v := e.to
       if v != p && !used[v] {
           dfsSize(v, u)
           subSize[u] += subSize[v]
       }
   }
}

// find centroid
func dfsCentroid(u, p, tot int) int {
   for _, e := range adj[u] {
       v := e.to
       if v != p && !used[v] && subSize[v] > tot/2 {
           return dfsCentroid(v, u, tot)
       }
   }
   return u
}

type nodeData struct{ f, b, depth int }

// collect f,b,depth for nodes in subtree
func dfsCollect(u, p, depth, f, b int, arr *[]nodeData) {
   *arr = append(*arr, nodeData{f, b, depth})
   for _, e := range adj[u] {
       v := e.to
       if v != p && !used[v] {
           nf := addMod(mulMod(f, 10), e.d)
           nb := addMod(mulMod(e.d, p10[depth]), b)
           dfsCollect(v, u, depth+1, nf, nb, arr)
       }
   }
}

func decompose(u int) {
   dfsSize(u, -1)
   c := dfsCentroid(u, -1, subSize[u])
   used[c] = true
   // maps for B_u and transformed F_u
   cntB := make(map[int]int)
   cntF := make(map[int]int)
   // include centroid itself
   cntB[0] = 1
   cntF[0] = 1
   for _, e := range adj[c] {
       v := e.to
       if used[v] {
           continue
       }
       var arr []nodeData
       // start from neighbor edge
       dfsCollect(v, c, 1, e.d % M, e.d % M, &arr)
       // count pairs
       for _, nd := range arr {
           // (u,v)
           needB := mulMod((M-nd.f)%M, invP10[nd.depth])
           if cnt, ok := cntB[needB]; ok {
               ans += int64(cnt)
           }
           // (v,u)
           needF := (M - nd.b) % M
           if cnt, ok := cntF[needF]; ok {
               ans += int64(cnt)
           }
       }
       // add to maps
       for _, nd := range arr {
           cntB[nd.b]++
           keyF := mulMod(nd.f, invP10[nd.depth])
           cntF[keyF]++
       }
   }
   // recurse
   for _, e := range adj[c] {
       if !used[e.to] {
           decompose(e.to)
       }
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   fmt.Fscan(in, &n, &M)
   adj = make([][]edge, n)
   used = make([]bool, n)
   subSize = make([]int, n)
   for i := 0; i < n-1; i++ {
       var u, v, w int
       fmt.Fscan(in, &u, &v, &w)
       adj[u] = append(adj[u], edge{v, w})
       adj[v] = append(adj[v], edge{u, w})
   }
   // precompute powers
   p10 = make([]int, n+1)
   invP10 = make([]int, n+1)
   p10[0] = 1
   for i := 1; i <= n; i++ {
       p10[i] = mulMod(p10[i-1], 10)
   }
   inv10 = modInv(10, M)
   invP10[0] = 1
   for i := 1; i <= n; i++ {
       invP10[i] = mulMod(invP10[i-1], inv10)
   }
   // decompose and solve
   decompose(0)
   fmt.Println(ans)
}
