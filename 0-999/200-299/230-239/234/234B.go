package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// This program reads n and k, then n integers, and outputs the k-th largest value
// and the 1-based indices of the k largest values in the original order of values sorted by value.
func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   type Pair struct { value, idx int }
   a := make([]Pair, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i].value)
       a[i].idx = i + 1
   }
   sort.Slice(a, func(i, j int) bool {
       return a[i].value < a[j].value
   })
   // k-th largest is at index n-k after sorting ascending
   threshold := a[n-k].value
   fmt.Fprintln(writer, threshold)
   for i := n - k; i < n; i++ {
       fmt.Fprint(writer, a[i].idx)
       if i < n-1 {
           fmt.Fprint(writer, " ")
       }
   }
   fmt.Fprintln(writer)
}
