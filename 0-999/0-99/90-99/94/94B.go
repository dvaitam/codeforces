package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var m int
   if _, err := fmt.Fscan(reader, &m); err != nil {
       return
   }
   // adjacency matrix for 5 people (1-based indexing)
   var g [6][6]bool
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       g[a][b] = true
       g[b][a] = true
   }
   // check all triples
   for i := 1; i <= 5; i++ {
       for j := i + 1; j <= 5; j++ {
           for k := j + 1; k <= 5; k++ {
               cnt := 0
               if g[i][j] { cnt++ }
               if g[i][k] { cnt++ }
               if g[j][k] { cnt++ }
               if cnt == 3 || cnt == 0 {
                   fmt.Println("WIN")
                   return
               }
           }
       }
   }
   fmt.Println("FAIL")
}
