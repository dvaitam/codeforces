package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var s string
   fmt.Fscan(in, &s)
   var x, y int64
   fmt.Fscan(in, &x, &y)
   n := len(s)
   // fixed bits counts
   var cnt0, cnt1 int64
   fixedCost := int64(0)
   // record for each '?' zerosBefore and onesBefore
   var zb, ob []int64
   for i := 0; i < n; i++ {
       switch s[i] {
       case '0':
           fixedCost += cnt1 * y
           cnt0++
       case '1':
           fixedCost += cnt0 * x
           cnt1++
       case '?':
           zb = append(zb, cnt0)
           ob = append(ob, cnt1)
       }
   }
   total0, total1 := cnt0, cnt1
   m := len(zb)
   // compute cost for each '?' if 0 or 1
   c0 := make([]int64, m)
   c1 := make([]int64, m)
   for i := 0; i < m; i++ {
       onesBefore := ob[i]
       zerosBefore := zb[i]
       onesAfter := total1 - onesBefore
       zerosAfter := total0 - zerosBefore
       // assign 0: with ones before gives (1,0) cost y, with ones after gives (0,1) cost x
       c0[i] = onesBefore*y + onesAfter*x
       // assign 1: with zeros before gives (0,1) cost x, with zeros after gives (1,0) cost y
       c1[i] = zerosBefore*x + zerosAfter*y
   }
   // prefix sums of c0 and c1
   pref0 := make([]int64, m+1)
   pref1 := make([]int64, m+1)
   for i := 0; i < m; i++ {
       pref0[i+1] = pref0[i] + c0[i]
       pref1[i+1] = pref1[i] + c1[i]
   }
   totalC1 := pref1[m]
   // choose K zeros among '?'s (first K assigned 0, rest 1)
   best := int64(1) << 62
   for k := 0; k <= m; k++ {
       // sum cost for first k zeros: pref0[k], and rest ones: totalC1 - pref1[k]
       cost := pref0[k] + (totalC1 - pref1[k])
       // between '?'s: pairs of 0 then 1 cost x each
       cost += x * int64(k) * int64(m-k)
       if cost < best {
           best = cost
       }
   }
   ans := fixedCost + best
   fmt.Println(ans)
}
