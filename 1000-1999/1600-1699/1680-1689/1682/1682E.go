package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// pair holds a neighbor node and the edge index
type pair struct { first, second int }

var idx []int
var sz, u int

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var N, M int
   fmt.Fscan(reader, &N, &M)
   arr := make([]int, N+1)
   for i := 1; i <= N; i++ {
       fmt.Fscan(reader, &arr[i])
   }
   f := make([][]pair, N+1)
   graph := make([][]int, M+1)
   for i := 1; i <= M; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
      // record edges: neighbor node and edge index
      f[x] = append(f[x], pair{first: y, second: i})
      f[y] = append(f[y], pair{first: x, second: i})
   }
   chk := make([]bool, N+1)
   idx = make([]int, N+1)
   for i := 1; i <= N; i++ {
       if chk[i] {
           continue
       }
       j := i
       var cycle []int
       for !chk[j] {
           chk[j] = true
           idx[j] = len(cycle)
           cycle = append(cycle, j)
           j = arr[j]
       }
       sz = len(cycle)
       if sz <= 1 {
           continue
       }
       for _, u0 := range cycle {
           u = u0
           // sort edges around u by cycle order
           sort.Slice(f[u], func(i1, i2 int) bool {
               a := f[u][i1]
               b := f[u][i2]
               da := (idx[a.first] - idx[u] + sz) % sz
               db := (idx[b.first] - idx[u] + sz) % sz
               return da < db
           })
           // add directed edges between edge indices
           for v := 0; v+1 < len(f[u]); v++ {
               from := f[u][v+1].second
               to := f[u][v].second
               graph[from] = append(graph[from], to)
           }
       }
   }
   visited := make([]bool, M+1)
   var dfs func(int)
   dfs = func(ind int) {
       if visited[ind] {
           return
       }
       visited[ind] = true
       for _, e := range graph[ind] {
           dfs(e)
       }
       fmt.Fprintln(writer, ind)
   }
   for i := 1; i <= M; i++ {
       if !visited[i] {
           dfs(i)
       }
   }
}
