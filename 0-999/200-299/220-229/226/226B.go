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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // sort descending
   sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
   // prefix sums P[i]=sum of a[0..i-1]
   P := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       P[i] = P[i-1] + a[i-1]
   }
   // answer for k=1 (chain)
   var ans1 int64
   for i := 2; i <= n; i++ {
       ans1 += a[i-1] * int64(i-1)
   }
   var q int
   fmt.Fscan(reader, &q)
   ks := make([]int, q)
   uniq := make(map[int]struct{}, q)
   for i := 0; i < q; i++ {
       fmt.Fscan(reader, &ks[i])
       uniq[ks[i]] = struct{}{}
   }
   // compute answers for unique k
   ansMap := make(map[int]int64, len(uniq))
   rem := n - 1
   for k := range uniq {
       if k <= 1 {
           ansMap[k] = ans1
       } else {
           var curCount int
           var kPow int64 = 1
           var ans int64
           d := 1
           for curCount < rem {
               kPow *= int64(k)
               // number of nodes at this level
               level := kPow
               var use int64
               if level > int64(rem-curCount) {
                   use = int64(rem - curCount)
               } else {
                   use = level
               }
               // positions p from curCount+1 to curCount+use map to a indices [p+1]
               lIdx := curCount + 1 + 1
               rIdx := curCount + int(use) + 1
               if lIdx < 2 {
                   lIdx = 2
               }
               if rIdx > n {
                   rIdx = n
               }
               sumB := P[rIdx] - P[lIdx-1]
               ans += sumB * int64(d)
               curCount += int(use)
               d++
           }
           ansMap[k] = ans
       }
   }
   // output in order
   for i, k := range ks {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(fmt.Sprintf("%d", ansMap[k]))
   }
   writer.WriteByte('\n')
}
