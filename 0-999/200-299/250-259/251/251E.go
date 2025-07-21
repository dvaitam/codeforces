package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func add(a, b int) int {
   a += b
   if a >= MOD {
       a -= MOD
   }
   return a
}

func mul(a, b int) int {
   return int((int64(a) * int64(b)) % MOD)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   m := 2 * n
   adj := make([][]int, m)
   deg := make([]int, m)
   for i := 0; i < m-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       u--; v--
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
       deg[u]++
       deg[v]++
   }
   if n == 1 {
       // only one edge, two ways to place vertically
       fmt.Println(2)
       return
   }
   // Prepare tree order
   parent := make([]int, m)
   order := make([]int, 0, m)
   parent[0] = -1
   // BFS or stack
   stack := []int{0}
   for len(stack) > 0 {
       u := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       order = append(order, u)
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           parent[v] = u
           stack = append(stack, v)
       }
   }
   // dp0[u], dp1[u] - each [8] ints: index = b*4 + k (k capped at 3)
   dp0 := make([][8]int, m)
   dp1 := make([][8]int, m)
   // Process in reverse order (post-order)
   for i := len(order) - 1; i >= 0; i-- {
       u := order[i]
       // children
       var children []int
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           children = append(children, v)
       }
       csz := len(children)
       // prefix and suffix products for dp[c][0]
       // A[j] = product dp0 of children[0:j]
       A := make([][8]int, csz+1)
       B := make([][8]int, csz+1)
       // init A[0] = identity (b=0,k=0)
       A[0][0] = 1
       for j := 0; j < csz; j++ {
           // A[j+1] = A[j] * dp0[child[j]]
           var tmp [8]int
           for s1 := 0; s1 < 8; s1++ {
               if A[j][s1] == 0 {
                   continue
               }
               b1 := s1 >> 2; k1 := s1 & 3
               for s2 := 0; s2 < 8; s2++ {
                   v := dp0[children[j]][s2]
                   if v == 0 {
                       continue
                   }
                   b2 := s2 >> 2; k2 := s2 & 3
                   nb := b1 | b2
                   nk := k1 + k2
                   if nk > 3 {
                       nk = 3
                   }
                   idx := (nb<<2) | nk
                   tmp[idx] = (tmp[idx] + int((int64(A[j][s1]) * int64(v)) % MOD)) % MOD
               }
           }
           A[j+1] = tmp
       }
       // init B[csz] = identity
       B[csz][0] = 1
       for j := csz - 1; j >= 0; j-- {
           var tmp [8]int
           for s1 := 0; s1 < 8; s1++ {
               if B[j+1][s1] == 0 {
                   continue
               }
               b1 := s1 >> 2; k1 := s1 & 3
               for s2 := 0; s2 < 8; s2++ {
                   v := dp0[children[j]][s2]
                   if v == 0 {
                       continue
                   }
                   b2 := s2 >> 2; k2 := s2 & 3
                   nb := b1 | b2
                   nk := k1 + k2
                   if nk > 3 {
                       nk = 3
                   }
                   idx := (nb<<2) | nk
                   tmp[idx] = (tmp[idx] + int((int64(B[j+1][s1]) * int64(v)) % MOD)) % MOD
               }
           }
           B[j] = tmp
       }
       // dp1[u] = A[csz]
       dp1[u] = A[csz]
       // dp0[u]: match u with exactly one child
       var res0 [8]int
       for j, v := range children {
           // combine others = A[j] * B[j+1]
           var others [8]int
           for s1 := 0; s1 < 8; s1++ {
               if A[j][s1] == 0 {
                   continue
               }
               b1 := s1 >> 2; k1 := s1 & 3
               for s2 := 0; s2 < 8; s2++ {
                   c2 := B[j+1][s2]
                   if c2 == 0 {
                       continue
                   }
                   b2 := s2 >> 2; k2 := s2 & 3
                   nb := b1 | b2
                   nk := k1 + k2
                   if nk > 3 {
                       nk = 3
                   }
                   idx := (nb<<2) | nk
                   others[idx] = (others[idx] + int((int64(A[j][s1]) * int64(c2)) % MOD)) % MOD
               }
           }
           // now match u-v: use dp1[v]
           // edge u-v has c_e = deg[u]+deg[v]-2
           ce := deg[u] + deg[v] - 2
           for s := 0; s < 8; s++ {
               if others[s] == 0 {
                   continue
               }
               b0 := s >> 2; k0 := s & 3
               for s2 := 0; s2 < 8; s2++ {
                   cnt2 := dp1[v][s2]
                   if cnt2 == 0 {
                       continue
                   }
                   b2 := s2 >> 2; k2 := s2 & 3
                   nb := b0 | b2
                   nk := k0 + k2
                   if nk > 3 {
                       nk = 3
                   }
                   // apply effect of ce
                   if ce == 0 {
                       nb = 1
                   } else if ce == 1 {
                       nk++
                       if nk > 3 {
                           nk = 3
                       }
                   }
                   idx := (nb<<2) | nk
                   // multiply others[s] * cnt2
                   res0[idx] = (res0[idx] + int((int64(others[s]) * int64(cnt2)) % MOD)) % MOD
               }
           }
       }
       dp0[u] = res0
   }
   // root u=0, need dp0[0] with b=0 and k=2
   ans := dp0[0][0*4+2]
   // multiply by 4 for placements orientation and row flip
   ans = int((int64(ans) * 4) % MOD)
   fmt.Println(ans)
}
