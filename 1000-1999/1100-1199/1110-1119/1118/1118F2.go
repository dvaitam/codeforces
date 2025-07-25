package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, k int
   fmt.Fscan(reader, &n, &k)
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // Build parent and depth by iterative DFS, record postorder
   parent := make([]int, n+1)
   depth := make([]int, n+1)
   const LOG = 19
   up := make([][]int, LOG)
   for i := 0; i < LOG; i++ {
       up[i] = make([]int, n+1)
   }
   order := make([]int, 0, n)
   // stack for DFS: v, parent, next child idx
   type stt struct{ v, p, idx int }
   stack := []stt{{1, 0, 0}}
   parent[1] = 0
   depth[1] = 0
   for len(stack) > 0 {
       cur := &stack[len(stack)-1]
       v, p, i := cur.v, cur.p, cur.idx
       if i < len(adj[v]) {
           w := adj[v][i]
           cur.idx++
           if w == p {
               continue
           }
           parent[w] = v
           depth[w] = depth[v] + 1
           stack = append(stack, stt{w, v, 0})
       } else {
           order = append(order, v)
           stack = stack[:len(stack)-1]
       }
   }
   // LCA binary lifting
   for v := 1; v <= n; v++ {
       up[0][v] = parent[v]
   }
   for i := 1; i < LOG; i++ {
       for v := 1; v <= n; v++ {
           up[i][v] = up[i-1][ up[i-1][v] ]
       }
   }
   var lca func(int, int) int
   lca = func(u, v int) int {
       if depth[u] < depth[v] {
           u, v = v, u
       }
       d := depth[u] - depth[v]
       for i := 0; i < LOG; i++ {
           if (d>>i)&1 == 1 {
               u = up[i][u]
           }
       }
       if u == v {
           return u
       }
       for i := LOG-1; i >= 0; i-- {
           if up[i][u] != up[i][v] {
               u = up[i][u]
               v = up[i][v]
           }
       }
       return parent[u]
   }
   // Paths for each color to mark safe edges
   nodesByColor := make([][]int, k+1)
   for v := 1; v <= n; v++ {
       if a[v] > 0 {
           nodesByColor[a[v]] = append(nodesByColor[a[v]], v)
       }
   }
   mark := make([]int, n+1)
   for c := 1; c <= k; c++ {
       list := nodesByColor[c]
       if len(list) < 2 {
           continue
       }
       r := list[0]
       for _, u := range list[1:] {
           w := lca(r, u)
           mark[r]++
           mark[u]++
           mark[w] -= 2
       }
   }
   // accumulate marks to identify safe edges
   for _, v := range order {
       if parent[v] != 0 {
           mark[parent[v]] += mark[v]
       }
   }
   // safe edges if mark[v] > 0 for edge parent[v]-v
   for v := 2; v <= n; v++ {
       if mark[v] > 1 {
           fmt.Fprint(writer, 0)
           return
       }
   }
   // build safe adjacency to form components
   safeAdj := make([][]int, n+1)
   for v := 2; v <= n; v++ {
       if mark[v] == 1 {
           u := parent[v]
           safeAdj[u] = append(safeAdj[u], v)
           safeAdj[v] = append(safeAdj[v], u)
       }
   }
   // assign component IDs
   compID := make([]int, n+1)
   compCnt := 0
   for v := 1; v <= n; v++ {
       if compID[v] == 0 {
           compCnt++
           // BFS
           q := []int{v}
           compID[v] = compCnt
           for i := 0; i < len(q); i++ {
               x := q[i]
               for _, y := range safeAdj[x] {
                   if compID[y] == 0 {
                       compID[y] = compCnt
                       q = append(q, y)
                   }
               }
           }
       }
   }
   // build component tree using non-safe edges
   compAdj := make([][]int, compCnt+1)
   for v := 2; v <= n; v++ {
       if mark[v] == 0 {
           u := parent[v]
           c1, c2 := compID[v], compID[u]
           compAdj[c1] = append(compAdj[c1], c2)
           compAdj[c2] = append(compAdj[c2], c1)
       }
   }
   // identify special components
   hasColor := make([]bool, compCnt+1)
   for v := 1; v <= n; v++ {
       if a[v] > 0 {
           hasColor[compID[v]] = true
       }
   }
   totalSpecial := k
   // find a root component that is special
   root := 0
   for v := 1; v <= n; v++ {
       if a[v] > 0 {
           root = compID[v]
           break
       }
   }
   // DFS on component tree to get postorder and parent
   compParent := make([]int, compCnt+1)
   compOrder := make([]int, 0, compCnt)
   stack2 := []stt{{root, 0, 0}}
   for len(stack2) > 0 {
       cur := &stack2[len(stack2)-1]
       v, p, i := cur.v, cur.p, cur.idx
       if i < len(compAdj[v]) {
           w := compAdj[v][i]
           cur.idx++
           if w == p {
               continue
           }
           compParent[w] = v
           stack2 = append(stack2, stt{w, v, 0})
       } else {
           compOrder = append(compOrder, v)
           stack2 = stack2[:len(stack2)-1]
       }
   }
   // count special in subtrees and count dangerous edges
   spCnt := make([]int, compCnt+1)
   for c := 1; c <= compCnt; c++ {
       if hasColor[c] {
           spCnt[c] = 1
       }
   }
   dangerous := 0
   for _, v := range compOrder {
       if v == root {
           continue
       }
       if spCnt[v] > 0 && totalSpecial-spCnt[v] > 0 {
           dangerous++
       }
       spCnt[compParent[v]] += spCnt[v]
   }
   // compute combinations C(dangerous, k-1)
   mod := 998244353
   if dangerous < k-1 {
       fmt.Fprint(writer, 0)
       return
   }
   maxN := compCnt
   fact := make([]int, maxN+1)
   ifact := make([]int, maxN+1)
   fact[0] = 1
   for i := 1; i <= maxN; i++ {
       fact[i] = fact[i-1] * i % mod
   }
   ifact[maxN] = modInv(fact[maxN], mod)
   for i := maxN; i >= 1; i-- {
       ifact[i-1] = ifact[i] * i % mod
   }
   choose := func(n, r int) int {
       if r < 0 || r > n {
           return 0
       }
       return int((int64(fact[n]) * int64(ifact[r]) % int64(mod)) * int64(ifact[n-r]) % int64(mod))
   }
   res := choose(dangerous, k-1)
   fmt.Fprint(writer, res)
}

// modInv computes modular inverse of a mod m
func modInv(a, m int) int {
   return modPow(a, m-2, m)
}

// modPow computes a^e mod m
func modPow(a, e, m int) int {
   res := 1
   x := a % m
   for e > 0 {
       if e&1 == 1 {
           res = int(int64(res) * int64(x) % int64(m))
       }
       x = int(int64(x) * int64(x) % int64(m))
       e >>= 1
   }
   return res
}
