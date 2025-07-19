package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   neigh := make([][]int, n)
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       x--
       y--
       neigh[x] = append(neigh[x], y)
       neigh[y] = append(neigh[y], x)
   }
   type Node struct { id int; d []int }
   nodes := make([]Node, n)
   for i := 0; i < n; i++ {
       // sort adjacency lists
       sort.Ints(neigh[i])
       // copy slice for node
       nodes[i] = Node{id: i, d: neigh[i]}
   }
   // sort nodes by adjacency list lexicographically
   sort.Slice(nodes, func(i, j int) bool {
       a, b := nodes[i].d, nodes[j].d
       la, lb := len(a), len(b)
       for k := 0; k < la && k < lb; k++ {
           if a[k] != b[k] {
               return a[k] < b[k]
           }
       }
       return la < lb
   })
   // must have non-empty adjacency for first group
   if len(nodes[0].d) == 0 {
       fmt.Fprintln(writer, -1)
       return
   }
   labels := make([]int, n)
   cnt := 1
   labels[nodes[0].id] = cnt
   // assign groups
   for i := 1; i < n; i++ {
       prev, cur := nodes[i-1].d, nodes[i].d
       if len(prev) == len(cur) {
           same := true
           for k := 0; k < len(prev); k++ {
               if prev[k] != cur[k] {
                   same = false
                   break
               }
           }
           if !same {
               cnt++
           }
       } else {
           cnt++
       }
       labels[nodes[i].id] = cnt
   }
   if cnt != 3 {
       fmt.Fprintln(writer, -1)
       return
   }
   // output labels (1-indexed groups)
   for i := 0; i < n; i++ {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(fmt.Sprint(labels[i]))
   }
   writer.WriteByte('\n')
}
