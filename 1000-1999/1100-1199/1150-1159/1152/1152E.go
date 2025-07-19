package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Edge represents an undirected edge by its endpoint and identifier
type Edge struct {
   to int
   id int
}

func find(parent []int, x int) int {
   if parent[x] != x {
       parent[x] = find(parent, parent[x])
   }
   return parent[x]
}

func union(parent []int, x, y int) {
   fx := find(parent, x)
   fy := find(parent, y)
   if fx != fy {
       parent[fx] = fy
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   b := make([]int, n-1)
   c := make([]int, n-1)
   for i := 0; i < n-1; i++ {
       fmt.Fscan(reader, &b[i])
   }
   for i := 0; i < n-1; i++ {
       fmt.Fscan(reader, &c[i])
   }
   // Validate pairs
   for i := 0; i < n-1; i++ {
       if b[i] > c[i] {
           fmt.Fprintln(writer, -1)
           return
       }
   }
   // Coordinate compress
   S := make([]int, 0, 2*(n-1))
   S = append(S, b...)
   S = append(S, c...)
   sort.Ints(S)
   // unique
   s2 := make([]int, 0, len(S))
   for i, v := range S {
       if i == 0 || v != S[i-1] {
           s2 = append(s2, v)
       }
   }
   tot := len(s2)
   // remap b, c to indices
   for i := 0; i < n-1; i++ {
       b[i] = sort.SearchInts(s2, b[i])
       c[i] = sort.SearchInts(s2, c[i])
   }
   // Union-Find for connectivity
   parent := make([]int, tot)
   for i := 0; i < tot; i++ {
       parent[i] = i
   }
   deg := make([]int, tot)
   for i := 0; i < n-1; i++ {
       deg[b[i]]++
       deg[c[i]]++
       union(parent, b[i], c[i])
   }
   // Check connectivity
   root0 := find(parent, 0)
   for i := 0; i < tot; i++ {
       if find(parent, i) != root0 {
           fmt.Fprintln(writer, -1)
           return
       }
   }
   // Find odd degree vertices
   odds := make([]int, 0, 2)
   for i := 0; i < tot; i++ {
       if deg[i]&1 == 1 {
           odds = append(odds, i)
       }
   }
   if len(odds) != 0 && len(odds) != 2 {
       fmt.Fprintln(writer, -1)
       return
   }
   // Build adjacency
   adj := make([][]Edge, tot)
   for i := 0; i < n-1; i++ {
       u, v := b[i], c[i]
       adj[u] = append(adj[u], Edge{to: v, id: i})
       adj[v] = append(adj[v], Edge{to: u, id: i})
   }
   // Hierholzer's algorithm
   start := 0
   if len(odds) == 2 {
       start = odds[0]
   }
   vis := make([]bool, n-1)
   idx := make([]int, tot)
   stack := make([]int, 0, n)
   stack = append(stack, start)
   pathRev := make([]int, 0, n)
   for len(stack) > 0 {
       v := stack[len(stack)-1]
       // find next unused edge
       for idx[v] < len(adj[v]) && vis[adj[v][idx[v]].id] {
           idx[v]++
       }
       if idx[v] == len(adj[v]) {
           // backtrack
           stack = stack[:len(stack)-1]
           pathRev = append(pathRev, v)
       } else {
           e := adj[v][idx[v]]
           vis[e.id] = true
           stack = append(stack, e.to)
           idx[v]++
       }
   }
   if len(pathRev) != n {
       fmt.Fprintln(writer, -1)
       return
   }
   // Output sequence (reversed Euler path)
   for i := 0; i < n; i++ {
       writer.WriteString(fmt.Sprintf("%d", s2[pathRev[i]]))
       if i+1 < n {
           writer.WriteByte(' ')
       }
   }
   writer.WriteByte('\n')
}
