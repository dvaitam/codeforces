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

   var n, k int
   fmt.Fscan(reader, &n, &k)
   xs := make([]int64, n)
   ys := make([]int64, n)
   for i := 0; i < k; i++ {
       fmt.Fscan(reader, &xs[i], &ys[i])
   }
   var a, b, c, d int64
   fmt.Fscan(reader, &a, &b, &c, &d)
   mod := int64(1e9 + 9)
   for i := k; i < n; i++ {
       xs[i] = (a*xs[i-1] + b) % mod
       ys[i] = (c*ys[i-1] + d) % mod
   }
   // sort by x
   idx := make([]int, n)
   for i := range idx {
       idx[i] = i
   }
   sort.Slice(idx, func(i, j int) bool {
       return xs[idx[i]] < xs[idx[j]]
   })
   sortedX := make([]int64, n)
   for i, id := range idx {
       sortedX[i] = xs[id]
   }
   var m int
   fmt.Fscan(reader, &m)
   for qi := 0; qi < m; qi++ {
       var L, R int64
       fmt.Fscan(reader, &L, &R)
       l := sort.Search(n, func(i int) bool { return sortedX[i] >= L })
       r := sort.Search(n, func(i int) bool { return sortedX[i] > R }) - 1
       cnt := 0
       if l < n && r >= l {
           cnt = r - l + 1
       }
       fmt.Fprintln(writer, cnt/2)
   }
}
