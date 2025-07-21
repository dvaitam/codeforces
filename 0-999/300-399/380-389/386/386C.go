package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   fmt.Fscan(reader, &s)
   n := len(s)
   // compute total distinct letters in s
   seen := make([]bool, 26)
   d := 0
   for i := 0; i < n; i++ {
       idx := s[i] - 'a'
       if !seen[idx] {
           seen[idx] = true
           d++
       }
   }
   // f[k] = number of substrings with at most k distinct letters
   f := make([]int64, d+1)
   for k := 1; k <= d; k++ {
       var cnt [26]int
       distinct := 0
       l := 0
       var res int64
       for r := 0; r < n; r++ {
           idx := s[r] - 'a'
           if cnt[idx] == 0 {
               distinct++
           }
           cnt[idx]++
           for distinct > k {
               idxl := s[l] - 'a'
               cnt[idxl]--
               if cnt[idxl] == 0 {
                   distinct--
               }
               l++
           }
           res += int64(r - l + 1)
       }
       f[k] = res
   }
   // output results: t[k] = number of substrings with exactly k distinct letters
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, d)
   for k := 1; k <= d; k++ {
       var tk int64
       if k == 1 {
           tk = f[1]
       } else {
           tk = f[k] - f[k-1]
       }
       fmt.Fprintln(writer, tk)
   }
}
