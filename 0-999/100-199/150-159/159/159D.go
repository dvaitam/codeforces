package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   // count of palindromic substrings ending at i and starting at i
   countEnd := make([]int64, n)
   countStart := make([]int64, n)
   // expand around centers for odd-length palindromes
   for center := 0; center < n; center++ {
       l, r := center, center
       for l >= 0 && r < n && s[l] == s[r] {
           countStart[l]++
           countEnd[r]++
           l--
           r++
       }
   }
   // even-length palindromes
   for center := 0; center < n-1; center++ {
       l, r := center, center+1
       for l >= 0 && r < n && s[l] == s[r] {
           countStart[l]++
           countEnd[r]++
           l--
           r++
       }
   }
   // suffix sums of countStart
   suffix := make([]int64, n+1)
   for i := n - 1; i >= 0; i-- {
       suffix[i] = suffix[i+1] + countStart[i]
   }
   // compute answer: for each end index b, count pairs with start x > b
   var ans int64
   for b := 0; b < n-1; b++ {
       if countEnd[b] > 0 {
           ans += countEnd[b] * suffix[b+1]
       }
   }
   // output
   fmt.Println(ans)
}
