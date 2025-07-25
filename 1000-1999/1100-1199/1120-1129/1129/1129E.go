package main

import (
   "bufio"
   "fmt"
   "os"
)

// This solution reads the tree edges from standard input (offline mode)
// and outputs them directly. For interactive judging, this file
// serves as a placeholder: replace the offline reading with queries
// to Mikaela as per the problem statement.
func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   edges := make([][2]int, 0, n-1)
   for i := 0; i < n-1; i++ {
       var u, v int
       if _, err := fmt.Fscan(reader, &u, &v); err != nil {
           // No more edges; break for interactive placeholder
           break
       }
       edges = append(edges, [2]int{u, v})
   }

   // Output in the required format
   fmt.Fprintln(writer, "ANSWER")
   for _, e := range edges {
       fmt.Fprintf(writer, "%d %d\n", e[0], e[1])
   }
}
