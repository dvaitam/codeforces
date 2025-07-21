package main

import (
   "bufio"
   "fmt"
   "os"
)

// DSU for components
type DSU struct {
   p []int
}
func NewDSU(n int) *DSU {
   d := &DSU{p: make([]int, n)}
   for i := range d.p {
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
       d.p[ry] = rx
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   // store edges
   u := make([]int, m)
   v := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &u[i], &v[i])
       u[i]--
       v[i]--
   }
   dsu := NewDSU(n)
   for i := 0; i < m; i++ {
       dsu.Union(u[i], v[i])
   }
   // map root to vertices, edge count
   compV := make([][]int, n)
   compE := make([]int, n)
   for i := 0; i < n; i++ {
       compV[dsu.Find(i)] = append(compV[dsu.Find(i)], i)
   }
   for i := 0; i < m; i++ {
       compE[dsu.Find(u[i])]++
   }
   // check cyclomatic
   for r, verts := range compV {
       if len(verts) == 0 {
           continue
       }
       if compE[r]-len(verts)+1 > 1 {
           fmt.Println("NO")
           return
       }
   }
   // detect cycle vertices
   inCycle := make([]bool, n)
   // build adjacency with edge index
   adj := make([][]int, n)
   for i := 0; i < m; i++ {
       // store edge index*2 for u->v, i*2+1 for v->u
       adj[u[i]] = append(adj[u[i]], i*2)
       if u[i] != v[i] {
           adj[v[i]] = append(adj[v[i]], i*2+1)
       } else {
           // loop, mark immediately
           inCycle[u[i]] = true
       }
   }
   used := make([]bool, n)
   var seenEdge = make([]bool, 2*m)
   var stack []int
   var parentEdge = make([]int, n)
   var found bool
   var cycleEdges []int
   var dfs func(int)
   dfs = func(x int) {
       if found {
           return
       }
       used[x] = true
       for _, ei := range adj[x] {
           if seenEdge[ei] {
               continue
           }
           seenEdge[ei] = true
           // other direction index
           other := ei ^ 1
           seenEdge[other] = true
           y := u[ei/2]
           if ei%2 == 1 {
               y = v[ei/2]
           }
           if !used[y] {
               parentEdge[y] = ei
               dfs(y)
               if found {
                   return
               }
           } else {
               // visited and not parent, found cycle
               if parentEdge[x] == ei {
                   continue
               }
               // record cycle edges
               found = true
               // backtrack from x to y
               cycleEdges = append(cycleEdges, ei)
               cur := x
               for cur != y {
                   pe := parentEdge[cur]
                   cycleEdges = append(cycleEdges, pe)
                   // move cur to parent of pe
                   px := u[pe/2]
                   if pe%2 == 1 {
                       px = v[pe/2]
                   }
                   // px is parent of cur? if not, swap
                   if px == cur {
                       // go to other end
                       if u[pe/2] == cur {
                           px = v[pe/2]
                       } else {
                           px = u[pe/2]
                       }
                   }
                   cur = px
               }
               return
           }
       }
   }
   for r, verts := range compV {
       if len(verts) == 0 || compE[r]-len(verts)+1 != 1 {
           continue
       }
       // find cycle in this comp
       // reset
       for _, vtx := range verts {
           used[vtx] = false
       }
       for _, ei := range cycleEdges {
           _ = ei
       }
       cycleEdges = cycleEdges[:0]
       found = false
       // start dfs at first vertex
       dfs(verts[0])
       // mark cycle nodes
       for _, ei := range cycleEdges {
           a := u[ei/2]
           b := v[ei/2]
           inCycle[a] = true
           inCycle[b] = true
       }
   }
   // collect added edges
   var addU, addV []int
   for i := 0; i < n; i++ {
       if !inCycle[i] {
           addU = append(addU, i)
           addV = append(addV, i)
       }
   }
   fmt.Println("YES")
   fmt.Println(len(addU))
   for i := range addU {
       // output 1-based
       fmt.Printf("%d %d\n", addU[i]+1, addV[i]+1)
   }
}
