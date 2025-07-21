package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   var a int64
   if _, err := fmt.Fscan(reader, &n, &m, &a); err != nil {
       return
   }
   b := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &b[i])
   }
   p := make([]int64, m)
   for j := 0; j < m; j++ {
       fmt.Fscan(reader, &p[j])
   }
   // sort personal budgets descending
   sort.Slice(b, func(i, j int) bool {
       return b[i] > b[j]
   })
   // sort bike prices ascending
   sort.Slice(p, func(i, j int) bool {
       return p[i] < p[j]
   })
   // prefix sums
   maxk := n
   if m < n {
       maxk = m
   }
   B := make([]int64, maxk+1)
   for i := 1; i <= maxk; i++ {
       B[i] = B[i-1] + b[i-1]
   }
   P := make([]int64, maxk+1)
   for i := 1; i <= maxk; i++ {
       P[i] = P[i-1] + p[i-1]
   }
   best := 0
   var bestSpend int64 = 0
   for r := 1; r <= maxk; r++ {
       totalRent := P[r]
       // personal needed after using shared budget
       need := totalRent - a
       if need < 0 {
           need = 0
       }
       // check if top r personal budgets can cover need
       if B[r] >= need {
           best = r
           bestSpend = need
       } else {
           // further r will only increase totalRent faster than B[r], so break
           // but to be safe, continue checking
           // continue
       }
   }
   fmt.Fprintln(writer, best, bestSpend)
}
