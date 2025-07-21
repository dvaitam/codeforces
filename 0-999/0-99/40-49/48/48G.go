package main

import (
   "bufio"
   "fmt"
   "os"
)

// Edge represents a graph edge with weight
type Edge struct { to, w int }

var (
   n            int
   adj          [][]Edge
   deg          []int
   inCycle      []bool
   sizeDown     []int           // size of subtree for dfs
   sumDown      []int64         // sum distances from node to descendants
   sumDepth     []int64         // sum of depths (from cycle root) in subtree
   depthDist    []int64         // distance from cycle root
   cycleNodes   []int           // ordered cycle nodes
   cycleEdgeW   []int           // weights on cycle edges
   tMass        []int64         // mass at cycle positions
   sumCycleDist []int64         // sum of distances on cycle for each cycle node
   ans          []int64         // final answers
)

func main() {
   in := bufio.NewReader(os.Stdin)
   fmt.Fscan(in, &n)
   adj = make([][]Edge, n+1)
   deg = make([]int, n+1)
   // read graph
   for i := 0; i < n; i++ {
       var a, b, t int
       fmt.Fscan(in, &a, &b, &t)
       adj[a] = append(adj[a], Edge{b, t})
       adj[b] = append(adj[b], Edge{a, t})
       deg[a]++
       deg[b]++
   }
   // find cycle nodes by peeling leaves
   inCycle = make([]bool, n+1)
   removed := make([]bool, n+1)
   queue := make([]int, 0, n)
   for u := 1; u <= n; u++ {
       if deg[u] == 1 {
           removed[u] = true
           queue = append(queue, u)
       }
   }
   for i := 0; i < len(queue); i++ {
       u := queue[i]
       for _, e := range adj[u] {
           v := e.to
           if removed[v] {
               continue
           }
           deg[v]--
           if deg[v] == 1 {
               removed[v] = true
               queue = append(queue, v)
           }
       }
   }
   for u := 1; u <= n; u++ {
       if !removed[u] {
           inCycle[u] = true
       }
   }
   // extract cycle in order
   var start int
   for u := 1; u <= n; u++ {
       if inCycle[u] {
           start = u
           break
       }
   }
   // walk the cycle
   prev := -1
   u := start
   for {
       cycleNodes = append(cycleNodes, u)
       // find next cycle neighbor
       var next, w int
       for _, e := range adj[u] {
           if !inCycle[e.to] || e.to == prev {
               continue
           }
           next, w = e.to, e.w
           break
       }
       cycleEdgeW = append(cycleEdgeW, w)
       if next == start {
           break
       }
       prev, u = u, next
   }
   k := len(cycleNodes)
   // prepare DP arrays
   sizeDown = make([]int, n+1)
   sumDown = make([]int64, n+1)
   sumDepth = make([]int64, n+1)
   depthDist = make([]int64, n+1)
   ans = make([]int64, n+1)
   tMass = make([]int64, k)
   // compute subtree stats for each cycle node
   for i, r := range cycleNodes {
       depthDist[r] = 0
       dfs1(r, -1)
       tMass[i] = int64(sizeDown[r])
   }
   // compute cycle distance contributions
   sumCycleDist = make([]int64, k)
   // prefix distances along cycle
   S := make([]int64, k+1)
   for i := 0; i < k; i++ {
       S[i+1] = S[i] + int64(cycleEdgeW[i])
   }
   C := S[k]
   // build doubled arrays
   pos := make([]int64, 2*k)
   mass := make([]int64, 2*k)
   for i := 0; i < 2*k; i++ {
       idx := i % k
       pos[i] = S[idx]
       if i >= k {
           pos[i] += C
       }
       mass[i] = tMass[idx]
   }
   // prefix sums
   Pm := make([]int64, 2*k+1)
   Pp := make([]int64, 2*k+1)
   for i := 0; i < 2*k; i++ {
       Pm[i+1] = Pm[i] + mass[i]
       Pp[i+1] = Pp[i] + mass[i]*pos[i]
   }
   // two-pointer over cycle
   R := 0
   for i := 0; i < k; i++ {
       if R < i {
           R = i
       }
       for R < i+k && 2*(pos[R]-pos[i]) <= C {
           R++
       }
       // [i..R-1] close forward, [R..i+k-1] backward
       sumLmass := Pm[R] - Pm[i]
       sumLpm := Pp[R] - Pp[i]
       sumL := sumLpm - pos[i]*sumLmass
       sumRmass := Pm[i+k] - Pm[R]
       sumRpm := Pp[i+k] - Pp[R]
       sumR := (C + pos[i])*sumRmass - sumRpm
       sumCycleDist[i] = sumL + sumR
   }
   // assign answers for cycle nodes
   for i, u := range cycleNodes {
       ans[u] = sumDown[u] + sumCycleDist[i]
   }
   // fill answers for tree nodes
   for i, r := range cycleNodes {
       dfsAssign(r, -1, i)
   }
   // output
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for i := 1; i <= n; i++ {
       if i > 1 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, ans[i])
   }
   fmt.Fprintln(out)
}

// dfs1 computes sizeDown, sumDown, sumDepth, depthDist for tree rooted at u
func dfs1(u, p int) {
   sizeDown[u] = 1
   sumDown[u] = 0
   sumDepth[u] = depthDist[u]
   for _, e := range adj[u] {
       v, w := e.to, e.w
       if v == p || inCycle[v] {
           continue
       }
       depthDist[v] = depthDist[u] + int64(w)
       dfs1(v, u)
       sizeDown[u] += sizeDown[v]
       sumDown[u] += sumDown[v] + int64(sizeDown[v])*int64(w)
       sumDepth[u] += sumDepth[v]
   }
}

// dfsAssign computes ans for tree nodes under cycle node indexed ci
func dfsAssign(u, p, ci int) {
   root := cycleNodes[ci]
   for _, e := range adj[u] {
       v, _ := e.to, e.w
       if v == p || inCycle[v] {
           continue
       }
       // nodes outside v's subtree: rem
       rem := int64(n - sizeDown[v])
       ans[v] = sumDown[v] + rem*depthDist[v] + ans[root] - sumDepth[v]
       dfsAssign(v, u, ci)
   }
}
