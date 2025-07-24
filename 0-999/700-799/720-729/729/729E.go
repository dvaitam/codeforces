package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, s int
   if _, err := fmt.Fscan(reader, &n, &s); err != nil {
       return
   }
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // cost for chief if misreported
   costChief := 0
   if a[s] != 0 {
       costChief = 1
   }
   // count reports among non-chief
   r := make([]int, n+2)
   for i := 1; i <= n; i++ {
       if i == s {
           continue
       }
       ai := a[i]
       if ai < 0 {
           ai = 0
       }
       if ai > n {
           ai = n
       }
       r[ai]++
   }
   // r0: count of non-chief reporting depth 0
   r0 := r[0]
   // total non-chief nodes
   total := n - 1
   // big: count of reporting depth >=1
   big := total - r0
   bestP := total + 5 // large
   z := 0
   // Sweep D from 1..n
   for D := 1; D <= n; D++ {
       // remove r[D] from big
       if D <= n {
           big -= r[D]
       }
       // missing level if no one reported D
       if D < len(r) {
           if r[D] == 0 {
               z++
           }
       } else {
           // beyond reported range, r[D]==0
           z++
       }
       // mistakes P = max(missing levels, nodes not matchable)
       // nodes not matchable = r0 + big
       A := r0 + big
       P := z
       if A > P {
           P = A
       }
       if P < bestP {
           bestP = P
       }
   }
   // total mistakes
   ans := costChief + bestP
   fmt.Fprintln(writer, ans)
}
