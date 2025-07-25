package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var q int
   fmt.Fscan(reader, &q)
   const pat = "RGB"
   for qi := 0; qi < q; qi++ {
       var n, k int
       var s string
       fmt.Fscan(reader, &n, &k)
       fmt.Fscan(reader, &s)
       ans := k
       // Try all 3 possible alignments
       for offset := 0; offset < 3; offset++ {
           // prefix sum of mismatches
           pref := make([]int, n+1)
           for i := 0; i < n; i++ {
               exp := pat[(i+offset)%3]
               if s[i] != exp {
                   pref[i+1] = pref[i] + 1
               } else {
                   pref[i+1] = pref[i]
               }
           }
           // slide window of length k
           for i := k; i <= n; i++ {
               changes := pref[i] - pref[i-k]
               if changes < ans {
                   ans = changes
               }
           }
       }
       fmt.Fprintln(writer, ans)
   }
}
