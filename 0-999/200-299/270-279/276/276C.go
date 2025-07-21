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

   var n, q int
   fmt.Fscan(reader, &n, &q)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   // Difference array for query frequencies
   freqDiff := make([]int, n+1)
   for i := 0; i < q; i++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       freqDiff[l-1]++
       freqDiff[r]--
   }

   // Build actual frequencies
   freq := make([]int, n)
   curr := 0
   for i := 0; i < n; i++ {
       curr += freqDiff[i]
       freq[i] = curr
   }

   // Sort array values and frequencies
   sort.Ints(a)
   sort.Ints(freq)

   // Compute maximum sum
   var ans int64
   for i := 0; i < n; i++ {
       ans += int64(a[i]) * int64(freq[i])
   }

   fmt.Fprint(writer, ans)
}
