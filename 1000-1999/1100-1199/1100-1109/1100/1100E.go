package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Edge struct {
   u, v, c, id int
}

var (
   n, m   int
   edges  []Edge
   adj    [][]int
   inDeg  []int
   tid    []int
)

// check whether suffix edges[mid..m-1] form a DAG
// if fill is true, record topological order in tid
func check(mid int, fill bool) bool {
   for i := 1; i <= n; i++ {
       adj[i] = adj[i][:0]
       inDeg[i] = 0
   }
   for i := mid; i < m; i++ {
       e := edges[i]
       adj[e.u] = append(adj[e.u], e.v)
       inDeg[e.v]++
   }
   queue := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       if inDeg[i] == 0 {
           queue = append(queue, i)
       }
   }
   cnt := 0
   for head := 0; head < len(queue); head++ {
       u := queue[head]
       cnt++
       if fill {
           tid[u] = cnt
       }
       for _, v := range adj[u] {
           inDeg[v]--
           if inDeg[v] == 0 {
               queue = append(queue, v)
           }
       }
   }
   return cnt == n
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n, &m)
   edges = make([]Edge, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &edges[i].u, &edges[i].v, &edges[i].c)
       edges[i].id = i + 1
   }
   sort.Slice(edges, func(i, j int) bool {
       return edges[i].c < edges[j].c
   })
   adj = make([][]int, n+1)
   for i := range adj {
       adj[i] = make([]int, 0)
   }
   inDeg = make([]int, n+1)
   tid = make([]int, n+1)
   lo, hi := 0, m
   for lo < hi {
       mid := (lo + hi) / 2
       if check(mid, false) {
           hi = mid
       } else {
           lo = mid + 1
       }
   }
   ans := lo
   check(ans, true)
   idx := ans
   if idx >= m {
       idx = m - 1
   }
   threshold := edges[idx].c
   var rem []int
   for i := 0; i < ans; i++ {
       e := edges[i]
       if tid[e.u] > tid[e.v] {
           rem = append(rem, e.id)
       }
   }
   sort.Ints(rem)
   fmt.Fprintf(writer, "%d %d\n", threshold, len(rem))
   for _, id := range rem {
       fmt.Fprintf(writer, "%d ", id)
   }
   fmt.Fprint(writer, "\n")
}
