package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   // read the concatenated string of digits
   fmt.Fscan(reader, &s)
   n := len(s)
   // cnt[r] = number of substrings ending at previous position with remainder r
   cnt := make([]int64, 11)
   ncnt := make([]int64, 11)
   var res int64
   for i := 0; i < n; i++ {
       d := int(s[i] - '0')
       // reset new counts
       for j := 0; j < 11; j++ {
           ncnt[j] = 0
       }
       // extend previous substrings
       for j := d + 1; j < 11; j++ {
           // compute new remainder after adding digit d
           // r' = (j*(j-1)/2 + d) mod 11
           nj := ((int64(j)*(int64(j)-1)/2) + int64(d) + 10) % 11
           ncnt[nj] += cnt[j]
       }
       // swap cnt and ncnt
       cnt, ncnt = ncnt, cnt
       // start new substring at i if digit non-zero
       if d != 0 {
           cnt[d]++
       }
       // accumulate count of inadequate substrings ending at i
       for j := 0; j < 11; j++ {
           res += cnt[j]
       }
   }
   // output result
   fmt.Println(res)
}
