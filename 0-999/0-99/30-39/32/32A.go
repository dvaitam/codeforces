package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var d int64
   if _, err := fmt.Fscan(in, &n, &d); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
   var count int64
   r := 0
   for i := 0; i < n; i++ {
       if r < i+1 {
           r = i + 1
       }
       for r < n && a[r]-a[i] <= d {
           r++
       }
       // now r is first index >i where a[r]-a[i] > d, so valid j are i+1..r-1
       cnt := int64(r - i - 1)
       if cnt > 0 {
           count += cnt
       }
   }
   // ordered pairs: each unordered pair counts twice
   fmt.Println(count * 2)
}
