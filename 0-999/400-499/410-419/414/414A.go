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

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   // Special case: n=1 and k=0
   if k == 0 && n == 1 {
       fmt.Fprintln(writer, 1)
       return
   }
   // Impossible cases
   if k < n/2 || n == 1 {
       fmt.Fprintln(writer, -1)
       return
   }
   // Determine first pair
   base := k - ((n/2) - 1)
   x := base * 2
   res := make([]int, 0, n)
   res = append(res, base, x)
   // Remaining pairs
   for i := 4; i <= n; i += 2 {
       res = append(res, x+1, x+2)
       x += 2
   }
   // If odd, add last element
   if n%2 == 1 {
       res = append(res, x+1)
   }
   // Output result
   for i, v := range res {
       if i > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}
