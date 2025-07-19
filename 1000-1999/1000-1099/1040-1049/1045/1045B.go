package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   A := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &A[i])
   }
   // special case n == 1
   if n == 1 {
       fmt.Fprintln(out, 1)
       fmt.Fprintln(out, (A[0]+A[0])%m)
       return
   }
   // build difference pattern C of length n-1
   C := make([]int, n-1)
   for i := 0; i < n-1; i++ {
       C[i] = A[i+1] - A[i]
   }
   // build text B of length 2*n-2
   B := make([]int, 2*n-2)
   // reversed differences
   for i := 0; i < n-1; i++ {
       B[i] = C[n-2-i]
   }
   // wrap-around difference
   B[n-1] = (A[0] - A[n-1] + m) % m
   // repeat reversed diffs except last
   for i := 0; i < n-2; i++ {
       B[n+i] = B[i]
   }
   // build KMP failure function for C
   fail := make([]int, n-1)
   for i := 1; i < n-1; i++ {
       j := fail[i-1]
       for j > 0 && C[j] != C[i] {
           j = fail[j-1]
       }
       if C[j] == C[i] {
           j++
       }
       fail[i] = j
   }
   // search
   var ansVals []int
   now := 0
   for i, v := range B {
       for now > 0 && C[now] != v {
           now = fail[now-1]
       }
       if C[now] == v {
           now++
       }
       if now == n-1 {
           // match ends at i (0-based), compute rotation index
           // t = 2*n - i - 2 (1-based), index = t-1 = 2*n - i - 3
           idx := 2*n - i - 3
           if idx >= 0 && idx < n {
               ansVals = append(ansVals, (A[0]+A[idx])%m)
           }
           // continue search
           now = fail[now-1]
       }
   }
   sort.Ints(ansVals)
   // output
   fmt.Fprintln(out, len(ansVals))
   if len(ansVals) > 0 {
       for i, v := range ansVals {
           if i > 0 {
               out.WriteByte(' ')
           }
           fmt.Fprint(out, v)
       }
       out.WriteByte('\n')
   }
}
