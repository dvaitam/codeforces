package main

import (
   "fmt"
  "sort"
)

func main() {
   var n int
   var k int64
   if _, err := fmt.Scan(&n, &k); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&a[i])
   }
   sort.Ints(a)
   // compute zero-based indices
   idx1 := int((k-1) / int64(n))
   idx2 := int((k-1) % int64(n))
   fmt.Println(a[idx1], a[idx2])
}
