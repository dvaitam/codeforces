package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   if n < 3 {
       fmt.Println("NO")
       return
   }
   sort.Ints(a)
   // build unique values and counts
   uniq := make([]int, 0, n)
   cnt := make([]int, 0, n)
   last := a[0]
   uniq = append(uniq, last)
   cnt = append(cnt, 1)
   for i := 1; i < n; i++ {
       if a[i] == last {
           cnt[len(cnt)-1]++
       } else {
           // new value
           if a[i] != last+1 {
               fmt.Println("NO")
               return
           }
           last = a[i]
           uniq = append(uniq, last)
           cnt = append(cnt, 1)
       }
   }
   m := len(cnt)
   // dynamic edge counts e between uniq[i] and uniq[i+1]
   // ePrev corresponds to edge count before uniq[0], must be 0
   var ePrev int
   // At uniq[0], eCur = 2*cnt[0] - ePrev
   for i := 0; i < m; i++ {
       eCur := 2*cnt[i] - ePrev
       if eCur < 0 {
           fmt.Println("NO")
           return
       }
       ePrev = eCur
   }
   // After last, ePrev must be 0 to close cycle
   if ePrev != 0 {
       fmt.Println("NO")
   } else {
       fmt.Println("YES")
   }
}
