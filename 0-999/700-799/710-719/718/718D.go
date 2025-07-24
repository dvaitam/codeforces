package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// DSU for vertices 1..n
type DSU struct {
   p, r []int
}
func NewDSU(n int) *DSU {
   p := make([]int, n+1)
   r := make([]int, n+1)
   for i := range p {
       p[i] = i
   }
   return &DSU{p, r}
}
func (d *DSU) Find(x int) int {
   if d.p[x] != x {
       d.p[x] = d.Find(d.p[x])
   }
   return d.p[x]
}
func (d *DSU) Union(a, b int) {
   a = d.Find(a); b = d.Find(b)
   if a == b {
       return
   }
   if d.r[a] < d.r[b] {
       a, b = b, a
   }
   d.p[b] = a
   if d.r[a] == d.r[b] {
       d.r[a]++
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   adj := make([][]int, n+1)
   deg := make([]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
       deg[u]++
       deg[v]++
   }
   // find centroids: compute subtree sizes and max component sizes
   parent := make([]int, n+1)
   sz := make([]int, n+1)
   order := make([]int, 0, n)
   // post-order DFS to record parent and order
   type Item struct{u, p, idx int}
   stack := []Item{{1, 0, 0}}
   parent[1] = 0
   for len(stack) > 0 {
       it := &stack[len(stack)-1]
       u, p := it.u, it.p
       if it.idx < len(adj[u]) {
           v := adj[u][it.idx]
           it.idx++
           if v == p {
               continue
           }
           parent[v] = u
           stack = append(stack, Item{v, u, 0})
       } else {
           order = append(order, u)
           stack = stack[:len(stack)-1]
       }
   }
   // compute sizes
   for _, u := range order {
       sz[u] = 1
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           sz[u] += sz[v]
       }
   }
   // find centroids based on largest component size
   minSize := n + 1
   centroids := make([]int, 0, 2)
   for u := 1; u <= n; u++ {
       maxPart := n - sz[u]
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           if sz[v] > maxPart {
               maxPart = sz[v]
           }
       }
       if maxPart < minSize {
           minSize = maxPart
           centroids = centroids[:0]
           centroids = append(centroids, u)
       } else if maxPart == minSize {
           centroids = append(centroids, u)
       }
   }
   // if two centroids, remove the edge between them and use dummy root 0
   if len(centroids) == 2 {
       c1, c2 := centroids[0], centroids[1]
       // remove c2 from adj[c1]
       for i, v := range adj[c1] {
           if v == c2 {
               adj[c1] = append(adj[c1][:i], adj[c1][i+1:]...)
               break
           }
       }
       // remove c1 from adj[c2]
       for i, v := range adj[c2] {
           if v == c1 {
               adj[c2] = append(adj[c2][:i], adj[c2][i+1:]...)
               break
           }
       }
   }
   // build rooted tree children
   var root int
   children := make([][]int, n+1)
   if len(centroids) == 1 {
       root = centroids[0]
   } else {
       // two centroids: use 0 as dummy root
       root = 0
       children = make([][]int, n+1) // index 0..n
   }
   // rebuild children with BFS/stack
   st2 := []int{root}
   par2 := make([]int, len(children))
   par2[root] = -1
   for len(st2) > 0 {
       u := st2[len(st2)-1]
       st2 = st2[:len(st2)-1]
       if u == 0 {
           // connect centroids
           c1, c2 := centroids[0], centroids[1]
           children[u] = append(children[u], c1, c2)
           par2[c1] = 0
           par2[c2] = 0
           st2 = append(st2, c1, c2)
           continue
       }
       for _, v := range adj[u] {
           if v == par2[u] {
               continue
           }
           par2[v] = u
           children[u] = append(children[u], v)
           st2 = append(st2, v)
       }
   }
   // post-order
   order = []int{}
   type SI struct{u, idx int}
   st3 := []SI{{root, 0}}
   for len(st3) > 0 {
       top := &st3[len(st3)-1]
       u := top.u
       if top.idx < len(children[u]) {
           v := children[u][top.idx]
           top.idx++
           st3 = append(st3, SI{v, 0})
       } else {
           order = append(order, u)
           st3 = st3[:len(st3)-1]
       }
   }
   // subtree type IDs
   subtype := make([]int, len(children))
   nextID := 1
   key2id := make(map[string]int)
   dsu := NewDSU(n)
   // process nodes in post-order
   for _, u := range order {
       // gather child subtype IDs
       ids := make([]int, 0, len(children[u]))
       for _, v := range children[u] {
           ids = append(ids, subtype[v])
       }
       sort.Ints(ids)
       key := fmt.Sprint(ids)
       id, ok := key2id[key]
       if !ok {
           id = nextID
           key2id[key] = id
           nextID++
       }
       subtype[u] = id
       // group symmetric children and map isomorphic subtrees
       groups := make(map[int][]int)
       for _, v := range children[u] {
           cid := subtype[v]
           groups[cid] = append(groups[cid], v)
       }
       for _, vs := range groups {
           if len(vs) <= 1 {
               continue
           }
           root0 := vs[0]
           for i := 1; i < len(vs); i++ {
               // map subtree root0 and vs[i]
               queue := [][2]int{{root0, vs[i]}}
               for len(queue) > 0 {
                   pair := queue[len(queue)-1]
                   queue = queue[:len(queue)-1]
                   a, b := pair[0], pair[1]
                   dsu.Union(a, b)
                   ca := children[a]
                   cb := children[b]
                   // copy and sort children by subtype
                   la := len(ca)
                   aList := make([]int, la)
                   copy(aList, ca)
                   bList := make([]int, la)
                   copy(bList, cb)
                   sort.Slice(aList, func(i, j int) bool { return subtype[aList[i]] < subtype[aList[j]] })
                   sort.Slice(bList, func(i, j int) bool { return subtype[bList[i]] < subtype[bList[j]] })
                   for k := 0; k < la; k++ {
                       queue = append(queue, [2]int{aList[k], bList[k]})
                   }
               }
           }
       }
   }
   // count orbits with deg<4
   seen := make(map[int]bool)
   res := 0
   for i := 1; i <= n; i++ {
       if deg[i] < 4 {
           leader := dsu.Find(i)
           if !seen[leader] {
               seen[leader] = true
               res++
           }
       }
   }
   fmt.Fprintln(out, res)
}
