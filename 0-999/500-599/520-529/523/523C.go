package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s, t string
   // Read s and t
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   m := len(s)
   n := len(t)
   // forward match positions L for s
   L := make([]int, m)
   pj := 0
   for i := 0; i < n && pj < m; i++ {
       if t[i] == s[pj] {
           L[pj] = i
           pj++
       }
   }
   if pj < m {
       // s cannot be matched even in whole t
       fmt.Println(0)
       return
   }
   // backward match positions R for s
   R := make([]int, m)
   pj = m - 1
   for i := n - 1; i >= 0 && pj >= 0; i-- {
       if t[i] == s[pj] {
           R[pj] = i
           pj--
       }
   }
   if pj >= 0 {
       fmt.Println(0)
       return
   }
   // count cuts i such that prefix ends at >= L[m-1] and suffix starts at <= R[0]
   // valid i range: [L[m-1]+1, R[0]] inclusive
   leftEnd := L[m-1]
   rightStart := R[0]
   ans := rightStart - leftEnd
   if ans < 0 {
       ans = 0
   }
   fmt.Println(ans)
}
