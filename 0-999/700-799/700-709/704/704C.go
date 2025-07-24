package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   // unary flags
   hasPos := make([]bool, m+1)
   hasNeg := make([]bool, m+1)
   // binary edges
   type edge struct{ u, v int; su, sv int }
   edges := make([]edge, 0, n)
   adj := make([][]int, m+1)
   globalConst := 0
   for i := 0; i < n; i++ {
       var k int
       fmt.Fscan(reader, &k)
       if k == 1 {
           var a int
           fmt.Fscan(reader, &a)
           v := abs(a)
           if a > 0 {
               hasPos[v] = !hasPos[v]
           } else {
               hasNeg[v] = !hasNeg[v]
           }
       } else {
           var a, b int
           fmt.Fscan(reader, &a, &b)
           u := abs(a)
           v := abs(b)
           su := 0; if a < 0 { su = 1 }
           sv := 0; if b < 0 { sv = 1 }
           if u == v {
               if su == sv {
                   // OR(x,x) = x
                   if su == 0 {
                       hasPos[u] = !hasPos[u]
                   } else {
                       hasNeg[u] = !hasNeg[u]
                   }
               } else {
                   // OR(x, !x) = 1
                   globalConst ^= 1
               }
           } else {
               idx := len(edges)
               edges = append(edges, edge{u, v, su, sv})
               adj[u] = append(adj[u], idx)
               adj[v] = append(adj[v], idx)
           }
       }
   }
   // process unary per var
   unary := make([]int, m+1)
   for i := 1; i <= m; i++ {
       sp := 0; if hasPos[i] { sp = 1 }
       sq := 0; if hasNeg[i] { sq = 1 }
       if sp == 1 && sq == 1 {
           globalConst ^= 1
           unary[i] = 0
       } else if sp+sq == 1 {
           if sp == 1 {
               unary[i] = 1 // Ci = xi
           } else {
               unary[i] = 2 // Ci = 1-xi
           }
       }
   }
   // visited for binary components
   visited := make([]bool, m+1)
   ways0, ways1 := int64(1), int64(0)
   // helper for unary term
   unaryTerm := func(i, x int) int {
       switch unary[i] {
       case 1:
           return x
       case 2:
           return 1 - x
       }
       return 0
   }
   // use global abs function for absolute value
   // iterate vars
   for i := 1; i <= m; i++ {
       if visited[i] || len(adj[i]) == 0 {
           continue
       }
       // collect component nodes
       stack := []int{i}
       comp := []int(nil)
       visited[i] = true
       for len(stack) > 0 {
           u := stack[len(stack)-1]; stack = stack[:len(stack)-1]
           comp = append(comp, u)
           for _, ei := range adj[u] {
               e := edges[ei]
               v := e.u ^ e.v ^ u
               if !visited[v] {
                   visited[v] = true
                   stack = append(stack, v)
               }
           }
       }
       // find start: node with deg1 or comp[0]
       start := comp[0]
       for _, u := range comp {
           if len(adj[u]) == 1 {
               start = u; break
           }
       }
       // traverse to get order
       var nodes []int
       var spPrev, spNext []int
       prev := -1; cur := start
       // determine if cycle
       isCycle := len(adj[start]) == 2
       for {
           nodes = append(nodes, cur)
           var found bool
           for _, ei := range adj[cur] {
               e := edges[ei]
               v := e.u ^ e.v ^ cur
               if v == prev {
                   continue
               }
               // choose this edge
               // record signs
               var s1, s2 int
               if cur == e.u {
                   s1, s2 = e.su, e.sv
               } else {
                   s1, s2 = e.sv, e.su
               }
               spPrev = append(spPrev, s1)
               spNext = append(spNext, s2)
               prev, cur = cur, v
               found = true
               break
           }
           if !found {
               break
           }
           if isCycle && cur == start {
               break
           }
       }
       k := len(nodes)
       var comp0, comp1 int64
       // path
       if !isCycle {
           // dp[x][p]
           dp0 := [2][2]int64{}
           // init
           for x0 := 0; x0 < 2; x0++ {
               p0 := unaryTerm(nodes[0], x0)
               dp0[x0][p0] = 1
           }
           // transitions
           for i2 := 0; i2 < k-1; i2++ {
               var dp1 [2][2]int64
               for xi := 0; xi < 2; xi++ {
                   for p := 0; p < 2; p++ {
                       cnt := dp0[xi][p]
                       if cnt == 0 {
                           continue
                       }
                       for xj := 0; xj < 2; xj++ {
                           // clause OR
                           v1 := xi
                           if spPrev[i2] == 1 {
                               v1 = 1 - v1
                           }
                           v2 := xj
                           if spNext[i2] == 1 {
                               v2 = 1 - v2
                           }
                           c := v1 | v2
                           u2 := unaryTerm(nodes[i2+1], xj)
                           q := c ^ u2
                           np := p ^ q
                           dp1[xj][np] = (dp1[xj][np] + cnt) % MOD
                       }
                   }
               }
               dp0 = dp1
           }
           // collect
           for xk := 0; xk < 2; xk++ {
               comp0 = (comp0 + dp0[xk][0]) % MOD
               comp1 = (comp1 + dp0[xk][1]) % MOD
           }
       } else {
           // cycle: use separate for each x0
           for x0 := 0; x0 < 2; x0++ {
               // dp[x][p]
               dp0 := [2][2]int64{}
               p0 := unaryTerm(nodes[0], x0)
               dp0[x0][p0] = 1
               // transitions except last
               for i2 := 0; i2 < k-1; i2++ {
                   var dp1 [2][2]int64
                   for xi := 0; xi < 2; xi++ {
                       for p := 0; p < 2; p++ {
                           cnt := dp0[xi][p]
                           if cnt == 0 {
                               continue
                           }
                           for xj := 0; xj < 2; xj++ {
                               v1 := xi
                               if spPrev[i2] == 1 { v1 = 1 - v1 }
                               v2 := xj
                               if spNext[i2] == 1 { v2 = 1 - v2 }
                               c := v1 | v2
                               u2 := unaryTerm(nodes[i2+1], xj)
                               q := c ^ u2
                               np := p ^ q
                               dp1[xj][np] = (dp1[xj][np] + cnt) % MOD
                           }
                       }
                   }
                   dp0 = dp1
               }
               // last edge between last node and start
               i2 := k-1
               for xk := 0; xk < 2; xk++ {
                   for p := 0; p < 2; p++ {
                       cnt := dp0[xk][p]
                       if cnt == 0 { continue }
                       v1 := xk
                       if spPrev[i2] == 1 { v1 = 1 - v1 }
                       v2 := x0
                       if spNext[i2] == 1 { v2 = 1 - v2 }
                       c := v1 | v2
                       np := p ^ c
                       if np == 0 {
                           comp0 = (comp0 + cnt) % MOD
                       } else {
                           comp1 = (comp1 + cnt) % MOD
                       }
                   }
               }
           }
       }
       // combine
       nw0 := (ways0*comp0 + ways1*comp1) % MOD
       nw1 := (ways0*comp1 + ways1*comp0) % MOD
       ways0, ways1 = nw0, nw1
   }
   // handle isolated vars
   for i := 1; i <= m; i++ {
       if len(adj[i]) == 0 {
           // comp for single var
           var c0, c1 int64
           switch unary[i] {
           case 0:
               c0 = 2; c1 = 0
           case 1, 2:
               c0 = 1; c1 = 1
           }
           // combine
           nw0 := (ways0*c0 + ways1*c1) % MOD
           nw1 := (ways0*c1 + ways1*c0) % MOD
           ways0, ways1 = nw0, nw1
       }
   }
   // target parity
   target := globalConst ^ 1
   if target == 0 {
       fmt.Println(ways0)
   } else {
       fmt.Println(ways1)
   }
}
// helper abs
func abs(x int) int { if x<0 { return -x }; return x }
