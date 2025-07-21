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

   var s string
   fmt.Fscan(in, &s)
   n := len(s)
   // Z-function
   z := make([]int, n)
   l, r := 0, 0
   for i := 1; i < n; i++ {
       if i <= r {
           if r-i+1 < z[i-l] {
               z[i] = r - i + 1
           } else {
               z[i] = z[i-l]
           }
       }
       for i+z[i] < n && s[z[i]] == s[i+z[i]] {
           z[i]++
       }
       if i+z[i]-1 > r {
           l = i
           r = i + z[i] - 1
       }
   }
   // count occurrences of prefixes
   cnt := make([]int, n+2)
   for i := 1; i < n; i++ {
       cnt[1]++
       cnt[z[i]+1]--
   }
   for i := 1; i <= n; i++ {
       cnt[i] += cnt[i-1]
       // include prefix at position 0
       cnt[i]++
   }
   // prefix-function for borders
   pi := make([]int, n)
   for i := 1; i < n; i++ {
       j := pi[i-1]
       for j > 0 && s[i] != s[j] {
           j = pi[j-1]
       }
       if s[i] == s[j] {
           j++
       }
       pi[i] = j
   }
   // collect border lengths
   var borders []int
   k := n
   for k > 0 {
       borders = append(borders, k)
       if k == n {
           k = pi[n-1]
       } else {
           k = pi[k-1]
       }
   }
   // reverse to increasing order
   for i, j := 0, len(borders)-1; i < j; i, j = i+1, j-1 {
       borders[i], borders[j] = borders[j], borders[i]
   }
   // output
   fmt.Fprintln(out, len(borders))
   for _, length := range borders {
       fmt.Fprintln(out, length, cnt[length])
   }
}
