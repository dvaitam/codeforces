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

   var t, n, m, a, b int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       fmt.Fscan(reader, &n, &m, &a, &b)
       // Check feasibility
       if n*a != m*b {
           fmt.Fprintln(writer, "NO")
           continue
       }
       // Initialize answer matrix
       ans := make([][]byte, n)
       for i := 0; i < n; i++ {
           ans[i] = make([]byte, m)
       }
       // Fill ones
       j := 0
       for i := 0; i < n; i++ {
           for k := 0; k < a; k++ {
               ans[i][j] = '1'
               j++
               if j == m {
                   j = 0
               }
           }
       }
       // Output
       fmt.Fprintln(writer, "YES")
       for i := 0; i < n; i++ {
           // ensure each column sum equals b (optional)
           writer.Write(ans[i])
           writer.WriteByte('\n')
       }
   }
}
