package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func main() {
   defer writer.Flush()
   // read number of vertices
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // adjacency marker
   graph := make([][]bool, n)
   for i := 0; i < n; i++ {
       graph[i] = make([]bool, n)
   }
   // for each vertex j, query neighbors by setting d_j=1, others 0
   for j := 0; j < n; j++ {
       // output query
       fmt.Fprint(writer, "?")
       for i := 0; i < n; i++ {
           if i == j {
               fmt.Fprint(writer, " 1")
           } else {
               fmt.Fprint(writer, " 0")
           }
       }
       fmt.Fprint(writer, "\n")
       writer.Flush()
       // read response string of 0/1
       resp, err := reader.ReadString('\n')
       if err != nil {
           return
       }
       resp = strings.TrimSpace(resp)
       // record edges
       for i := 0; i < n && i < len(resp); i++ {
           if i != j && resp[i] == '1' {
               graph[j][i] = true
           }
       }
   }
   // output reconstructed tree
   fmt.Fprintln(writer, "!")
   for i := 0; i < n; i++ {
       for k := i + 1; k < n; k++ {
           if graph[i][k] || graph[k][i] {
               fmt.Fprintf(writer, "%d %d\n", i+1, k+1)
           }
       }
   }
}
