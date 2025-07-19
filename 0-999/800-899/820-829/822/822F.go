package main

import (
   "bufio"
   "fmt"
   "os"
)

// Edge represents a tree edge with destination node and its index.
type Edge struct {
   to int
   id int
}

var (
   N      int
   edges  [][]Edge
   writer *bufio.Writer
)

// dfs traverses the tree, printing instructions based on length.
func dfs(x, parent int, length float64) {
   deg := len(edges[x])
   step := 2.0 / float64(deg)
   cur := length
   for _, e := range edges[x] {
       if e.to == parent {
           continue
       }
       // update current position along cycle
       cur += step
       if cur > 1.0 {
           cur -= 2.0
       }
       id := e.id + 1
       if cur >= 0.0 {
           // traverse edge x -> e.to
           fmt.Fprintf(writer, "1 %d %d %d %.10f\n", id, x, e.to, cur)
           dfs(e.to, x, cur-1.0)
       } else {
           // traverse edge e.to -> x
           fmt.Fprintf(writer, "1 %d %d %d %.10f\n", id, e.to, x, 1.0+cur)
           dfs(e.to, x, cur+1.0)
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &N)
   edges = make([][]Edge, N+1)
   for i := 0; i < N-1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       edges[u] = append(edges[u], Edge{to: v, id: i})
       edges[v] = append(edges[v], Edge{to: u, id: i})
   }
   // print number of operations
   fmt.Fprintln(writer, N-1)
   dfs(1, -1, 0.0)
}
