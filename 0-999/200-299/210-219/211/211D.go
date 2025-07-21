package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // compute L: distance to previous strictly less
   L := make([]int, n)
   stack := make([]int, 0, n)
   for i := 0; i < n; i++ {
       for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
           stack = stack[:len(stack)-1]
       }
       if len(stack) == 0 {
           L[i] = i + 1
       } else {
           L[i] = i - stack[len(stack)-1]
       }
       stack = append(stack, i)
   }
   // compute R: distance to next less or equal
   R := make([]int, n)
   stack = stack[:0]
   for i := n - 1; i >= 0; i-- {
       for len(stack) > 0 && a[stack[len(stack)-1]] > a[i] {
           stack = stack[:len(stack)-1]
       }
       if len(stack) == 0 {
           R[i] = n - i
       } else {
           R[i] = stack[len(stack)-1] - i
       }
       stack = append(stack, i)
   }
   // difference arrays for alpha (da) and beta (db)
   da := make([]int64, n+3)
   db := make([]int64, n+3)
   for i := 0; i < n; i++ {
       A := L[i]
       B := R[i]
       mn := A
       mx := B
       if B < A {
           mn = B
           mx = A
       }
       total := int64(A + B - 1)
       ai := a[i]
       // [1..mn]: alpha += ai
       if mn > 0 {
           da[1] += ai
           da[mn+1] -= ai
       }
       // [mn+1..mx]: beta += ai * mn
       if mx > mn {
           v := ai * int64(mn)
           da[mn+1] += 0 // no-op to ensure indices in range
           db[mn+1] += v
           db[mx+1] -= v
       }
       // [mx+1..A+B-1]: alpha += -ai; beta += ai*(A+B)
       if int(mx)+1 <= int(total) {
           l := mx + 1
           r := int(total)
           da[l] -= ai
           da[r+1] += ai
           db[l] += ai * int64(A+B)
           db[r+1] -= ai * int64(A+B)
       }
   }
   // build prefix sums and compute S[k]
   S := make([]float64, n+1)
   var curA, curB int64
   for k := 1; k <= n; k++ {
       curA += da[k]
       curB += db[k]
       // S = alpha*k + beta
       S[k] = float64(curA)*float64(k) + float64(curB)
   }
   var m int
   fmt.Fscan(in, &m)
   // answer queries
   for i := 0; i < m; i++ {
       var k int
       fmt.Fscan(in, &k)
       denom := float64(n - k + 1)
       ans := S[k] / denom
       if i > 0 {
           out.WriteByte(' ')
       }
       // print with high precision
       fmt.Fprintf(out, "%.10f", ans)
   }
   out.WriteByte('\n')
}
