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
   var k int64
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })

   used := make(map[int64]bool, n)
   var cnt int
   for _, v := range a {
       if k != 0 && v%k == 0 {
           if used[v/k] {
               continue
           }
       }
       used[v] = true
       cnt++
   }
   fmt.Fprintln(writer, cnt)
}
