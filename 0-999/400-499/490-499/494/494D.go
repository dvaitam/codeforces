package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

type Edge struct { to, w int }

var (
   n int
   adj [][]Edge
   D []int64
   sizeSub []int
   sumDSub, sumD2Sub []int64
   G, T []int64
   parent [][]int
   level []int
   tin, tout []int
   timer int
)

func add(a, b int64) int64 { a += b; if a >= MOD { a -= MOD }; return a }
func sub(a, b int64) int64 { a -= b; if a < 0 { a += MOD }; return a }
func mul(a, b int64) int64 { return (a % MOD) * (b % MOD) % MOD }

func dfs1(u, p int) {
   tin[u] = timer
   timer++
   sizeSub[u] = 1
   sumDSub[u] = D[u]
   sumD2Sub[u] = mul(D[u], D[u])
   parent[0][u] = p
   for _, e := range adj[u] {
       v := e.to
       if v == p {
           continue
       }
       D[v] = add(D[u], int64(e.w))
       level[v] = level[u] + 1
       dfs1(v, u)
       sizeSub[u] += sizeSub[v]
       sumDSub[u] = add(sumDSub[u], sumDSub[v])
       sumD2Sub[u] = add(sumD2Sub[u], sumD2Sub[v])
   }
   tout[u] = timer - 1
}

func dfs2(u, p int) {
   for _, e := range adj[u] {
       v := e.to
       if v == p {
           continue
       }
       w := int64(e.w) % MOD
       // G[v] = G[u] + w*(n - 2*sizeSub[v])
       G[v] = add(G[u], mul(w, int64(n-2*sizeSub[v])))
       // sumDistSub = sum (D[x] - D[u]) over x in sub(v)
       sumDistSub := sub(sumDSub[v], mul(int64(sizeSub[v]), D[u]))
       // T[v] = T[u] + 2*w*(G[u] - 2*sumDistSub) + n*w*w
       tmp := sub(G[u], mul(2, sumDistSub))
       T[v] = add(add(T[u], mul(mul(2, w), tmp)), mul(int64(n), mul(w, w)))
       dfs2(v, u)
   }
}

func isInSub(u, v int) bool {
   return tin[v] <= tin[u] && tin[u] <= tout[v]
}

func lca(u, v int) int {
   if level[u] < level[v] {
       u, v = v, u
   }
   diff := level[u] - level[v]
   for k := 0; diff > 0; k++ {
       if diff&1 != 0 {
           u = parent[k][u]
       }
       diff >>= 1
   }
   if u == v {
       return u
   }
   for k := len(parent) - 1; k >= 0; k-- {
       pu := parent[k][u]
       pv := parent[k][v]
       if pu != pv {
           u = pu
           v = pv
       }
   }
   return parent[0][u]
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n)
   adj = make([][]Edge, n+1)
   for i := 0; i < n-1; i++ {
       var a, b, c int
       fmt.Fscan(in, &a, &b, &c)
       adj[a] = append(adj[a], Edge{b, c})
       adj[b] = append(adj[b], Edge{a, c})
   }
   D = make([]int64, n+1)
   sizeSub = make([]int, n+1)
   sumDSub = make([]int64, n+1)
   sumD2Sub = make([]int64, n+1)
   level = make([]int, n+1)
   tin = make([]int, n+1)
   tout = make([]int, n+1)
   // parent
   LOG := 1
   for (1<<LOG) <= n {
       LOG++
   }
   parent = make([][]int, LOG)
   for i := range parent {
       parent[i] = make([]int, n+1)
   }
   timer = 0
   dfs1(1, 0)
   // build parent
   for k := 1; k < LOG; k++ {
       for v := 1; v <= n; v++ {
           parent[k][v] = parent[k-1][ parent[k-1][v] ]
       }
   }
   G = make([]int64, n+1)
   T = make([]int64, n+1)
   G[1] = sumDSub[1] % MOD
   T[1] = sumD2Sub[1] % MOD
   dfs2(1, 0)
   var q int
   fmt.Fscan(in, &q)
   for i := 0; i < q; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       m := int64(sizeSub[v])
       S := sumDSub[v] % MOD
       S2 := sumD2Sub[v] % MOD
       D_u := D[u] % MOD
       var SumSub int64
       if !isInSub(u, v) {
           L := lca(u, v)
           D_L := D[L] % MOD
           // m*(D_u^2 +4D_L^2 -4D_u*D_L) + S2 +2D_u*S -4D_L*S
           part := add( add(mul(D_u, D_u), mul(4, mul(D_L, D_L))), sub(0, mul(4, mul(D_u, D_L))) )
           SumSub = add(add(mul(m, part), S2), add(mul(2, mul(D_u, S)), sub(0, mul(4, S*D_L%MOD))))
       } else {
           // u in sub(v)
           n1 := int64(sizeSub[u])
           Sx := sumDSub[u] % MOD
           S2x := sumD2Sub[u] % MOD
           m2 := m - n1
           S_2 := sub(S, Sx)
           S2_2 := sub(S2, S2x)
           D_v := D[v] % MOD
           // sum1 = S2x -2D_u*Sx + n1*D_u^2
           sum1 := add(sub(S2x, mul(2, mul(D_u, Sx))), mul(n1, mul(D_u, D_u)))
           // sum2 = m2*(D_u -2D_v)^2 + S2_2 +2*(D_u -2D_v)*S_2
           tmp := sub(D_u, mul(2, D_v))
           sum2 := add(add(mul(m2, mul(tmp, tmp)), S2_2), mul(2, mul(tmp, S_2)))
           SumSub = add(sum1, sum2)
       }
       ans := sub(mul(2, SumSub), T[u]%MOD)
       fmt.Fprintln(out, ans)
   }
}
