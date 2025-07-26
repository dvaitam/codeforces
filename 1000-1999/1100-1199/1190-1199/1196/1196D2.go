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
   // mapping: R->0, G->1, B->2
   conv := map[byte]int{'R': 0, 'G': 1, 'B': 2}
   for ; q > 0; q-- {
       var n, k int
       fmt.Fscan(reader, &n, &k)
       var s string
       fmt.Fscan(reader, &s)
       sbytes := []byte(s)
       // prefix sums of mismatches for each of 3 patterns
       // pref[c][i]: number of mismatches in s[0:i] (i chars) for pattern starting with c
       pref := make([][3]int, n+1)
       for i := 0; i < n; i++ {
           curr := conv[sbytes[i]]
           // copy previous
           pref[i+1] = pref[i]
           for c := 0; c < 3; c++ {
               expected := (c + i) % 3
               if curr != expected {
                   pref[i+1][c]++
               }
           }
       }
       ans := n
       // slide window of length k, start at l from 0 to n-k
       for l := 0; l + k <= n; l++ {
           r := l + k
           for c := 0; c < 3; c++ {
               // pattern c mismatches in [l, r)
               mismatches := pref[r][c] - pref[l][c]
               if mismatches < ans {
                   ans = mismatches
               }
           }
       }
       fmt.Fprintln(writer, ans)
   }
}
