package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// DSU
type DSU struct{ p []int }
func NewDSU(n int) *DSU {
   d := &DSU{p: make([]int, n+1)}
   for i := 1; i <= n; i++ {
       d.p[i] = i
   }
   return d
}
func (d *DSU) Find(x int) int {
   if d.p[x] != x {
       d.p[x] = d.Find(d.p[x])
   }
   return d.p[x]
}
func (d *DSU) Union(x, y int) {
   rx, ry := d.Find(x), d.Find(y)
   if rx != ry {
       d.p[rx] = ry
   }
}

// BFS result
type bfsRes struct{
   node, dist int
   parent []int
}

// BFS from start over adj list of size n+1
func bfs(start, n int, adj [][]int) bfsRes {
   dist := make([]int, n+1)
   parent := make([]int, n+1)
   for i := 1; i <= n; i++ {
       dist[i] = -1
   }
   q := make([]int, 0, n)
   q = append(q, start)
   dist[start] = 0
   parent[start] = -1
   resNode := start
   resDist := 0
   for i := 0; i < len(q); i++ {
       u := q[i]
       if dist[u] > resDist {
           resDist = dist[u]
           resNode = u
       }
       for _, v := range adj[u] {
           if dist[v] < 0 {
               dist[v] = dist[u] + 1
               parent[v] = u
               q = append(q, v)
           }
       }
   }
   return bfsRes{node: resNode, dist: resDist, parent: parent}
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(in, &n, &m)
   dsu := NewDSU(n)
   adj := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
       dsu.Union(u, v)
   }
   // components
   compMap := make(map[int][]int)
   for i := 1; i <= n; i++ {
       r := dsu.Find(i)
       compMap[r] = append(compMap[r], i)
   }
   type comp struct{ center, radius int }
   comps := make([]comp, 0, len(compMap))
   // find center and radius for each component
   for _, nodes := range compMap {
       u0 := nodes[0]
       // first BFS to find one end
       r1 := bfs(u0, n, adj)
       // second BFS from r1.node
       r2 := bfs(r1.node, n, adj)
       // build path from r2.node to r1.node
       path := []int{r2.node}
       for cur := r2.node; cur != r1.node; cur = r2.parent[cur] {
           path = append(path, r2.parent[cur])
       }
       // center at middle of path
       L := len(path)
       center := path[L/2]
       // radius is max dist to ends
       d1 := L/2
       d2 := L - 1 - L/2
       radius := d1
       if d2 > radius {
           radius = d2
       }
       comps = append(comps, comp{center: center, radius: radius})
   }
   // sort by descending radius
   sort.Slice(comps, func(i, j int) bool {
       return comps[i].radius > comps[j].radius
   })
   // add edges
   var newEdges [][2]int
   if len(comps) > 0 {
       mainCenter := comps[0].center
       for i := 1; i < len(comps); i++ {
           c := comps[i].center
           newEdges = append(newEdges, [2]int{mainCenter, c})
           adj[mainCenter] = append(adj[mainCenter], c)
           adj[c] = append(adj[c], mainCenter)
       }
   }
   // compute diameter of full tree
   // BFS from 1 (or any)
   start := 1
   if n > 0 {
       start = 1
   }
   r1 := bfs(start, n, adj)
   r2 := bfs(r1.node, n, adj)
   diameter := r2.dist
   // output
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fprintln(out, diameter)
   for _, e := range newEdges {
       fmt.Fprintln(out, e[0], e[1])
   }
}
