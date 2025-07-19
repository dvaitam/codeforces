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

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   // count occurrences
   a := make([]int, k+1)
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       if x >= 1 && x <= k {
           a[x]++
       }
   }
   // capacities
   b := make([]int, k+1)
   for i := 1; i <= k; i++ {
       fmt.Fscan(reader, &b[i])
   }
   // compute minimum number of groups
   suf := 0
   num := -1
   for i := k; i >= 1; i-- {
       suf += a[i]
       // ceil(suf / b[i])
       t := suf / b[i]
       if suf % b[i] != 0 {
           t++
       }
       if t > num {
           num = t
       }
   }
   if num <= 0 {
       fmt.Fprintln(writer, 0)
       return
   }
   // distribute elements into groups
   ans := make([][]int, num)
   id := -1
   for i := 1; i <= k; i++ {
       for j := 0; j < a[i]; j++ {
           id++
           idx := id % num
           ans[idx] = append(ans[idx], i)
       }
   }
   // output
   fmt.Fprintln(writer, num)
   for _, group := range ans {
       fmt.Fprint(writer, len(group))
       for _, v := range group {
           fmt.Fprint(writer, " ", v)
       }
       fmt.Fprintln(writer)
   }
}
