package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n1, n2, n3 int
   if _, err := fmt.Fscan(reader, &n1, &n2, &n3); err != nil {
       return
   }
   ns := []int{n1, n2, n3}
   // read three trees
   adjList := make([][][]int, 3)
   for t := 0; t < 3; t++ {
       n := ns[t]
       adj := make([][]int, n)
       for i := 0; i < n-1; i++ {
           var u, v int
           fmt.Fscan(reader, &u, &v)
           u--; v--
           adj[u] = append(adj[u], v)
           adj[v] = append(adj[v], u)
       }
       adjList[t] = adj
   }
   // compute info for each tree
   type info struct {
       n         int
       D0        int64
       sumDepth  []int64
       dist      [][]int
       maxDepth  int64
   }
   infos := make([]info, 3)
   for t := 0; t < 3; t++ {
       n := ns[t]
       adj := adjList[t]
       // D0 via DFS
       var D0 int64
       sz := make([]int, n)
       var dfs1 func(u, p int)
       dfs1 = func(u, p int) {
           sz[u] = 1
           for _, v := range adj[u] {
               if v == p {
                   continue
               }
               dfs1(v, u)
               sz[u] += sz[v]
               D0 += int64(sz[v]) * int64(n - sz[v])
           }
       }
       dfs1(0, -1)
       // compute dist matrix and sumDepth
       dist := make([][]int, n)
       sumDepth := make([]int64, n)
       for i := 0; i < n; i++ {
           // BFS from i
           d := make([]int, n)
           for j := 0; j < n; j++ {
               d[j] = -1
           }
           queue := make([]int, 0, n)
           d[i] = 0
           queue = append(queue, i)
           for qi := 0; qi < len(queue); qi++ {
               u := queue[qi]
               for _, v := range adj[u] {
                   if d[v] < 0 {
                       d[v] = d[u] + 1
                       queue = append(queue, v)
                   }
               }
           }
           dist[i] = d
           var sd int64
           for _, dd := range d {
               sd += int64(dd)
           }
           sumDepth[i] = sd
       }
       // find maxDepth
       var maxD int64
       for _, sd := range sumDepth {
           if sd > maxD {
               maxD = sd
           }
       }
       infos[t] = info{n: n, D0: D0, sumDepth: sumDepth, dist: dist, maxDepth: maxD}
   }
   // total internal distances
   var totalD0 int64
   for t := 0; t < 3; t++ {
       totalD0 += infos[t].D0
   }
   var best int64
   // try each tree as middle
   for j := 0; j < 3; j++ {
       // identify ends a,k
       var a, k int
       for t := 0; t < 3; t++ {
           if t != j {
               if a == 0 && t != 0 || a != 0 {
                   if a == 0 && t != j {
                       a = t
                   } else if a != 0 && k == 0 {
                       k = t
                   }
               } else {
                   a = t
               }
           }
       }
       // fix a,k selection, simpler do list
       idxs := []int{}
       for t := 0; t < 3; t++ {
           if t != j {
               idxs = append(idxs, t)
           }
       }
       a, k = idxs[0], idxs[1]
       ia, ij, ik := infos[a], infos[j], infos[k]
       sa, sj, sk := int64(ia.n), int64(ij.n), int64(ik.n)
       // best inner for middle tree: choose v,w
       var bestInner int64
       // prefetch s_a and s_k
       for v := 0; v < ij.n; v++ {
           for w := 0; w < ij.n; w++ {
               val := sa*ij.sumDepth[v] + sk*ij.sumDepth[w] + sa*sk*int64(ij.dist[v][w])
               if val > bestInner {
                   bestInner = val
               }
           }
       }
       // compute cross sum
       cross := ia.maxDepth*(sj+sk) + ik.maxDepth*(sj+sa) + bestInner + sa*sj + sj*sk + 2*sa*sk
       if cross > best {
           best = cross
       }
   }
   ans := totalD0 + best
   fmt.Println(ans)
}
