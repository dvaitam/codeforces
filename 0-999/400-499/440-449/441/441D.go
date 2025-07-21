package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(in, &n)
   p := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &p[i])
   }
   fmt.Fscan(in, &m)
   // compute initial cycles
   comp := make([]int, n+1)
   compList := make(map[int][]int)
   visited := make([]bool, n+1)
   compID := 0
   for i := 1; i <= n; i++ {
       if !visited[i] {
           compID++
           u := i
           for !visited[u] {
               visited[u] = true
               comp[u] = compID
               compList[compID] = append(compList[compID], u)
               u = p[u]
           }
       }
   }
   for id := 1; id <= compID; id++ {
       sort.Ints(compList[id])
   }
   cycles := compID
   targetCycles := n - m
   var ops [][2]int

   // rebuild components for given nodes
   rebuild := func(nodes []int) {
       for _, u := range nodes {
           visited[u] = false
           comp[u] = 0
       }
       for _, u := range nodes {
           if !visited[u] {
               compID++
               var group []int
               v := u
               for !visited[v] {
                   visited[v] = true
                   comp[v] = compID
                   group = append(group, v)
                   v = p[v]
               }
               sort.Ints(group)
               compList[compID] = group
           }
       }
   }

   // increase cycles by splitting
   for cycles < targetCycles {
       var ai, bi int
       for i := 1; i <= n; i++ {
           cl := compList[comp[i]]
           if len(cl) >= 2 {
               ai = i
               for _, x := range cl {
                   if x > i {
                       bi = x
                       break
                   }
               }
               break
           }
       }
       ops = append(ops, [2]int{ai, bi})
       p[ai], p[bi] = p[bi], p[ai]
       old := compList[comp[ai]]
       delete(compList, comp[ai])
       rebuild(old)
       cycles++
   }

   // decrease cycles by merging
   for cycles > targetCycles {
       var ai, bi int
       for i := 1; i <= n; i++ {
           for j := i + 1; j <= n; j++ {
               if comp[i] != comp[j] {
                   ai, bi = i, j
                   break
               }
           }
           if ai != 0 {
               break
           }
       }
       ops = append(ops, [2]int{ai, bi})
       p[ai], p[bi] = p[bi], p[ai]
       id1, id2 := comp[ai], comp[bi]
       nodes := append(compList[id1], compList[id2]...)
       delete(compList, id1)
       delete(compList, id2)
       rebuild(nodes)
       cycles--
   }

   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, len(ops))
   if len(ops) > 0 {
       for i, pr := range ops {
           if i > 0 {
               w.WriteByte(' ')
           }
           fmt.Fprintf(w, "%d %d", pr[0], pr[1])
       }
       w.WriteByte('\n')
   }
}
