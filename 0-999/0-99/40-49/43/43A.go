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

   var n int
   fmt.Fscan(reader, &n)
   counts := make(map[string]int)
   var team string
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &team)
       counts[team]++
   }
   winner := ""
   max := -1
   for name, cnt := range counts {
       if cnt > max {
           max = cnt
           winner = name
       }
   }
   fmt.Fprintln(writer, winner)
}
