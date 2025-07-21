package main

import (
   "bufio"
   "fmt"
   "os"
   "runtime"
   "strings"
   "strconv"
   "math/bits"
)

func main() {
   runtime.GOMAXPROCS(runtime.NumCPU())
   reader := bufio.NewReader(os.Stdin)
   line, _ := reader.ReadString('\n')
   parts := strings.Fields(line)
   if len(parts) < 2 {
       return
   }
   n, _ := strconv.Atoi(parts[0])
   m, _ := strconv.Atoi(parts[1])
   // adjacency list and edge list
   adj := make([][]int, n)
   edges := make([][2]int, 0, m)
   for i := 0; i < m; i++ {
       l, _ := reader.ReadString('\n')
       f := strings.Fields(l)
       if len(f) < 2 {
           i--
           continue
       }
       u, _ := strconv.Atoi(f[0]); v, _ := strconv.Atoi(f[1])
       u--; v--
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
       edges = append(edges, [2]int{u, v})
   }
   // bitsets
   w := (n + 63) >> 6
   bs := make([][]uint64, n)
   for i := 0; i < n; i++ {
       b := make([]uint64, w)
       // set self
       b[i>>6] |= 1 << (uint(i) & 63)
       for _, v := range adj[i] {
           b[v>>6] |= 1 << (uint(v) & 63)
       }
       bs[i] = b
   }
   // best coverage
   best := 0
   // single vertex
   for i := 0; i < n; i++ {
       cnt := 0
       for j := 0; j < w; j++ {
           cnt += bits.OnesCount64(bs[i][j])
       }
       if cnt > best {
           best = cnt
       }
   }
   // path of length 1 (edge)
   tmp := make([]uint64, w)
   for _, e := range edges {
       u, v := e[0], e[1]
       cnt := 0
       bu, bv := bs[u], bs[v]
       for j := 0; j < w; j++ {
           tmp[j] = bu[j] | bv[j]
           cnt += bits.OnesCount64(tmp[j])
       }
       if cnt > best {
           best = cnt
       }
   }
   // result: merges = n - best clusters
   res := n - best
   fmt.Println(res)
}
