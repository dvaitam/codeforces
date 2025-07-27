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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
       var ans int64
       for i, v := range a {
           // price = v, buyers = n - i
           cnt := int64(n - i)
           revenue := v * cnt
           if revenue > ans {
               ans = revenue
           }
       }
       fmt.Fprintln(writer, ans)
   }
}
