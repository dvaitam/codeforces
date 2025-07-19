package main

import (
   "bufio"
   "fmt"
   "os"
)

type pair struct {
   s int
   p int
}

func min(a, b int) int {
   if a < b {
      return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
      return
   }
   A := make([]pair, 0, n)
   B := make([]pair, 0, n)
   for i := 1; i <= n; i++ {
      var x, y int
      fmt.Fscan(in, &x, &y)
      if x == 1 {
         A = append(A, pair{y, i})
      } else {
         B = append(B, pair{y, i})
      }
   }
   // reverse B to match original order
   for i, j := 0, len(B)-1; i < j; i, j = i+1, j-1 {
      B[i], B[j] = B[j], B[i]
   }
   ai, bi := 0, 0
   for ai < len(A) && bi < len(B) {
      if A[ai].s == 0 {
         fmt.Fprintf(out, "%d %d 0\n", A[ai].p, B[bi].p)
         ai++
         continue
      }
      // match supply A[ai] with demand B[bi]
      for A[ai].s > 0 && bi < len(B) {
         l := min(A[ai].s, B[bi].s)
         A[ai].s -= l
         B[bi].s -= l
         fmt.Fprintf(out, "%d %d %d\n", A[ai].p, B[bi].p, l)
         if A[ai].s > 0 {
            bi++
         }
      }
      ai++
   }
   // remaining B entries, pair with last A id with 0
   if bi < len(B) {
      lastA := 0
      if ai > 0 {
         lastA = A[ai-1].p
      } else if len(A) > 0 {
         lastA = A[len(A)-1].p
      }
      for ; bi < len(B); bi++ {
         fmt.Fprintf(out, "%d %d 0\n", lastA, B[bi].p)
      }
   }
}
