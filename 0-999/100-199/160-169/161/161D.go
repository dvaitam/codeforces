package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, k int
   adj [][]int
   removed []bool
   sz []int
   ans int64
)

func dfsSize(u, p int) {
   sz[u] = 1
   for _, v := range adj[u] {
       if v != p && !removed[v] {
           dfsSize(v, u)
           sz[u] += sz[v]
       }
   }
}

func dfsCentroid(u, p, total int) int {
   for _, v := range adj[u] {
       if v != p && !removed[v] {
           if sz[v] > total/2 {
               return dfsCentroid(v, u, total)
           }
       }
   }
   return u
}

func dfsCount(u, p, depth int, cnt []int) {
   if depth > k {
       return
   }
   cnt[depth]++
   for _, v := range adj[u] {
       if v != p && !removed[v] {
           dfsCount(v, u, depth+1, cnt)
       }
   }
}

func solve(u int) {
   dfsSize(u, -1)
   c := dfsCentroid(u, -1, sz[u])
   removed[c] = true
   // freq[d] = number of nodes at distance d from centroid in processed subtrees
   freq := make([]int, k+1)
   freq[0] = 1
   for _, v := range adj[c] {
       if removed[v] {
           continue
       }
       cnt := make([]int, k+1)
       dfsCount(v, c, 1, cnt)
       // count pairs between this subtree and previous
       for d := 1; d <= k; d++ {
           if cnt[d] > 0 && k-d >= 0 {
               ans += int64(cnt[d]) * int64(freq[k-d])
           }
       }
       // merge cnt into freq
       for d := 1; d <= k; d++ {
           freq[d] += cnt[d]
       }
   }
   // recurse into subtrees
   for _, v := range adj[c] {
       if !removed[v] {
           solve(v)
       }
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   adj = make([][]int, n+1)
   removed = make([]bool, n+1)
   sz = make([]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   ans = 0
   solve(1)
   fmt.Fprintln(out, ans)
}
