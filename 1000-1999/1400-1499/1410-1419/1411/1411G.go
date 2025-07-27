package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const mod = 998244353

func add(a, b int) int {
   s := a + b
   if s >= mod {
       s -= mod
   }
   return s
}

func sub(a, b int) int {
   s := a - b
   if s < 0 {
       s += mod
   }
   return s
}

func mul(a, b int) int {
   return int((int64(a) * int64(b)) % mod)
}

func modpow(a, e int) int {
   res := 1
   x := a
   for e > 0 {
       if e&1 != 0 {
           res = mul(res, x)
       }
       x = mul(x, x)
       e >>= 1
   }
   return res
}

func inv(a int) int {
   // mod is prime
   return modpow(a, mod-2)
}

// fwht computes unnormalized xor transform of a in place
func fwht(a []int, n int) {
   for len1 := 1; len1 < n; len1 <<= 1 {
       for i := 0; i < n; i += len1 << 1 {
           for j := 0; j < len1; j++ {
               u := a[i+j]
               v := a[i+j+len1]
               a[i+j] = add(u, v)
               a[i+j+len1] = sub(u, v)
           }
       }
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   fmt.Fscan(in, &n, &m)
   adj := make([][]int, n)
   indeg := make([]int, n)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       u--
       v--
       adj[u] = append(adj[u], v)
       indeg[v]++
   }
   // topo sort
   q := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if indeg[i] == 0 {
           q = append(q, i)
       }
   }
   for i := 0; i < len(q); i++ {
       u := q[i]
       for _, v := range adj[u] {
           indeg[v]--
           if indeg[v] == 0 {
               q = append(q, v)
           }
       }
   }
   // compute grundy in reverse topo
   g := make([]int, n)
   maxg := 0
   tmp := make([]int, 0)
   for idx := n - 1; idx >= 0; idx-- {
       u := q[idx]
       tmp = tmp[:0]
       for _, v := range adj[u] {
           tmp = append(tmp, g[v])
       }
       if len(tmp) > 0 {
           sort.Ints(tmp)
           mex := 0
           for _, x := range tmp {
               if x == mex {
                   mex++
               } else if x > mex {
                   break
               }
           }
           g[u] = mex
       } else {
           g[u] = 0
       }
       if g[u] > maxg {
           maxg = g[u]
       }
   }
   // count frequencies
   size := 1
   for size <= maxg {
       size <<= 1
   }
   V := make([]int, size)
   invn := inv(n)
   for i := 0; i < n; i++ {
       V[g[i]] = add(V[g[i]], invn)
   }
   // FWHT
   fwht(V, size)
   // probabilities
   p_stop := inv(n + 1)
   p_add := mul(n, p_stop)
   // sum terms
   sum := 0
   for i := 0; i < size; i++ {
       // denom = 1 - p_add * V[i]
       denom := sub(1, mul(p_add, V[i]))
       sum = add(sum, inv(denom))
   }
   // E = p_stop * sum * inv(size)
   E := mul(p_stop, mul(sum, inv(size)))
   ans := sub(1, E)
   fmt.Fprintln(out, ans)
}
