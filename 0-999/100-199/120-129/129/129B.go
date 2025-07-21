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
   adj := make([][]int, n)
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       a--
       b--
       adj[a] = append(adj[a], b)
       adj[b] = append(adj[b], a)
   }
   deg := make([]int, n)
   removed := make([]bool, n)
   for i := 0; i < n; i++ {
       deg[i] = len(adj[i])
   }
   rounds := 0
   for {
       var toRemove []int
       for i := 0; i < n; i++ {
           if !removed[i] && deg[i] == 1 {
               toRemove = append(toRemove, i)
           }
       }
       if len(toRemove) == 0 {
           break
       }
       rounds++
       // mark removed
       for _, u := range toRemove {
           removed[u] = true
       }
       // update degrees
       for _, u := range toRemove {
           for _, v := range adj[u] {
               if !removed[v] {
                   deg[v]--
               }
           }
       }
   }
   fmt.Fprintln(writer, rounds)
}
