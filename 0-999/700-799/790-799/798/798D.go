package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Node represents an element with values x, y and original id
type Node struct {
   x, y int
   id   int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   nodes := make([]Node, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &nodes[i].x)
       nodes[i].id = i + 1
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &nodes[i].y)
   }
   // Sort by x descending
   sort.Slice(nodes, func(i, j int) bool {
       return nodes[i].x > nodes[j].x
   })
   // ans[id] marks selected nodes
   ans := make([]bool, n+1)
   // Pairwise selection based on y
   start := n % 2 // 0 if even, 1 if odd
   for i := start; i+1 < n; i += 2 {
       if nodes[i].y > nodes[i+1].y {
           ans[nodes[i].id] = true
       } else {
           ans[nodes[i+1].id] = true
       }
   }
   // Select one more to reach n/2+1
   for i := 0; i < n; i++ {
       if !ans[nodes[i].id] {
           ans[nodes[i].id] = true
           break
       }
   }
   // Output
   lim := n/2 + 1
   fmt.Fprintln(writer, lim)
   first := true
   for i := 1; i <= n; i++ {
       if ans[i] {
           if !first {
               writer.WriteByte(' ')
           }
           first = false
           fmt.Fprint(writer, i)
       }
   }
   fmt.Fprintln(writer)
}
