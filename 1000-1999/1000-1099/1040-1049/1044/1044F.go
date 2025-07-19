package main

import (
   "bufio"
   "fmt"
   "os"
)

const (
   SConst = 1 << 18
   MaxLog = 20
)

var (
   N, Q   int
   adj     [][]int
   stArr   []int
   enArr   []int
   depth   []int
   par     [][]int
   V       int
   S       = SConst
   segVal  []int
   segCnt  []int
)

func buildSeg() {
   segVal = make([]int, 2*S)
   segCnt = make([]int, 2*S)
   // init leaves
   for i := 0; i < S; i++ {
       if i < N {
           segCnt[S+i] = 1
       }
   }
   // build internal
   for i := S - 1; i >= 1; i-- {
       segCnt[i] = segCnt[2*i] + segCnt[2*i+1]
   }
}

func updateNode(i int) {
   if segVal[i] > 0 {
       segCnt[i] = 0
   } else if i >= S {
       idx := i - S
       if idx < N {
           segCnt[i] = 1
       } else {
           segCnt[i] = 0
       }
   } else {
       segCnt[i] = segCnt[2*i] + segCnt[2*i+1]
   }
}

// update range [a,b)
func update(a, b, v, i, st, en int) {
   if a == b {
       return
   }
   if en-st == b-a {
       segVal[i] += v
       updateNode(i)
       return
   }
   md := (st + en) >> 1
   if a < md {
       nb := b
       if nb > md {
           nb = md
       }
       update(a, nb, v, 2*i, st, md)
   }
   if b > md {
       na := a
       if na < md {
           na = md
       }
       update(na, b, v, 2*i+1, md, en)
   }
   updateNode(i)
}

func query() int {
   return segCnt[1]
}

func dfs(cur, prv int) {
   stArr[cur] = V
   V++
   depth[cur] = depth[prv] + 1
   par[cur][0] = prv
   // binary lifting
   for j := 0; par[cur][j] != 0; j++ {
       par[cur][j+1] = par[par[cur][j]][j]
   }
   for _, nxt := range adj[cur] {
       if nxt == prv {
           continue
       }
       dfs(nxt, cur)
   }
   enArr[cur] = V
}

func getAnc(a, d int) int {
   for i := 0; d > 0; i++ {
       if d&1 == 1 {
           a = par[a][i]
       }
       d >>= 1
   }
   return a
}

func doUpdate(a, b, v int) {
   if a == b {
       return
   }
   if depth[a] > depth[b] {
       a, b = b, a
   }
   d := depth[b] - depth[a]
   if getAnc(b, d) == a {
       // a is ancestor of b
       banc := getAnc(b, d-1)
       update(stArr[banc], stArr[b], v, 1, 0, S)
       update(enArr[b], enArr[banc], v, 1, 0, S)
   } else {
       if stArr[a] > stArr[b] {
           a, b = b, a
       }
       update(0, stArr[a], v, 1, 0, S)
       update(enArr[a], stArr[b], v, 1, 0, S)
       update(enArr[b], N, v, 1, 0, S)
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &N, &Q)
   adj = make([][]int, N+1)
   for i := 0; i < N-1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // build segment tree
   buildSeg()
   // prepare dfs structures
   stArr = make([]int, N+1)
   enArr = make([]int, N+1)
   depth = make([]int, N+1)
   par = make([][]int, N+1)
   for i := range par {
       par[i] = make([]int, MaxLog)
   }
   V = 0
   dfs(1, 0)
   // dynamic edge set
   edges := make(map[int64]bool)
   for i := 0; i < Q; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       if u > v {
           u, v = v, u
       }
       key := (int64(u) << 32) | int64(v)
       val := 1
       if edges[key] {
           delete(edges, key)
           val = -1
       } else {
           edges[key] = true
       }
       doUpdate(u, v, val)
       fmt.Fprintln(writer, query())
   }
}
