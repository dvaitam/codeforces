package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   // Compute binomial coefficient C(n, k)
   if k > n {
       fmt.Println(0)
       return
   }
   // Use symmetry property C(n, k) == C(n, n-k)
   if k > n-k {
       k = n - k
   }
   var res int64 = 1
   for i := 1; i <= k; i++ {
       res = res * int64(n-k+i) / int64(i)
   }
   fmt.Println(res)
}
