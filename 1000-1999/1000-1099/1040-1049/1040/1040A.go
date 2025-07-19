package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var d0, d1 int
   if _, err := fmt.Fscan(reader, &n, &d0, &d1); err != nil {
       return
   }
   D := []int{d0, d1}
   A := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &A[i])
   }
   var res int64
   // process pairs
   for i := 0; i < n/2; i++ {
       j := n - 1 - i
       ai, aj := A[i], A[j]
       if ai != 2 && aj != 2 {
           if ai != aj {
               fmt.Println(-1)
               return
           }
           continue
       }
       if ai == 2 && aj == 2 {
           // choose cheapest for both
           res += int64(2 * min(d0, d1))
           continue
       }
       if ai == 2 {
           // replace ai to match aj
           res += int64(D[aj])
           continue
       }
       // aj == 2
       res += int64(D[ai])
   }
   // middle element if odd
   if n%2 == 1 && A[n/2] == 2 {
       res += int64(min(d0, d1))
   }
   fmt.Println(res)
}
