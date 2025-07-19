package main

import (
   "bufio"
   "fmt"
   "os"
)

// Pair represents an operation between two nodes
type Pair struct { u, v int }

var (
   N, K   int
   adj     [][]int
   sz, mx  []int
   marked  []bool
   ct       int
   elem     [2][]int
   ans      []Pair
)

func max(a, b int) int {
   if a > b { return a }
   return b
}

// dfs computes subtree sizes and maximum child subtree size
func dfs(n, p int) int {
   sz[n] = 1
   mx[n] = 0
   for _, v := range adj[n] {
       if v == p || marked[v] {
           continue
       }
       csz := dfs(v, n)
       sz[n] += csz
       mx[n] = max(mx[n], csz)
   }
   return sz[n]
}

// dfs2 finds the centroid of the current subtree
func dfs2(n, p, tot int) {
   if 2*mx[n] <= tot && 2*(tot-sz[n]) <= tot {
       ct = n
   }
   for _, v := range adj[n] {
       if v == p || marked[v] {
           continue
       }
       dfs2(v, n, tot)
   }
}

// op adds operations connecting all nodes in a subtree to root r
func op(n, p, r int) {
   if r != p {
       ans = append(ans, Pair{n, r})
   }
   for _, v := range adj[n] {
       if v == p || marked[v] {
           continue
       }
       op(v, n, r)
   }
}

// op2 collects nodes by parity depth in elem
func op2(n, p, d int) {
   elem[d] = append(elem[d], n)
   for _, v := range adj[n] {
       if v == p || marked[v] {
           continue
       }
       op2(v, n, d^1)
   }
}

// solve performs centroid decomposition and records operations
func solve(n int) {
   tot := dfs(n, 0)
   ct = 0
   dfs2(n, 0, tot)
   cur := ct
   marked[cur] = true
   dfs(cur, 0)
   // find the largest child subtree
   bigSz, bigV := -1, -1
   for _, v := range adj[cur] {
       if marked[v] {
           continue
       }
       if sz[v] > bigSz {
           bigSz = sz[v]
           bigV = v
       }
   }
   // reset element lists
   elem[0] = elem[0][:0]
   elem[1] = elem[1][:0]
   for _, v := range adj[cur] {
       if marked[v] {
           continue
       }
       if v == bigV {
           op2(v, cur, 0)
       } else {
           op(v, cur, cur)
       }
   }
   // choose smaller parity group
   idx := 0
   if len(elem[1]) < len(elem[0]) {
       idx = 1
   }
   for _, v := range elem[idx] {
       if v != bigV {
           ans = append(ans, Pair{cur, v})
       }
   }
   // recurse on remaining subtrees
   for _, v := range adj[cur] {
       if marked[v] {
           continue
       }
       solve(v)
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &N, &K)
   adj = make([][]int, N+1)
   sz = make([]int, N+1)
   mx = make([]int, N+1)
   marked = make([]bool, N+1)
   for i := 1; i < N; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       adj[x] = append(adj[x], y)
       adj[y] = append(adj[y], x)
   }
   solve(1)
   // output operations
   fmt.Fprintln(writer, len(ans))
   for _, p := range ans {
       fmt.Fprintln(writer, p.u, p.v)
   }
}
