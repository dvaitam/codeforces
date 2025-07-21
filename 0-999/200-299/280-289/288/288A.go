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
   // Impossible cases
   if k > n || (k == 1 && n > 1) {
       fmt.Fprintln(writer, -1)
       return
   }
   // Special case: only one position
   if n == 1 {
       // k must be 1
       fmt.Fprintln(writer, "a")
       return
   }
   // For k >= 2
   // Length of alternating prefix using 'a' and 'b'
   altLen := n - (k - 2)
   res := make([]byte, n)
   for i := 0; i < altLen; i++ {
       if i%2 == 0 {
           res[i] = 'a'
       } else {
           res[i] = 'b'
       }
   }
   // Append the remaining distinct letters from 'c' onwards
   for j := 0; j < k-2; j++ {
       res[altLen+j] = byte('c' + j)
   }
   writer.Write(res)
}
