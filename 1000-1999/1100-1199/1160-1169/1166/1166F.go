package main

import (
   "bufio"
   "fmt"
   "os"
)

// DSU structure for union-find
type DSU struct {
   parent []int
   size   []int
}

// NewDSU creates a new DSU with n elements (1-based)
func NewDSU(n int) *DSU {
   parent := make([]int, n+1)
   size := make([]int, n+1)
   for i := 1; i <= n; i++ {
       parent[i] = i
       size[i] = 1
   }
   return &DSU{parent: parent, size: size}
}

// Find with path compression
func (d *DSU) Find(x int) int {
   if d.parent[x] != x {
       d.parent[x] = d.Find(d.parent[x])
   }
   return d.parent[x]
}

// Union by size
func (d *DSU) Union(a, b int) {
   ra := d.Find(a)
   rb := d.Find(b)
   if ra == rb {
       return
   }
   if d.size[ra] < d.size[rb] {
       ra, rb = rb, ra
   }
   d.parent[rb] = ra
   d.size[ra] += d.size[rb]
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, c, q int
   fmt.Fscan(reader, &n, &m, &c, &q)
   dsu := NewDSU(n)
   // adjacency: for each node, map color -> neighbors
   adj := make([]map[int][]int, n+1)
   for i := 1; i <= n; i++ {
       adj[i] = make(map[int][]int)
   }

   // helper to add edge and process meta-unions
   addEdge := func(x, y, z int) {
       // for each neighbor v of x with color z, union(v, y)
       if lst, ok := adj[x][z]; ok {
           for _, v := range lst {
               dsu.Union(v, y)
           }
       }
       // for each neighbor u of y with color z, union(u, x)
       if lst, ok := adj[y][z]; ok {
           for _, u := range lst {
               dsu.Union(u, x)
           }
       }
       // add to adjacency
       adj[x][z] = append(adj[x][z], y)
       adj[y][z] = append(adj[y][z], x)
   }

   // initial roads
   for i := 0; i < m; i++ {
       var x, y, z int
       fmt.Fscan(reader, &x, &y, &z)
       addEdge(x, y, z)
   }

   // process events
   for i := 0; i < q; i++ {
       var op byte
       fmt.Fscan(reader, &op)
       if op == '+' {
           var x, y, z int
           fmt.Fscan(reader, &x, &y, &z)
           addEdge(x, y, z)
       } else if op == '?' {
           var x, y int
           fmt.Fscan(reader, &x, &y)
           if dsu.Find(x) == dsu.Find(y) {
               fmt.Fprintln(writer, "Yes")
           } else {
               fmt.Fprintln(writer, "No")
           }
       }
   }
}
