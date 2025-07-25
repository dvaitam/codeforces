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
   fmt.Fscan(reader, &n, &m)
   // record pink edges
   pink := make(map[int]map[int]bool)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       if pink[u] == nil {
           pink[u] = make(map[int]bool)
       }
       pink[u][v] = true
   }
   candidate := 1
   for i := 2; i <= n; i++ {
       // if pink edge between candidate and i
       if pink[candidate] != nil && pink[candidate][i] {
           // candidate->i in pink, keep
           continue
       }
       if pink[i] != nil && pink[i][candidate] {
           // i->candidate in pink, switch
           candidate = i
           continue
       }
       // green edge, query direction
       fmt.Fprintf(writer, "? %d %d\n", candidate, i)
       writer.Flush()
       var ans int
       fmt.Fscan(reader, &ans)
       // ans == 1 means candidate->i, keep; else switch
       if ans == 0 {
           candidate = i
       }
   }
   // found candidate
   fmt.Fprintf(writer, "! %d\n", candidate)
}
