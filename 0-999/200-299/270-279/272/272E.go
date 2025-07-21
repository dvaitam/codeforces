package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   adj := make([][]int, n)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       u--
       v--
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // color 0 or 1
   col := make([]byte, n)
   // same_count[v]: number of neighbors with same color as v
   same := make([]int, n)
   // initialize same counts (all col=0)
   q := make([]int, 0, n)
   for i := 0; i < n; i++ {
       same[i] = len(adj[i])
       if same[i] >= 2 {
           q = append(q, i)
       }
   }
   // process bad nodes
   for qi := 0; qi < len(q); qi++ {
       v := q[qi]
       if same[v] < 2 {
           continue
       }
       // flip color of v
       old := col[v]
       col[v] ^= 1
       // update same_count for neighbors
       k := same[v]
       d := len(adj[v])
       same[v] = d - k
       for _, u := range adj[v] {
           // before flip, edge v-u contributed if col[u]==old
           if col[u] == old {
               same[u]--
           }
           // after flip, edge contributes if col[u]==col[v]
           if col[u] == col[v] {
               same[u]++
           }
           if same[u] == 2 {
               q = append(q, u)
           }
       }
       // after flip, v has same[v]<2 guaranteed
   }
   // check valid
   for i := 0; i < n; i++ {
       if same[i] >= 2 {
           fmt.Println(-1)
           return
       }
   }
   // output
   out := make([]byte, n)
   for i := 0; i < n; i++ {
       out[i] = '0' + col[i]
   }
   writer := bufio.NewWriter(os.Stdout)
   writer.Write(out)
   writer.WriteByte('\n')
   writer.Flush()
}
