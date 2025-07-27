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
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       l := make([]int, n)
       r := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &l[i], &r[i])
       }
       L := make([]int, n)
       R := make([]int, n)
       copy(L, l)
       copy(R, r)
       sort.Ints(L)
       sort.Ints(R)
       maxc := 0
       for i := 0; i < n; i++ {
           // count segments with left <= r[i]
           x := sort.Search(n, func(j int) bool { return L[j] > r[i] })
           // count segments with right < l[i]
           y := sort.Search(n, func(j int) bool { return R[j] >= l[i] })
           c := x - y
           if c > maxc {
               maxc = c
           }
       }
       // minimum deletions = n - maxc
       fmt.Fprintln(writer, n-maxc)
   }
}
