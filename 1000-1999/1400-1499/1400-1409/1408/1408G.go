package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const mod = 998244353

// Edge represents an undirected edge with weight
type Edge struct {
   u, v, w int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   // Read edges
   edges := make([]Edge, 0, n*(n-1)/2)
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           var a int
           fmt.Fscan(reader, &a)
           if i < j {
               edges = append(edges, Edge{i, j, a})
           }
       }
   }
   sort.Slice(edges, func(i, j int) bool { return edges[i].w < edges[j].w })

   // DSU initialization
   parent := make([]int, n)
   size := make([]int, n)
   edgesCnt := make([]int, n)
   dp := make([][]int, n)
   for i := 0; i < n; i++ {
       parent[i] = i
       size[i] = 1
       edgesCnt[i] = 0
       dp[i] = make([]int, 2)
       dp[i][1] = 1
   }

   // Process edges in increasing order
   for _, e := range edges {
       ru := find(parent, e.u)
       rv := find(parent, e.v)
       if ru != rv {
           // Union by size: attach smaller ru to larger rv
           if size[ru] > size[rv] {
               ru, rv = rv, ru
           }
           newSize := size[ru] + size[rv]
           // Convolution of dp[ru] and dp[rv]
           newDp := make([]int, newSize+1)
           for i := 1; i <= size[ru]; i++ {
               vi := dp[ru][i]
               if vi == 0 {
                   continue
               }
               for j := 1; j <= size[rv]; j++ {
                   newDp[i+j] = (newDp[i+j] + vi*dp[rv][j]) % mod
               }
           }
           parent[ru] = rv
           size[rv] = newSize
           edgesCnt[rv] += edgesCnt[ru] + 1
           dp[rv] = newDp
           // If component is now a clique, we can form one cluster
           if edgesCnt[rv] == newSize*(newSize-1)/2 {
               dp[rv][1] = (dp[rv][1] + 1) % mod
           }
       } else {
           // Edge inside component
           edgesCnt[ru]++
           if edgesCnt[ru] == size[ru]*(size[ru]-1)/2 {
               dp[ru][1] = (dp[ru][1] + 1) % mod
           }
       }
   }

   // The whole graph is connected, find root
   root := find(parent, 0)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // Output dp[root][k] for k = 1..n
   for k := 1; k <= n; k++ {
       if k < len(dp[root]) {
           fmt.Fprint(writer, dp[root][k])
       } else {
           fmt.Fprint(writer, 0)
       }
       if k < n {
           fmt.Fprint(writer, " ")
       }
   }
}

// find returns the representative of x with path compression
func find(parent []int, x int) int {
   if parent[x] != x {
       parent[x] = find(parent, parent[x])
   }
   return parent[x]
}
