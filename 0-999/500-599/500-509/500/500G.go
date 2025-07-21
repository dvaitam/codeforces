package main

import (
   "bufio"
   "fmt"
   "os"
)

const MAXN = 200005

var (
   n int
   adj [MAXN][]int
   up  [20][MAXN]int
   depth [MAXN]int
   tin, tout [MAXN]int
   timer int
)

func dfs(u, p int) {
   tin[u] = timer; timer++
   up[0][u] = p
   for i := 1; i < 20; i++ {
       up[i][u] = up[i-1][ up[i-1][u] ]
   }
   for _, v := range adj[u] {
       if v == p { continue }
       depth[v] = depth[u] + 1
       dfs(v, u)
   }
   tout[u] = timer; timer++
}

func isAncestor(u, v int) bool {
   return tin[u] <= tin[v] && tout[v] <= tout[u]
}

func lca(u, v int) int {
   if isAncestor(u, v) {
       return u
   }
   if isAncestor(v, u) {
       return v
   }
   for i := 19; i >= 0; i-- {
       if up[i][u] != 0 && !isAncestor(up[i][u], v) {
           u = up[i][u]
       }
   }
   return up[0][u]
}

func dist(u, v int) int {
   w := lca(u, v)
   return depth[u] + depth[v] - 2*depth[w]
}

// get k-th node on path u->v, 0-based (0 gives u)
func getKth(u, v, k int) int {
   w := lca(u, v)
   duw := depth[u] - depth[w]
   if k <= duw {
       // go up from u by k
       x := u
       for i := 0; i < 20; i++ {
           if k>>i & 1 == 1 {
               x = up[i][x]
           }
       }
       return x
   }
   // go down from w towards v
   kv := depth[v] - depth[w] - (k - duw)
   // move up from v by kv
   x := v
   for i := 0; i < 20; i++ {
       if kv>>i & 1 == 1 {
           x = up[i][x]
       }
   }
   return x
}

// extended gcd: returns g,x,y such that ax+by=g
func extgcd(a, b int) (g, x, y int) {
   if b == 0 {
       return a, 1, 0
   }
   g, x1, y1 := extgcd(b, a%b)
   return g, y1, x1 - (a/b)*y1
}

// solve t ≡ a mod n, t ≡ b mod m; return (t mod lcm), or (-1) if none
func crt(a, n, b, m int) int64 {
   // a + n * x = b (mod m)
   g, x, y := extgcd(n, m)
   if (b - a) % g != 0 {
       return -1
   }
   lcm := int64(n/g) * int64(m)
   // multiply x by (b-a)/g
   mul := (b - a) / g
   x0 := int64(x) * int64(mul) % int64(m/g)
   t := (int64(a) + int64(n)*x0) % lcm
   if t < 0 {
       t += lcm
   }
   return t
}

func main() {
   in := bufio.NewReader(os.Stdin)
   fmt.Fscan(in, &n)
   for i := 1; i < n; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // build depth, parent, tin/tout via iterative DFS to avoid recursion
   timer = 0
   depth[1] = 0
   type frame struct{u, p, i int}
   stack := []frame{{1, 0, 0}}
   for len(stack) > 0 {
       fr := &stack[len(stack)-1]
       u, p := fr.u, fr.p
       if fr.i == 0 {
           tin[u] = timer; timer++
           up[0][u] = p
           for k := 1; k < 20; k++ {
               up[k][u] = up[k-1][ up[k-1][u] ]
           }
       }
       // process children
       if fr.i < len(adj[u]) {
           v := adj[u][fr.i]
           fr.i++
           if v == p {
               continue
           }
           depth[v] = depth[u] + 1
           stack = append(stack, frame{v, u, 0})
           continue
       }
       // exiting u
       tout[u] = timer; timer++
       stack = stack[:len(stack)-1]
   }
   var t int
   fmt.Fscan(in, &t)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for tc := 0; tc < t; tc++ {
       var u, v, x, y int
       fmt.Fscan(in, &u, &v, &x, &y)
       // compute intersection endpoints
       cand := make([]int, 0, 4)
       for _, p := range []int{lca(u, x), lca(u, y), lca(v, x), lca(v, y)} {
           if dist(u, v) == dist(u, p) + dist(p, v) && dist(x, y) == dist(x, p) + dist(p, y) {
               cand = append(cand, p)
           }
       }
       if len(cand) == 0 {
           fmt.Fprintln(out, -1)
           continue
       }
       // find two endpoints s,t with max distance
       s, tnode := cand[0], cand[0]
       maxd := 0
       for i := 0; i < len(cand); i++ {
           for j := i; j < len(cand); j++ {
               d := dist(cand[i], cand[j])
               if d > maxd {
                   maxd = d
                   s, tnode = cand[i], cand[j]
               }
           }
       }
       // precompute periods and endpoint distances
       LA := dist(u, v)
       LB := dist(x, y)
       P1 := LA * 2
       P2 := LB * 2
       best := int64(-1)
       for _, c := range []int{s, tnode} {
           d1 := dist(u, c)
           d2 := dist(x, c)
           rs := []int{d1, (P1 - d1) % P1}
           rt := []int{d2, (P2 - d2) % P2}
           for _, r1 := range rs {
               for _, r2 := range rt {
                   t0 := crt(r1, P1, r2, P2)
                   if t0 >= 0 {
                       if best < 0 || t0 < best {
                           best = t0
                       }
                   }
               }
           }
       }
       if best < 0 {
           fmt.Fprintln(out, -1)
       } else {
           fmt.Fprintln(out, best)
       }
   }
}
