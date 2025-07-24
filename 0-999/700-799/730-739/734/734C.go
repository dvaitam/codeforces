package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int
   var x, s int64
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   fmt.Fscan(in, &x, &s)
   // type1 spells
   a := make([]int64, m+1)
   b := make([]int64, m+1)
   a[0] = x
   b[0] = 0
   for i := 1; i <= m; i++ {
       fmt.Fscan(in, &a[i])
   }
   for i := 1; i <= m; i++ {
       fmt.Fscan(in, &b[i])
   }
   // type2 spells
   c := make([]int64, k+1)
   d := make([]int64, k+1)
   c[0] = 0
   d[0] = 0
   for i := 1; i <= k; i++ {
       fmt.Fscan(in, &c[i])
   }
   for i := 1; i <= k; i++ {
       fmt.Fscan(in, &d[i])
   }

   var ans int64 = n * x
   // iterate type1 spells
   for i := 0; i <= m; i++ {
       cost1 := b[i]
       if cost1 > s {
           continue
       }
       rem := s - cost1
       // find max j such that d[j] <= rem
       idx := sort.Search(len(d), func(j int) bool { return d[j] > rem })
       j := idx - 1
       if j < 0 {
           j = 0
       }
       // potions left to make
       need := int64(n) - c[j]
       if need < 0 {
           need = 0
       }
       // time
       t := need * a[i]
       if t < ans {
           ans = t
       }
   }
   // output
   fmt.Println(ans)
}
