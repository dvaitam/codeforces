package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var k, m int64
   fmt.Fscan(reader, &n, &k, &m)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // initial average
   var sum0 float64
   for i := 0; i < n; i++ {
       sum0 += float64(a[i])
   }
   ans := sum0 / float64(n)
   // sort ascending
   sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
   // consider suffixes
   var suffixSum float64
   for j := n - 1; j >= 0; j-- {
       suffixSum += float64(a[j])
       // number of elements in suffix
       s := int64(n - j)
       // remaining operations after using one per smallest removed
       o := m - int64(j)
       if o < 0 {
           continue
       }
       // max total add for this suffix
       add := k * s
       if add > o {
           add = o
       }
       candidate := (suffixSum + float64(add)) / float64(s)
       if candidate > ans {
           ans = candidate
       }
   }
   // output result
   writer := bufio.NewWriter(os.Stdout)
   fmt.Fprintf(writer, "%.10f", ans)
   writer.Flush()
}
