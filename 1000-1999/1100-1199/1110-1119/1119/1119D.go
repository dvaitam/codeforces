package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
)

func main() {
   rdr := bufio.NewReader(os.Stdin)
   wtr := bufio.NewWriter(os.Stdout)
   defer wtr.Flush()

   var n int
   fmt.Fscan(rdr, &n)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(rdr, &a[i])
   }
   sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
   b := make([]int64, n-1)
   for i := 0; i < n-1; i++ {
       b[i] = a[i+1] - a[i]
   }
   sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
   ps := make([]int64, n)
   for i := 1; i < n; i++ {
       ps[i] = ps[i-1] + b[i-1]
   }
   var q int
   fmt.Fscan(rdr, &q)
   for i := 0; i < q; i++ {
       var l, r int64
       fmt.Fscan(rdr, &l, &r)
       k := r - l + 1
       idx := sort.Search(len(b), func(i int) bool { return b[i] > k })
       ans := ps[idx] + int64(n-idx)*k
       wtr.WriteString(strconv.FormatInt(ans, 10))
       wtr.WriteByte(' ')
   }
   wtr.WriteByte('\n')
}
