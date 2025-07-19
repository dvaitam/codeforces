package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   adj := make([][]int, n)
   for i := 0; i < m; i++ {
       var v, u int
       fmt.Fscan(reader, &v, &u)
       v--
       u--
       adj[v] = append(adj[v], u)
       adj[u] = append(adj[u], v)
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   pt := make([]int, n)
   queue := make([]int, 0, n)
   res := make([]int, 0, n)
   // initial nodes
   for i := 0; i < n; i++ {
       if pt[i] == a[i] {
           // schedule i
           pt[i]++
           queue = append(queue, i)
           res = append(res, i)
       }
   }
   // process
   for qi := 0; qi < len(queue); qi++ {
       v := queue[qi]
       for _, u := range adj[v] {
           pt[u]++
           if pt[u] == a[u] {
               pt[u]++
               queue = append(queue, u)
               res = append(res, u)
           }
       }
   }
   // output
   writer.WriteString(strconv.Itoa(len(res)))
   writer.WriteByte('\n')
   for i, v := range res {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(strconv.Itoa(v + 1))
   }
   writer.WriteByte('\n')
}
