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
   const inf = int(1e9)
   k := inf
   for t := 0; t < m; t++ {
       var i, j int
       fmt.Fscan(reader, &i, &j)
       if j-i < k {
           k = j - i
       }
   }
   k++
   fmt.Fprintln(writer, k)
   // output labels for vertices: from n-1 down to 0
   for i := n - 1; i >= 0; i-- {
       fmt.Fprintf(writer, "%d ", i%k)
   }
}
