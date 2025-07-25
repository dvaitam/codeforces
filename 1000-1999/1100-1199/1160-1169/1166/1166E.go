package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var m, n int
   if _, err := fmt.Fscan(reader, &m, &n); err != nil {
       return
   }
   days := make([][]int, m)
   for i := 0; i < m; i++ {
       var s int
       fmt.Fscan(reader, &s)
       days[i] = make([]int, s)
       for j := 0; j < s; j++ {
           fmt.Fscan(reader, &days[i][j])
           days[i][j]--
       }
       sort.Ints(days[i])
   }
   // Build directed graph: edge i->j if days[i] and days[j] disjoint
   adj := make([][]int, m)
   for i := 0; i < m; i++ {
       for j := 0; j < m; j++ {
           if i == j {
               continue
           }
           // check disjoint
           if disjoint(days[i], days[j]) {
               adj[i] = append(adj[i], j)
           }
       }
   }
   // detect cycle in digraph
   vis := make([]int, m) // 0=unseen,1=visiting,2=done
   var dfs func(int) bool
   dfs = func(u int) bool {
       vis[u] = 1
       for _, v := range adj[u] {
           if vis[v] == 1 {
               return true // cycle
           }
           if vis[v] == 0 {
               if dfs(v) {
                   return true
               }
           }
       }
       vis[u] = 2
       return false
   }
   for i := 0; i < m; i++ {
       if vis[i] == 0 {
           if dfs(i) {
               fmt.Println("impossible")
               return
           }
       }
   }
   fmt.Println("possible")
}

// disjoint returns true if a and b have no common element (both sorted)
func disjoint(a, b []int) bool {
   i, j := 0, 0
   na, nb := len(a), len(b)
   for i < na && j < nb {
       if a[i] == b[j] {
           return false
       }
       if a[i] < b[j] {
           i++
       } else {
           j++
       }
   }
   return true
}
