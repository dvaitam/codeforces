package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   // Count total hamsters (H)
   totalH := 0
   n = len(s)
   bad := make([]int, n)
   for i, ch := range s {
       if ch == 'H' {
           totalH++
       } else if ch == 'T' {
           bad[i] = 1
       }
   }
   // Sliding window of size totalH on circle: count T's in window
   // totalH >=1 and <= n-1 as at least one of each present
   m := totalH
   // initial window [0..m-1]
   curr := 0
   for i := 0; i < m; i++ {
       curr += bad[i]
   }
   minSwaps := curr
   // slide window start from 1 to n-1
   for i := 1; i < n; i++ {
       // remove index i-1, add index (i+m-1)%n
       curr -= bad[i-1]
       addIdx := i + m - 1
       if addIdx >= n {
           addIdx -= n
       }
       curr += bad[addIdx]
       if curr < minSwaps {
           minSwaps = curr
       }
   }
   fmt.Println(minSwaps)
}
