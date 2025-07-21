package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s, t string
   if _, err := fmt.Fscan(reader, &s, &t); err != nil {
       return
   }
   n, m := len(s), len(t)
   // count of s
   cnt0 := [26]int{}
   for i := 0; i < n; i++ {
       cnt0[s[i]-'a']++
   }
   maxPre := m
   if n < m {
       maxPre = n
   }
   // cnt[j]: counts after consuming t[0..j-1]
   cnt := make([][26]int, maxPre+1)
   ok := make([]bool, maxPre+1)
   cnt[0] = cnt0
   ok[0] = true
   for j := 0; j < maxPre; j++ {
       if !ok[j] {
           break
       }
       c := t[j] - 'a'
       if cnt[j][c] > 0 {
           cnt[j+1] = cnt[j]
           cnt[j+1][c]--
           ok[j+1] = true
       } else {
           break
       }
   }
   // Case: prefix matches full t and s longer -> use t as prefix
   if n > m && ok[m] {
       res := make([]byte, n)
       // prefix equal t
       for i := 0; i < m; i++ {
           res[i] = t[i]
       }
       // position m: smallest available
       for k := 0; k < 26; k++ {
           if cnt[m][k] > 0 {
               res[m] = byte('a' + k)
               cnt[m][k]--
               break
           }
       }
       // fill rest
       idx := m + 1
       for k := 0; k < 26; k++ {
           for cnt[m][k] > 0 {
               res[idx] = byte('a' + k)
               cnt[m][k]--
               idx++
           }
       }
       fmt.Println(string(res))
       return
   }
   // general: find largest j where we can increase
   for j := maxPre; j >= 1; j-- {
       if !ok[j-1] {
           continue
       }
       baseCnt := cnt[j-1]
       x := t[j-1] - 'a'
       // find smallest k > x
       var kSel int = -1
       for k := int(x) + 1; k < 26; k++ {
           if baseCnt[k] > 0 {
               kSel = k
               break
           }
       }
       if kSel < 0 {
           continue
       }
       // build result
       res := make([]byte, n)
       // prefix t[0..j-2]
       for i := 0; i < j-1; i++ {
           res[i] = t[i]
       }
       // position j-1: choose kSel
       res[j-1] = byte('a' + kSel)
       baseCnt[kSel]--
       // fill rest positions
       idx := j
       for kk := 0; kk < 26; kk++ {
           for baseCnt[kk] > 0 {
               res[idx] = byte('a' + kk)
               baseCnt[kk]--
               idx++
           }
       }
       fmt.Println(string(res))
       return
   }
   // no solution
   fmt.Println("-1")
}
