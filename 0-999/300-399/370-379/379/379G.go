package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = 1000000000

type Edge struct { u, v int }

var (
   n int
   adj [][]int
   dfn, low []int
   bccStack []Edge
   blocks [][]int
   timeDfs int
   bcAdj [][]int
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func dfsBCC(u, parent int) {
   timeDfs++
   dfn[u] = timeDfs
   low[u] = timeDfs
   for _, v := range adj[u] {
       if v == parent {
           // skip only one parent edge
           parent = -1
           continue
       }
       if dfn[v] == 0 {
           bccStack = append(bccStack, Edge{u, v})
           dfsBCC(v, u)
           low[u] = min(low[u], low[v])
           if low[v] >= dfn[u] {
               // pop edges until u-v
               var comp []int
               seen := make(map[int]bool)
               for {
                   e := bccStack[len(bccStack)-1]
                   bccStack = bccStack[:len(bccStack)-1]
                   if !seen[e.u] {
                       seen[e.u] = true
                       comp = append(comp, e.u)
                   }
                   if !seen[e.v] {
                       seen[e.v] = true
                       comp = append(comp, e.v)
                   }
                   if e.u == u && e.v == v {
                       break
                   }
               }
               blocks = append(blocks, comp)
           }
       } else if dfn[v] < dfn[u] {
           // back edge
           bccStack = append(bccStack, Edge{u, v})
           low[u] = min(low[u], dfn[v])
       }
   }
}

// dfs on block-cut tree, u < n: vertex, u >= n: block node
func dfs(u, parent int) ([]int, []int) {
   if u < n {
       // vertex node
       // dp0, dp1 slices, initial size 2
       dp0 := []int{0, INF}
       dp1 := []int{INF, 0}
       for _, v := range bcAdj[u] {
           if v == parent {
               continue
           }
           c0, c1 := dfs(v, u)
           // merge
           // new sizes
           l0, l1 := len(dp0), len(dp1)
           lc := len(c0)
           ndp0 := make([]int, l0+lc-1)
           ndp1 := make([]int, l1+len(c1)-1)
           for i := range ndp0 {
               ndp0[i] = INF
           }
           for i := range ndp1 {
               ndp1[i] = INF
           }
           // dp0 with c0
           for i := 0; i < l0; i++ {
               if dp0[i] >= INF {
                   continue
               }
               for j := 0; j < lc; j++ {
                   if c0[j] >= INF {
                       continue
                   }
                   v := dp0[i] + c0[j]
                   if v < ndp0[i+j] {
                       ndp0[i+j] = v
                   }
               }
           }
           // dp1 with c1
           for i := 0; i < l1; i++ {
               if dp1[i] >= INF {
                   continue
               }
               for j := 0; j < len(c1); j++ {
                   if c1[j] >= INF {
                       continue
                   }
                   v := dp1[i] + c1[j]
                   if v < ndp1[i+j] {
                       ndp1[i+j] = v
                   }
               }
           }
           dp0, dp1 = ndp0, ndp1
       }
       return dp0, dp1
   }
   // block node
   bid := u - n
   verts := blocks[bid]
   // build set for block
   inB := make(map[int]bool)
   for _, v := range verts {
       inB[v] = true
   }
   // find parent vertex
   rp := parent
   // build cycle order
   // adjacency within block
   neigh := make(map[int][]int)
   for _, v := range verts {
       for _, w := range adj[v] {
           if inB[w] {
               neigh[v] = append(neigh[v], w)
           }
       }
   }
   // sequence arr, start from rp
   m := len(verts)
   arr := make([]int, m)
   arr[0] = rp
   // pick next
   var prev = -1
   if len(neigh[rp]) > 0 {
       // if cycle, degree 2; if edge, degree1
       prev = rp
       arr[1] = neigh[rp][0]
       prev = rp
   }
   // traverse
   for i := 2; i < m; i++ {
       cur := arr[i-1]
       // pick neighbor not equal prev
       for _, w := range neigh[cur] {
           if w != prev {
               arr[i] = w
               prev = cur
               break
           }
       }
   }
   // now arr holds rp, v1..v[m-1]
   // get child dp for positions 1..m-1
   childDP0 := make([][]int, m)
   childDP1 := make([][]int, m)
   for i := 1; i < m; i++ {
       v := arr[i]
       dp0, dp1 := dfs(v, u)
       childDP0[i] = dp0
       childDP1[i] = dp1
   }
   // DP for each c0
   dpB0 := make([]int, m)
   dpB1 := make([]int, m)
   for k := 0; k < m; k++ {
       dpB0[k], dpB1[k] = INF, INF
   }
   // for c0 = 0,1
   for c0 := 0; c0 <= 1; c0++ {
       // cur_dp[c][s] = b
       cur := make([][]int, 2)
       // initial at pos0
       cur[c0] = make([]int, 1)
       cur[1-c0] = make([]int, 0)
       cur[c0][0] = 0
       // iterate positions
       for pos := 1; pos < m; pos++ {
           next := make([][]int, 2)
           for cc := 0; cc <= 1; cc++ {
               next[cc] = make([]int, len(cur[0])+len(cur[1])) // overestimate
               for i := range next[cc] {
                   next[cc][i] = INF
               }
           }
           for prevC := 0; prevC <= 1; prevC++ {
               dpPrev := cur[prevC]
               if len(dpPrev) == 0 {
                   continue
               }
               for currC := 0; currC <= 1; currC++ {
                   dpChild := childDP0[pos]
                   if currC == 1 {
                       dpChild = childDP1[pos]
                   }
                   addE := 0
                   if prevC != currC {
                       addE = 1
                   }
                   for s1, b1 := range dpPrev {
                       if b1 >= INF {
                           continue
                       }
                       for s2, b2 := range dpChild {
                           if b2 >= INF {
                               continue
                           }
                           sNew := s1 + s2
                           v := b1 + b2 + addE
                           if v < next[currC][sNew] {
                               next[currC][sNew] = v
                           }
                       }
                   }
               }
           }
           cur = next
       }
       // close cycle edge
       for lastC := 0; lastC <= 1; lastC++ {
           dpLast := cur[lastC]
           if len(dpLast) == 0 {
               continue
           }
           addE := 0
           if lastC != c0 {
               addE = 1
           }
           for s, b := range dpLast {
               if b >= INF {
                   continue
               }
               totalB := b + addE
               if totalB < []int{dpB0[s], dpB1[s]}[c0] {
                   if c0 == 0 {
                       dpB0[s] = totalB
                   } else {
                       dpB1[s] = totalB
                   }
               }
           }
       }
   }
   return dpB0, dpB1
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var m int
   fmt.Fscan(reader, &n, &m)
   adj = make([][]int, n)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       u--
       v--
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   dfn = make([]int, n)
   low = make([]int, n)
   timeDfs = 0
   dfsBCC(0, -1)
   // build block-cut tree
   B := len(blocks)
   bcAdj = make([][]int, n+B)
   for i, comp := range blocks {
       bid := n + i
       for _, v := range comp {
           bcAdj[bid] = append(bcAdj[bid], v)
           bcAdj[v] = append(bcAdj[v], bid)
       }
   }
   // DP
   dp0, dp1 := dfs(0, -1)
   // output results
   // dp0, dp1 lengths
   res := make([]int, n+1)
   for a := 0; a <= n; a++ {
       b0 := INF
       if a < len(dp0) {
           b0 = dp0[a]
       }
       b1 := INF
       if a < len(dp1) {
           b1 = dp1[a]
       }
       b := min(b0, b1)
       res[a] = n - a - b
   }
   for i, v := range res {
       if i > 0 {
           fmt.Fprint(writer, ' ')
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}
