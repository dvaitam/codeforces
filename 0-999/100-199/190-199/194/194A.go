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
   // x2 is the number of exams with mark 2 (fails), minimize x2
   for x2 := 0; x2 <= n; x2++ {
       // minimum sum with x2 fails: 2*x2 + 3*(n-x2) = 3*n - x2
       // maximum sum with x2 fails: 2*x2 + 5*(n-x2) = 5*n - 3*x2
       minSum := 3*n - x2
       maxSum := 5*n - 3*x2
       if k >= minSum && k <= maxSum {
           fmt.Println(x2)
           return
       }
   }
}
