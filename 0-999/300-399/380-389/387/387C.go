package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var p string
   // Read input p (may be very long)
   if _, err := fmt.Fscan(reader, &p); err != nil {
       return
   }
   n := len(p)
   if n == 0 {
       fmt.Println(0)
       return
   }
   // Count potential splits: each non-zero starting digit yields a new leaf
   // Leaves = 1 (first) + count of p[i] != '0' for i>=1
   leaves := 1
   for i := 1; i < n; i++ {
       if p[i] != '0' {
           leaves++
       }
   }
   // If the first possible split is at position 1 (p[1] != '0')
   // and p[0] < p[1], then making both single-digit leaves violates bi >= bj
   if n >= 2 && p[1] != '0' && p[0] < p[1] {
       leaves--
   }
   fmt.Println(leaves)
}
