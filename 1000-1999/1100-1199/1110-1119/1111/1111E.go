package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const (
   MOD = 1000000007
   X   = 18
)

var (
   adj [][]int
   par [][X]int
   dep []int
   st, en []int
   tim int
   grp [][]pair
   radj [][]int
   isThere []int
)

type pair struct{ first, second int }
type vpair struct{ dep, node int }
type qpair struct{ v, d int }

func dfs(v, p int) {
   tim++
   st[v] = tim
   dep[v] = dep[p] + 1
   par[v][0] = p
   for i := 1; i < X; i++ {
       par[v][i] = par[par[v][i-1]][i-1]
   }
   for _, u := range adj[v] {
       if u == p {
           continue
       }
       dfs(u, v)
   }
   en[v] = tim
}

func lca(a, b int) int {
   if dep[b] > dep[a] {
       a, b = b, a
   }
   // lift a
   for i := X-1; i >= 0; i-- {
       if dep[a]-(1<<i) >= dep[b] {
           a = par[a][i]
       }
   }
   if a == b {
       return a
   }
   for i := X-1; i >= 0; i-- {
       if par[a][i] != par[b][i] {
           a = par[a][i]
           b = par[b][i]
       }
   }
   return par[a][0]
}

// build virtual tree for grp of a node
func makeTree(V []pair, l, r, p int) {
   if l > r {
       return
   }
   v := V[l].second
   if p != v {
       radj[p] = append(radj[p], v)
   }
   // split children under v by subtree range
   cnt := r - l
   // find first index in V[l+1:l+1+cnt] with first > en[v]
   k := sort.Search(cnt, func(i int) bool { return V[l+1+i].first > en[v] })
   p2 := l + k
   // recurse on left [l+1, p2] and right [p2+1, r]
   makeTree(V, l+1, p2, v)
   makeTree(V, p2+1, r, p)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, q int
   if _, err := fmt.Fscan(reader, &n, &q); err != nil {
       return
   }
   adj = make([][]int, n+1)
   for i := 1; i < n; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   par = make([][X]int, n+1)
   dep = make([]int, n+1)
   st = make([]int, n+1)
   en = make([]int, n+1)
   tim = 0
   dfs(1, 0)
   grp = make([][]pair, n+1)
   radj = make([][]int, n+1)
   isThere = make([]int, n+1)

   for ; q > 0; q-- {
       var k, m, r int
       fmt.Fscan(reader, &k, &m, &r)
       A := make([]int, k)
       V := make([]vpair, 0, k)
       for i := 0; i < k; i++ {
           fmt.Fscan(reader, &A[i])
           isThere[A[i]] = 1
           l := lca(A[i], r)
           V = append(V, vpair{dep: dep[l], node: l})
           grp[l] = append(grp[l], pair{first: st[A[i]], second: A[i]})
       }
       // unique and sort V by dep
       sort.Slice(V, func(i, j int) bool {
           if V[i].dep != V[j].dep {
               return V[i].dep < V[j].dep
           }
           return V[i].node < V[j].node
       })
       // remove duplicates
       u := 0
       for i := 1; i < len(V); i++ {
           if V[i].node != V[u].node {
               u++
               V[u] = V[i]
           }
       }
       V = V[:u+1]
       // build trees
       for _, vp := range V {
           // sort grp[vp.node] by first
           g := grp[vp.node]
           sort.Slice(g, func(i, j int) bool { return g[i].first < g[j].first })
           grp[vp.node] = g
           makeTree(grp[vp.node], 0, len(g)-1, vp.node)
       }
       // link between components
       for i := 1; i < len(V); i++ {
           radj[V[i].node] = append(radj[V[i].node], V[i-1].node)
       }
       // attach root
       radj[0] = append(radj[0], V[len(V)-1].node)

       dp := make([]int64, m+1)
       tdp := make([]int64, m+1)
       dp[0] = 1
       // BFS
       queue := make([]qpair, 0, len(V)+1)
       queue = append(queue, qpair{v: 0, d: 0})
       for head := 0; head < len(queue); head++ {
           cur := queue[head]
           v := cur.v
           d := cur.d
           for _, to := range radj[v] {
               queue = append(queue, qpair{v: to, d: d + isThere[v]})
           }
           if isThere[v] == 0 {
               continue
           }
           for i := 1; i <= m; i++ {
               var val int64 = dp[i-1]
               if d < i {
                   val = (val + dp[i]*int64(i-d)) % MOD
               }
               tdp[i] = val
           }
           for i := 0; i <= m; i++ {
               dp[i] = tdp[i]
               tdp[i] = 0
           }
       }
       var ans int64
       for i := 0; i <= m; i++ {
           ans = (ans + dp[i]) % MOD
       }
       fmt.Fprintln(writer, ans)
       // cleanup
       for _, x := range A {
           isThere[x] = 0
           radj[x] = radj[x][:0]
       }
       for _, vp := range V {
           grp[vp.node] = grp[vp.node][:0]
           radj[vp.node] = radj[vp.node][:0]
       }
       radj[0] = radj[0][:0]
   }
}
