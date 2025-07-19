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
   var m int64
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([]int64, n)
   var sum int64
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       sum += a[i]
   }
   if sum < m {
       fmt.Fprintln(writer, -1)
       return
   }
   sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
   // try days from 1 to m
   for day := int64(1); day <= m; day++ {
       var tot int64
       var cnt int64
       var t int64
       // greedily count capacity
       for i := n - 1; i >= 0; i-- {
           if a[i] < cnt {
               break
           }
           tot += a[i] - cnt
           t++
           if t == day {
               cnt++
               t = 0
           }
           if tot >= m {
               break
           }
       }
       if tot >= m {
           fmt.Fprintln(writer, day)
           return
       }
   }
   fmt.Fprintln(writer, -1)
}
