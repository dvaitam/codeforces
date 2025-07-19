package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func max(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   maxlen := make([]int64, 26)
   for k := 1; k <= n; k++ {
       var s string
       fmt.Fscan(reader, &s)
       m := len(s)
       if k != 1 {
           // compute prefix same count
           first := s[0] - 'a'
           last := s[m-1] - 'a'
           l := 1
           for l < m && s[l] == s[0] {
               l++
           }
           // compute suffix same count
           rcount := 1
           for i := m - 2; i >= 0 && s[i] == s[m-1]; i-- {
               rcount++
           }
           if l == m {
               // all characters same
               for i := 0; i < 26; i++ {
                   if int(s[0]-'a') != i {
                       maxlen[i] = min(maxlen[i], 1)
                   }
               }
               c := int(s[0] - 'a')
               if maxlen[c] > 0 {
                   maxlen[c] = (maxlen[c] + 1) * int64(m+1) - 1
               } else {
                   maxlen[c] = int64(m)
               }
           } else {
               for i := 0; i < 26; i++ {
                   maxlen[i] = min(maxlen[i], 1)
               }
               maxlen[first] += int64(l)
               maxlen[last] += int64(rcount)
           }
       }
       // update from runs in s
       var i = 0
       for i < m {
           j := i + 1
           for j < m && s[j] == s[i] {
               j++
           }
           c := int(s[i] - 'a')
           run := int64(j - i)
           if run > maxlen[c] {
               maxlen[c] = run
           }
           i = j
       }
   }
   var ans int64
   for i := 0; i < 26; i++ {
       if maxlen[i] > ans {
           ans = maxlen[i]
       }
   }
   fmt.Println(ans)
}
