package main

import (
   "bufio"
   "fmt"
   "os"
)

type edge struct {
   to, w, id int
}

var (
   n int
   adj [][]edge
   edges [][3]int
   // working arrays
   visited []bool
   compNodes []int
   sizeArr []int
   downSum []int64
   fullSum []int64
)

// collect nodes of component starting at u, skipping edge skipID
func collect(u, skipID int) {
   visited[u] = true
   compNodes = append(compNodes, u)
   for _, e := range adj[u] {
       if e.id == skipID || visited[e.to] {
           continue
       }
       collect(e.to, skipID)
   }
}

// dfs1: compute sizeArr[u] and downSum[u] for component of root r
func dfs1(u, parent, skipID int) {
   sizeArr[u] = 1
   downSum[u] = 0
   for _, e := range adj[u] {
       v := e.to
       if e.id == skipID || v == parent || !visited[v] {
           continue
       }
       dfs1(v, u, skipID)
       sizeArr[u] += sizeArr[v]
       downSum[u] += downSum[v] + int64(sizeArr[v]) * int64(e.w)
   }
}

// dfs2: compute fullSum for component, track min
func dfs2(u, parent, compSize int, skipID int, minPtr *int64) {
   // update min
   if fullSum[u] < *minPtr {
       *minPtr = fullSum[u]
   }
   for _, e := range adj[u] {
       v := e.to
       if e.id == skipID || v == parent || !visited[v] {
           continue
       }
       // reroot from u to v
       // distance sums adjust by: fullSum[v] = fullSum[u] + w*(compSize - 2*sizeArr[v])
       fullSum[v] = fullSum[u] + int64(e.w) * int64(compSize - 2*sizeArr[v])
       dfs2(v, u, compSize, skipID, minPtr)
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n)
   adj = make([][]edge, n)
   edges = make([][3]int, n-1)
   for i := 0; i < n-1; i++ {
       var a, b, w int
       fmt.Fscan(in, &a, &b, &w)
       a--
       b--
       edges[i] = [3]int{a, b, w}
       adj[a] = append(adj[a], edge{b, w, i})
       adj[b] = append(adj[b], edge{a, w, i})
   }
   // compute original fullSum
   visited = make([]bool, n)
   sizeArr = make([]int, n)
   downSum = make([]int64, n)
   fullSum = make([]int64, n)
   // full tree: skipID = -1, mark all visited
   for i := range visited {
       visited[i] = true
   }
   // dfs1 from 0
   dfs1(0, -1, -1)
   // reroot dp
   fullSum[0] = downSum[0]
   minTmp := fullSum[0]
   dfs2(0, -1, n, -1, &minTmp)
   // total sum of pairs is sum(fullSum)/2
   var sumAll int64
   for i := 0; i < n; i++ {
       sumAll += fullSum[i]
   }
   sumAll /= 2
   // prepare visited for components
   // will reuse visited slice, clear to false
   for i := range visited {
       visited[i] = false
   }
   sizeArr = make([]int, n)
   downSum = make([]int64, n)
   fullSum = make([]int64, n)
   bestDelta := int64(0)
   // try each edge
   for eid, e := range edges {
       u, v, _ := e[0], e[1], e[2]
       // component B starting at u
       compNodes = compNodes[:0]
       collect(u, eid)
       compSizeB := len(compNodes)
       // compute DP for B
       // dfs1 on comp B
       dfs1(u, -1, eid)
       fullSum[u] = downSum[u]
       minTmpB := fullSum[u]
       dfs2(u, -1, compSizeB, eid, &minTmpB)
       sumBu := fullSum[u]
       // reset visited for B
       for _, x := range compNodes {
           visited[x] = false
       }
       // component A at v
       compNodes = compNodes[:0]
       collect(v, eid)
       compSizeA := len(compNodes)
       dfs1(v, -1, eid)
       fullSum[v] = downSum[v]
       minTmpA := fullSum[v]
       dfs2(v, -1, compSizeA, eid, &minTmpA)
       sumAv := fullSum[v]
       // reset visited for A
       for _, x := range compNodes {
           visited[x] = false
       }
       // compute delta = b*(minA - sumAv) + a*(minB - sumBu)
       a := int64(compSizeA)
       b := int64(compSizeB)
       delta := b*(minTmpA - sumAv) + a*(minTmpB - sumBu)
       if eid == 0 || delta < bestDelta {
           bestDelta = delta
       }
   }
   // output answer
   fmt.Fprintln(out, sumAll+bestDelta)
}
