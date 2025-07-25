package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // Map sequence values to indices 0..5
   idxMap := map[int]int{4: 0, 8: 1, 15: 2, 16: 3, 23: 4, 42: 5}
   // cnt[i]: number of subsequences currently at stage i
   cnt := make([]int, 6)
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(in, &x)
       idx := idxMap[x]
       if idx == 0 {
           // start a new subsequence
           cnt[0]++
       } else if cnt[idx-1] > 0 {
           // extend a subsequence at previous stage
           cnt[idx-1]--
           cnt[idx]++
       }
       // else: cannot use this element, it will be removed
   }
   // Number of complete sequences is cnt[5]
   complete := cnt[5]
   // Minimum removals = total elements - used elements
   removals := n - complete*6
   fmt.Fprintln(out, removals)
}
