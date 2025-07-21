package main

import (
   "bufio"
   "fmt"
   "os"
)

// Solve: orient edges to maximize reachable pairs count
func main() {
   rdr := bufio.NewReader(os.Stdin)
   wtr := bufio.NewWriter(os.Stdout)
   defer wtr.Flush()

   // fast integer reader
   readInt := func() int {
       var x int
       var c byte
       var err error
       // skip non-digit
       for {
           c, err = rdr.ReadByte()
           if err != nil {
               return 0
           }
           if c >= '0' && c <= '9' {
               break
           }
       }
       // read number
       for c >= '0' && c <= '9' {
           x = x*10 + int(c-'0')
           c, err = rdr.ReadByte()
           if err != nil {
               break
           }
       }
       return x
   }

   n := readInt()
   m := readInt()
   // original graph
   type Edge struct{ to, id int }
   adj := make([][]Edge, n)
   U := make([]int, m)
   V := make([]int, m)
   for i := 0; i < m; i++ {
       u := readInt() - 1
       v := readInt() - 1
       U[i], V[i] = u, v
       adj[u] = append(adj[u], Edge{v, i})
       adj[v] = append(adj[v], Edge{u, i})
   }
   // find bridges
   disc := make([]int, n)
   low := make([]int, n)
   isBridge := make([]bool, m)
   timer := 1
   var dfsBridge func(u, peid int)
   dfsBridge = func(u, peid int) {
       disc[u] = timer
       low[u] = timer
       timer++
       for _, e := range adj[u] {
           v, eid := e.to, e.id
           if eid == peid {
               continue
           }
           if disc[v] == 0 {
               dfsBridge(v, eid)
               if low[v] > disc[u] {
                   isBridge[eid] = true
               }
               if low[v] < low[u] {
                   low[u] = low[v]
               }
           } else if disc[v] < low[u] {
               low[u] = disc[v]
           }
       }
   }
   dfsBridge(0, -1)
   // build 2-edge-connected components
   compId := make([]int, n)
   for i := range compId {
       compId[i] = -1
   }
   compCnt := 0
   for i := 0; i < n; i++ {
       if compId[i] != -1 {
           continue
       }
       // DFS stack
       stack := []int{i}
       compId[i] = compCnt
       for len(stack) > 0 {
           u := stack[len(stack)-1]
           stack = stack[:len(stack)-1]
           for _, e := range adj[u] {
               if compId[e.to] == -1 && !isBridge[e.id] {
                   compId[e.to] = compCnt
                   stack = append(stack, e.to)
               }
           }
       }
       compCnt++
   }
   // component sizes
   compSize := make([]int64, compCnt)
   for i := 0; i < n; i++ {
       compSize[compId[i]]++
   }
   // build component tree
   cadj := make([][]int, compCnt)
   for i := 0; i < m; i++ {
       if isBridge[i] {
           u, v := compId[U[i]], compId[V[i]]
           cadj[u] = append(cadj[u], v)
           cadj[v] = append(cadj[v], u)
       }
   }
   // subtree sums and initial DP
   S := make([]int64, compCnt)
   par := make([]int, compCnt)
   var f0 int64
   var dfs1 func(u, p int)
   dfs1 = func(u, p int) {
       par[u] = p
       S[u] = compSize[u]
       for _, v := range cadj[u] {
           if v == p {
               continue
           }
           dfs1(v, u)
           S[u] += S[v]
       }
       f0 += compSize[u] * S[u]
   }
   dfs1(0, -1)
   // reroot DP
   total := int64(n)
   F := make([]int64, compCnt)
   F[0] = f0
   ans := F[0]
   var dfs2 func(u int)
   dfs2 = func(u int) {
       for _, v := range cadj[u] {
           if v == par[u] {
               continue
           }
           F[v] = F[u] + compSize[v]*(total-S[v]) - compSize[u]*S[v]
           if F[v] > ans {
               ans = F[v]
           }
           dfs2(v)
       }
   }
   dfs2(0)
   // output result
   fmt.Fprintln(wtr, ans)
}
