package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, h int
   if _, err := fmt.Fscan(reader, &n, &m, &h); err != nil {
       return
   }
   s := make([]int, m+1)
   for i := 1; i <= m; i++ {
       fmt.Fscan(reader, &s[i])
   }
   total := 0
   for i := 1; i <= m; i++ {
       total += s[i]
   }
   if total < n {
       fmt.Println(-1)
       return
   }
   // number of other spots to fill
   B := n - 1
   // total non-department players
   A := total - s[h]
   var probNone float64 = 1.0
   if B > 0 {
       if A < B {
           probNone = 0.0
       } else {
           // probability of choosing no teammate from his department
           for i := 0; i < B; i++ {
               probNone *= float64(A - i) / float64(total-1 - i)
           }
       }
   }
   P := 1.0 - probNone
   fmt.Printf("%.10f\n", P)
}
