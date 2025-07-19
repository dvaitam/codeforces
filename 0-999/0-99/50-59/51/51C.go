package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

var n int
var a []int64
var st [3]int

// test checks if we can choose 3 elements in a with gaps <= x between consecutive picks
func test(x int64) bool {
   b := 0
   for i := 0; i < 3; i++ {
       st[i] = b
       key := a[b] + x
       // upper_bound: first index with a[j] > key
       pos := sort.Search(n, func(j int) bool { return a[j] > key })
       b = pos
       if b == n {
           return false
       }
   }
   return true
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   // read input
   fmt.Fscan(reader, &n)
   a = make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })

   // binary search on maximum x
   var l, r int64 = 0, 1000000000
   for l < r {
       mid := l + (r-l)/2
       if test(mid) {
           l = mid + 1
       } else {
           r = mid
       }
   }
   // set st for final value
   test(l)

   // output result: half of l, with 6 decimal places
   d := float64(l) / 2.0
   fmt.Fprintf(writer, "%.6f\n", d)
   // output positions: a[st[i]] + d
   for i := 0; i < 3; i++ {
       pos := float64(a[st[i]]) + d
       fmt.Fprintf(writer, "%.6f ", pos)
   }
   fmt.Fprintln(writer)
}
