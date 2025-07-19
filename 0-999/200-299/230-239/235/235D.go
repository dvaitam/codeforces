package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n        int
   adj      [][]int
   visited  []bool
   pd       []bool
   dArr     []int
   uArr     []int
   cycle    []int
   a        []float64
   ans      float64
)

// findCycle finds the unique cycle in the graph and returns its nodes in order
func findCycle(u, parent int, path *[]int) []int {
   visited[u] = true
   *path = append(*path, u)
   for _, v := range adj[u] {
       if v == parent {
           continue
       }
       if visited[v] {
           // cycle detected: extract from first occurrence of v
           idx := 0
           for i, x := range *path {
               if x == v {
                   idx = i
                   break
               }
           }
           cycle := make([]int, len(*path)-idx)
           copy(cycle, (*path)[idx:])
           return cycle
       }
       if res := findCycle(v, u, path); len(res) > 0 {
           return res
       }
   }
   *path = (*path)[:len(*path)-1]
   return nil
}

// dfs2 labels each node in a tree attached to cycle node compIdx, computing depth
func dfs2(u, parent, compIdx, depth int) {
   dArr[u] = compIdx
   uArr[u] = depth
   for _, v := range adj[u] {
       if v == parent || pd[v] {
           continue
       }
       dfs2(v, u, compIdx, depth+1)
   }
}

// dfs3 accumulates harmonic sums for distances from start node
func dfs3(u, parent, depth int) {
   ans += a[depth]
   for _, v := range adj[u] {
       if v != parent && !pd[v] {
           dfs3(v, u, depth+1)
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   fmt.Fscan(reader, &n)
   adj = make([][]int, n)
   for i := 0; i < n; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       // convert to 0-based
       u--
       v--
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   visited = make([]bool, n)
   pd = make([]bool, n)
   dArr = make([]int, n)
   uArr = make([]int, n)
   // find cycle
   path := make([]int, 0, n)
   cycle = findCycle(0, -1, &path)
   lenC := len(cycle)
   for i, node := range cycle {
       pd[node] = true
       dArr[node] = i
       uArr[node] = 0
   }
   // prepare harmonic numbers up to enough size
   maxA := 2*n + 3
   a = make([]float64, maxA)
   for i := 1; i < maxA; i++ {
       a[i] = 1.0 / float64(i)
   }
   // process tree components attached to cycle
   for i, node := range cycle {
       dfs2(node, -1, i, 0)
   }
   // accumulate distances within same component
   ans = 0.0
   for i := 0; i < n; i++ {
       root := cycle[dArr[i]]
       pd[root] = false
       dfs3(i, -1, 1)
       pd[root] = true
   }
   // accumulate distances between different components via cycle paths
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if dArr[i] != dArr[j] {
               k1 := uArr[i] + uArr[j]
               diff := abs(dArr[i]-dArr[j]) + 1
               k2 := diff
               k3 := lenC - diff + 2
               ans += a[k1+k2] + a[k1+k3] - a[k1+lenC]
           }
       }
   }
   fmt.Printf("%.11f\n", ans)
}

// abs returns absolute value of x
func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}
