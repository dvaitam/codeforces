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
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   children := make([][]int, n+1)
   parent := make([]int, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       children[u] = append(children[u], v)
       parent[v] = u
   }
   gift := make([]int, n+1)
   listVal := make([]bool, n+1)
   q := make([]int, n)
   head, tail := 0, 0
   // read gifts and enqueue roots
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &gift[i])
       if gift[i] >= 1 && gift[i] <= n {
           listVal[gift[i]] = true
       }
       if parent[i] == 0 {
           // BFS from root i
           q[tail] = i
           tail++
           for head < tail {
               u := q[head]
               head++
               for _, v := range children[u] {
                   q[tail] = v
                   tail++
               }
           }
       }
   }
   var ans []int
   // process in reverse order
   for i := tail - 1; i >= 0; i-- {
       u := q[i]
       if gift[u] != u {
           p := parent[u]
           if p == 0 || gift[p] != gift[u] {
               fmt.Fprintln(writer, -1)
               return
           }
       }
       if listVal[u] {
           ans = append(ans, u)
       }
   }
   fmt.Fprintln(writer, len(ans))
   for _, u := range ans {
       fmt.Fprintln(writer, u)
   }
}
