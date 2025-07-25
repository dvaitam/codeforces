package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
   // count duplicates
   dupCount := 0
   dupIdx := -1
   for i := 1; i < n; i++ {
       if a[i] == a[i-1] {
           dupCount++
           dupIdx = i
       }
   }
   // invalid if more than one duplicate
   if dupCount > 1 {
       fmt.Fprintln(out, "cslnb")
       return
   }
   if dupCount == 1 {
       v := a[dupIdx]
       // can't have duplicate zeros
       if v == 0 {
           fmt.Fprintln(out, "cslnb")
           return
       }
       // can't reduce one duplicate to v-1 if v-1 already exists
       target := v - 1
       // binary search for target
       idx := sort.Search(len(a), func(i int) bool { return a[i] >= target })
       if idx < n && a[idx] == target {
           fmt.Fprintln(out, "cslnb")
           return
       }
   }
   // compute total safe moves to reach [0,1,2,...]
   var moves int64
   for i := 0; i < n; i++ {
       moves += a[i] - int64(i)
   }
   if moves%2 == 1 {
       fmt.Fprintln(out, "sjfnb")
   } else {
       fmt.Fprintln(out, "cslnb")
   }
}
