package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, k int
   fmt.Fscan(reader, &n, &m, &k)
   c := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &c[i])
       c[i]--
   }
   // adjacency for intra-group edges
   adj := make([][]int, n)
   type Edge struct{u, v int}
   cross := make([]Edge, 0, m)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       u--
       v--
       if c[u] == c[v] {
           adj[u] = append(adj[u], v)
           adj[v] = append(adj[v], u)
       } else {
           cross = append(cross, Edge{u, v})
       }
   }
   // local colors for bipartiteness
   local := make([]int, n)
   for i := 0; i < n; i++ {
       local[i] = -1
   }
   groupGood := make([]bool, k)
   for i := 0; i < k; i++ {
       groupGood[i] = true
   }
   // BFS for each component
   queue := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if local[i] != -1 {
           continue
       }
       // start BFS
       local[i] = 0
       queue = queue[:0]
       queue = append(queue, i)
       for qi := 0; qi < len(queue); qi++ {
           u := queue[qi]
           for _, v := range adj[u] {
               if local[v] == -1 {
                   local[v] = local[u] ^ 1
                   queue = append(queue, v)
               } else if local[v] == local[u] {
                   // conflict in group
                   groupGood[c[u]] = false
               }
           }
       }
   }
   // count good groups
   goodCount := 0
   for i := 0; i < k; i++ {
       if groupGood[i] {
           goodCount++
       }
   }
   // total possible pairs
   total := int64(goodCount) * int64(goodCount-1) / 2

   // collect constraints for cross edges between good groups
   type cons struct{ga, gb int; w int}
   cs := make([]cons, 0, len(cross))
   for _, e := range cross {
       u, v := e.u, e.v
       ga, gb := c[u], c[v]
       if ga == gb || !groupGood[ga] || !groupGood[gb] {
           continue
       }
       w := local[u] ^ local[v] ^ 1
       if ga > gb {
           ga, gb = gb, ga
       }
       cs = append(cs, cons{ga, gb, w})
   }
   // sort by group pair
   sort.Slice(cs, func(i, j int) bool {
       if cs[i].ga != cs[j].ga {
           return cs[i].ga < cs[j].ga
       }
       return cs[i].gb < cs[j].gb
   })
   // scan for conflicts
   i := 0
   for i < len(cs) {
       j := i + 1
       ga, gb := cs[i].ga, cs[i].gb
       w0 := cs[i].w
       bad := false
       for j < len(cs) && cs[j].ga == ga && cs[j].gb == gb {
           if cs[j].w != w0 {
               bad = true
               break
           }
           j++
       }
       if bad {
           total--
       }
       // skip remaining with same pair
       for j < len(cs) && cs[j].ga == ga && cs[j].gb == gb {
           j++
       }
       i = j
   }
   fmt.Fprintln(writer, total)
}
