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
   // adjacency lists per color: colors are 1..m
   adjacency := make([][][]int, m+1)
   for c := 1; c <= m; c++ {
       adjacency[c] = make([][]int, n+1)
   }
   // read edges
   for i := 0; i < m; i++ {
       var a, b, c int
       fmt.Fscan(reader, &a, &b, &c)
       adjacency[c][a] = append(adjacency[c][a], b)
       adjacency[c][b] = append(adjacency[c][b], a)
   }
   // components per color
   comp := make([][]int, m+1)
   for c := 1; c <= m; c++ {
       comp[c] = make([]int, n+1)
       cid := 0
       // find connected components in color c
       for u := 1; u <= n; u++ {
           if comp[c][u] == 0 && len(adjacency[c][u]) > 0 {
               cid++
               // BFS/DFS from u
               stack := []int{u}
               comp[c][u] = cid
               for len(stack) > 0 {
                   v := stack[len(stack)-1]
                   stack = stack[:len(stack)-1]
                   for _, w := range adjacency[c][v] {
                       if comp[c][w] == 0 {
                           comp[c][w] = cid
                           stack = append(stack, w)
                       }
                   }
               }
           }
       }
   }
   // process queries
   var q int
   fmt.Fscan(reader, &q)
   for i := 0; i < q; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       cnt := 0
       for c := 1; c <= m; c++ {
           if comp[c][u] != 0 && comp[c][u] == comp[c][v] {
               cnt++
           }
       }
       fmt.Println(cnt)
   }
}
