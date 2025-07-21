package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   colors := make([]int, n+1)
   const maxC = 100000
   present := make([]bool, maxC+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &colors[i])
       present[colors[i]] = true
   }
   adj := make([][]int, maxC+1)
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       ca := colors[a]
       cb := colors[b]
       if ca != cb {
           adj[ca] = append(adj[ca], cb)
           adj[cb] = append(adj[cb], ca)
       }
   }
   bestCount := -1
   bestColor := 0
   for k := 1; k <= maxC; k++ {
       if !present[k] {
           continue
       }
       neis := adj[k]
       if len(neis) == 0 {
           if bestCount < 0 {
               bestCount = 0
               bestColor = k
           }
           continue
       }
       sort.Ints(neis)
       uniq := 1
       for i := 1; i < len(neis); i++ {
           if neis[i] != neis[i-1] {
               uniq++
           }
       }
       if uniq > bestCount {
           bestCount = uniq
           bestColor = k
       }
   }
   fmt.Println(bestColor)
}
