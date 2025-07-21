package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var m int
   fmt.Fscan(in, &m)
   minQ := int(1e9)
   for i := 0; i < m; i++ {
       var q int
       fmt.Fscan(in, &q)
       if q < minQ {
           minQ = q
       }
   }
   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   sort.Sort(sort.Reverse(sort.IntSlice(a)))
   k := minQ + 2
   var res int64
   for i, v := range a {
       if i%k < minQ {
           res += int64(v)
       }
   }
   fmt.Println(res)
}
