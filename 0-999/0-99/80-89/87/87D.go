package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// DSU for union-find with sizes
type DSU struct {
   parent []int
   size   []int
}

func NewDSU(n int) *DSU {
   p := make([]int, n)
   sz := make([]int, n)
   for i := 0; i < n; i++ {
       p[i] = i
       sz[i] = 1
   }
   return &DSU{parent: p, size: sz}
}

func (d *DSU) Find(x int) int {
   if d.parent[x] != x {
       d.parent[x] = d.Find(d.parent[x])
   }
   return d.parent[x]
}

func (d *DSU) Union(x, y int) bool {
   rx := d.Find(x)
   ry := d.Find(y)
   if rx == ry {
       return false
   }
   // union by size
   if d.size[rx] < d.size[ry] {
       rx, ry = ry, rx
   }
   d.parent[ry] = rx
   d.size[rx] += d.size[ry]
   return true
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   type Edge struct{ u, v int; w int; idx int }
   edges := make([]Edge, n-1)
   for i := 0; i < n-1; i++ {
       var u, v int
       var w int
       fmt.Fscan(reader, &u, &v, &w)
       edges[i] = Edge{u - 1, v - 1, w, i}
   }
   sort.Slice(edges, func(i, j int) bool {
       return edges[i].w < edges[j].w
   })
   ans := make([]int64, n-1)
   dsu := NewDSU(n)
   // process groups by weight
   for i := 0; i < len(edges); {
       j := i
       w := edges[i].w
       for j < len(edges) && edges[j].w == w {
           j++
       }
       // collect DSU roots
       k := j - i
       roots := make([]int, 0, 2*k)
       for t := i; t < j; t++ {
           ru := dsu.Find(edges[t].u)
           rv := dsu.Find(edges[t].v)
           roots = append(roots, ru, rv)
       }
       sort.Ints(roots)
       uniq := roots[:0]
       for _, x := range roots {
           if len(uniq) == 0 || uniq[len(uniq)-1] != x {
               uniq = append(uniq, x)
           }
       }
       m := len(uniq)
       // map root to id
       idmap := make(map[int]int, m)
       for id, r := range uniq {
           idmap[r] = id
       }
       compSize := make([]int64, m)
       for id, r := range uniq {
           compSize[id] = int64(dsu.size[r])
       }
       // build adjacency for H
       type Info struct{ to, eid int }
       adj := make([][]Info, m)
       for t := i; t < j; t++ {
           ru := dsu.Find(edges[t].u)
           rv := dsu.Find(edges[t].v)
           u := idmap[ru]
           v := idmap[rv]
           adj[u] = append(adj[u], Info{v, edges[t].idx})
           adj[v] = append(adj[v], Info{u, edges[t].idx})
       }
       // visited for components
       visited := make([]bool, m)
       subSum := make([]int64, m)
       // traverse each comp in H
       for id := 0; id < m; id++ {
           if visited[id] {
               continue
           }
           // collect nodes in this H-tree
           stack := []int{id}
           visited[id] = true
           compNodes := []int{id}
           for len(stack) > 0 {
               u := stack[len(stack)-1]
               stack = stack[:len(stack)-1]
               for _, ei := range adj[u] {
                   if !visited[ei.to] {
                       visited[ei.to] = true
                       stack = append(stack, ei.to)
                       compNodes = append(compNodes, ei.to)
                   }
               }
           }
           // total sum of this H-tree
           var totalSum int64
           for _, u := range compNodes {
               totalSum += compSize[u]
           }
           // DFS post-order to compute subtree sums
           type Frame struct{ u, parent, next int; eid int }
           stk := []Frame{{u: id, parent: -1, next: 0, eid: -1}}
           for len(stk) > 0 {
               fr := &stk[len(stk)-1]
               if fr.next < len(adj[fr.u]) {
                   ei := adj[fr.u][fr.next]
                   fr.next++
                   if ei.to == fr.parent {
                       continue
                   }
                   stk = append(stk, Frame{u: ei.to, parent: fr.u, next: 0, eid: ei.eid})
               } else {
                   // process u
                   u := fr.u
                   sum := compSize[u]
                   for _, ei := range adj[u] {
                       if ei.to == fr.parent {
                           continue
                       }
                       sum += subSum[ei.to]
                   }
                   subSum[u] = sum
                   if fr.eid != -1 {
                       // edge fr.eid connects parent and u
                       s := sum
                       ans[fr.eid] = 2 * s * (totalSum - s)
                   }
                   stk = stk[:len(stk)-1]
               }
           }
       }
       // union DSU
       for t := i; t < j; t++ {
           dsu.Union(edges[t].u, edges[t].v)
       }
       i = j
   }
   // find max
   var maxCnt int64
   for _, v := range ans {
       if v > maxCnt {
           maxCnt = v
       }
   }
   var ids []int
   for i, v := range ans {
       if v == maxCnt {
           ids = append(ids, i+1)
       }
   }
   sort.Ints(ids)
   fmt.Fprintf(writer, "%d %d\n", maxCnt, len(ids))
   for i, v := range ids {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprintf(writer, "%d", v)
   }
   writer.WriteByte('\n')
}
