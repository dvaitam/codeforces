package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   // Read the transformed string t
   line, err := reader.ReadString('\n')
   if err != nil {
       return
   }
   t := strings.TrimSpace(line)
   n := len(t)
   // Prepare three queries encoding positions in base-26
   qs := make([][]byte, 3)
   for k := 0; k < 3; k++ {
       qs[k] = make([]byte, n)
   }
   for i := 0; i < n; i++ {
       x := i
       d0 := x % 26
       x /= 26
       d1 := x % 26
       x /= 26
       d2 := x % 26
       qs[0][i] = byte('a' + d0)
       qs[1][i] = byte('a' + d1)
       qs[2][i] = byte('a' + d2)
   }
   // Responses
   res := make([][]byte, 3)
   for k := 0; k < 3; k++ {
       // Query
       fmt.Fprintf(writer, "? %s\n", string(qs[k]))
       writer.Flush()
       // Read response
       resp, err := reader.ReadString('\n')
       if err != nil {
           return
       }
       res[k] = []byte(strings.TrimSpace(resp))
   }
   // Reconstruct mapping: for each final position j, find original index
   // P^{-1}(j) = orig index
   orig := make([]int, n)
   for j := 0; j < n; j++ {
       d0 := int(res[0][j] - 'a')
       d1 := int(res[1][j] - 'a')
       d2 := int(res[2][j] - 'a')
       idx := d0 + d1*26 + d2*26*26
       orig[j] = idx
   }
   // Build s: s[orig[j]] = t[j]
   s := make([]byte, n)
   for j := 0; j < n; j++ {
       if orig[j] >= 0 && orig[j] < n {
           s[orig[j]] = t[j]
       }
   }
   // Output answer
   fmt.Fprintf(writer, "! %s\n", string(s))
   writer.Flush()
}
