package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var k int64
   fmt.Fscan(reader, &n, &k)
   w := make([]int64, n)
   for i := 1; i < n; i++ {
       fmt.Fscan(reader, &w[i])
   }
   g := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &g[i])
   }
   // prefix sums
   W := make([]int64, n+2)
   G := make([]int64, n+2)
   for i := 1; i <= n; i++ {
       W[i] = W[i-1]
       if i >= 2 {
           W[i] += w[i-1]
       }
       G[i] = G[i-1] + g[i]
   }
   // A[j] = W[j+1] - G[j], B[j] = W[j] - G[j]
   A := make([]int64, n)
   B := make([]int64, n)
   for j := 1; j < n; j++ {
       A[j] = W[j+1] - G[j]
       B[j] = W[j] - G[j]
   }

   // check if any segment of length L feasible
   check := func(L int) bool {
       if L <= 1 {
           return true
       }
       d := L - 1
       // deque for A max and B min, store indices
       qa := make([]int, 0, n)
       qb := make([]int, 0, n)
       for j := 1; j < n; j++ {
           // remove old
           start := j - d + 1
           if len(qa) > 0 && qa[0] < start {
               qa = qa[1:]
           }
           if len(qb) > 0 && qb[0] < start {
               qb = qb[1:]
           }
           // push j into qa (max)
           for len(qa) > 0 && A[qa[len(qa)-1]] <= A[j] {
               qa = qa[:len(qa)-1]
           }
           qa = append(qa, j)
           // push j into qb (min)
           for len(qb) > 0 && B[qb[len(qb)-1]] >= B[j] {
               qb = qb[:len(qb)-1]
           }
           qb = append(qb, j)
           if j >= d {
               l := j - d + 1
               r := l + L - 1
               // r = j+1
               // compute deficits
               ma := A[qa[0]]
               mb := B[qb[0]]
               C := W[l] - G[l-1]
               E := W[r] - G[r]
               fdef := ma - C
               if fdef < 0 {
                   fdef = 0
               }
               bdef := E - mb
               if bdef < 0 {
                   bdef = 0
               }
               need := fdef
               if bdef > need {
                   need = bdef
               }
               if need <= k {
                   return true
               }
           }
       }
       return false
   }

   // binary search max L
   low, high := 1, n
   ans := 1
   for low <= high {
       mid := (low + high) / 2
       if check(mid) {
           ans = mid
           low = mid + 1
       } else {
           high = mid - 1
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans)
}
