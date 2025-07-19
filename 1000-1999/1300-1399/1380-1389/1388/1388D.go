package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)

   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   b := make([]int, n)
   for i := 0; i < n; i++ {
       var bi int
       fmt.Fscan(reader, &bi)
       if bi == -1 {
           b[i] = -1
       } else {
           b[i] = bi - 1
       }
   }

   // f values
   f := make([]int64, n)
   copy(f, a)
   // indegree for initial graph
   indeg := make([]int, n)
   for i := 0; i < n; i++ {
       if b[i] >= 0 {
           indeg[b[i]]++
       }
   }
   // build edges
   graph := make([][]int, n)
   // indegree for DAG
   deg := make([]int, n)

   // initial queue
   q := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if indeg[i] == 0 {
           q = append(q, i)
       }
   }
   // process in top order
   for head := 0; head < len(q); head++ {
       x := q[head]
       if b[x] >= 0 {
           p := b[x]
           if f[x] < 0 {
               // process x after p: edge p->x
               graph[p] = append(graph[p], x)
               deg[x]++
           } else {
               // process x before p: edge x->p, and add f[x]
               f[p] += f[x]
               graph[x] = append(graph[x], p)
               deg[p]++
           }
           indeg[p]--
           if indeg[p] == 0 {
               q = append(q, p)
           }
       }
   }
   // compute answer
   var ans int64
   for i := 0; i < n; i++ {
       ans += f[i]
   }
   fmt.Fprintln(writer, ans)

   // final order topological sort
   q2 := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if deg[i] == 0 {
           q2 = append(q2, i)
       }
   }
   order := make([]int, 0, n)
   for head := 0; head < len(q2); head++ {
       x := q2[head]
       order = append(order, x)
       for _, y := range graph[x] {
           deg[y]--
           if deg[y] == 0 {
               q2 = append(q2, y)
           }
       }
   }
   // print order 1-based
   for i, v := range order {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v+1)
   }
   writer.WriteByte('\n')
}
