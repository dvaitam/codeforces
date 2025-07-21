package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   f := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &f[i])
   }
   adj := make([][]int, n+1)
   radj := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       adj[a] = append(adj[a], b)
       radj[b] = append(radj[b], a)
   }
   // Forward reachability from assigns without passing other assigns
   fwd := make([]bool, n+1)
   q := make([]int, 0, n)
   // init with assign nodes
   for i := 1; i <= n; i++ {
       if f[i] == 1 {
           fwd[i] = true
           q = append(q, i)
       }
   }
   for qi := 0; qi < len(q); qi++ {
       u := q[qi]
       for _, v := range adj[u] {
           if f[v] == 1 {
               continue
           }
           if !fwd[v] {
               fwd[v] = true
               q = append(q, v)
           }
       }
   }
   // Backward reachability from uses without passing assigns
   bakNonAssign := make([]bool, n+1)
   q = q[:0]
   // seeds: use nodes (f==2), which are non-assign
   for i := 1; i <= n; i++ {
       if f[i] == 2 {
           bakNonAssign[i] = true
           q = append(q, i)
       }
   }
   for qi := 0; qi < len(q); qi++ {
       u := q[qi]
       for _, v := range radj[u] {
           if f[v] == 1 {
               continue
           }
           if !bakNonAssign[v] {
               bakNonAssign[v] = true
               q = append(q, v)
           }
       }
   }
   // Compute final interesting: forward reachable and backward reachable
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 1; i <= n; i++ {
       interesting := false
       if fwd[i] {
           if f[i] == 1 {
               // assign node: check if any neighbor is non-assign reachable to use
               for _, v := range adj[i] {
                   if bakNonAssign[v] {
                       interesting = true
                       break
                   }
               }
           } else {
               // non-assign node
               if bakNonAssign[i] {
                   interesting = true
               }
           }
       }
       if interesting {
           writer.WriteByte('1')
       } else {
           writer.WriteByte('0')
       }
       if i < n {
           writer.WriteByte(' ')
       }
   }
   writer.WriteByte('\n')
}
