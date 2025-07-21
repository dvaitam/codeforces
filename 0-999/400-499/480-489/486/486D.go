package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const mod = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   var d, n int
   fmt.Fscan(reader, &d, &n)
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
   // sort nodes by value
   type node struct{ v, id int }
   nodes := make([]node, n)
   for i := 1; i <= n; i++ {
       nodes[i-1] = node{a[i], i}
   }
   sort.Slice(nodes, func(i, j int) bool { return nodes[i].v < nodes[j].v })
   vals := make([]int, n)
   ids := make([]int, n)
   for i := 0; i < n; i++ {
       vals[i] = nodes[i].v
       ids[i] = nodes[i].id
   }
   allowed := make([]bool, n+1)
   var ans int64
   // DP function
   var dfs func(u, p int) int64
   dfs = func(u, p int) int64 {
       res := int64(1)
       for _, v := range adj[u] {
           if v == p || !allowed[v] {
               continue
           }
           c := dfs(v, u)
           res = res * (c + 1) % mod
       }
       return res
   }
   // enumerate minimal value index
   for L := 0; L < n; L++ {
       minVal := vals[L]
       limit := minVal + d
       // find first index > limit
       R := sort.Search(n, func(i int) bool { return vals[i] > limit })
       // mark allowed nodes
       for i := L; i < R; i++ {
           allowed[ids[i]] = true
       }
       // compute connected subgraphs containing ids[L]
       root := ids[L]
       cnt := dfs(root, 0)
       ans = (ans + cnt) % mod
       // reset allowed
       for i := L; i < R; i++ {
           allowed[ids[i]] = false
       }
   }
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans)
}
