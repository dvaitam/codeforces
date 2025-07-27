package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var s string
       fmt.Fscan(reader, &s)
       n := len(s)
       dif := make([]int, n+1)
       for i := 0; i < n; i++ {
           dif[i+1] = dif[i]
           if s[i] == '0' {
               dif[i+1]++
           } else {
               dif[i+1]--
           }
       }
       // map diff value to positions of '1'
       mp := make(map[int][]int)
       for i := 0; i < n; i++ {
           if s[i] == '1' {
               mp[dif[i+1]] = append(mp[dif[i+1]], i)
           }
       }
       ptr := make(map[int]int)
       matchedRight := make([]bool, n)
       ans := make([]byte, 0, n)
       i := 0
       for i < n {
           if matchedRight[i] {
               ans = append(ans, '0')
               i++
               continue
           }
           if s[i] == '0' {
               ans = append(ans, '0')
               i++
               continue
           }
           d := dif[i]
           lst := mp[d]
           p := ptr[d]
           for p < len(lst) && (lst[p] <= i || matchedRight[lst[p]]) {
               p++
           }
           ptr[d] = p
           if p < len(lst) {
               j := lst[p]
               ans = append(ans, '0')
               matchedRight[j] = true
               ptr[d] = p + 1
               i++
               continue
           }
           // no further match, append remainder
           for k := i; k < n; k++ {
               if matchedRight[k] {
                   ans = append(ans, '0')
               } else {
                   ans = append(ans, s[k])
               }
           }
           break
       }
       fmt.Println(string(ans))
   }
}
