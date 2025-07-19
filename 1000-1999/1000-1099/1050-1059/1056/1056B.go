package main

import (
   "fmt"
   "os"
)

func main() {
   var n, m int64
   if _, err := fmt.Fscan(os.Stdin, &n, &m); err != nil {
       return
   }
   // Determine cycle length
   var sz int64
   if n < m {
       sz = n
   } else {
       sz = m
   }
   // Compute square residues for one cycle
   C := make([]int, 0, sz)
   for i := int64(1); i <= sz; i++ {
       r := (i * i) % m
       C = append(C, int(r))
   }
   // Count occurrences in full cycles and remainder
   cant := n / sz
   rem := n % sz
   mInt := int(m)
   rems := make([]int64, mInt)
   for _, r := range C {
       rems[r] += cant
   }
   for i := int64(0); i < rem; i++ {
       rems[C[i]]++
   }
   // Sum pairs where (i + j) % m == 0
   var ans int64
   for i := 0; i < mInt; i++ {
       opp := (mInt - i) % mInt
       ans += rems[i] * rems[opp]
   }
   fmt.Fprintln(os.Stdout, ans)
}
