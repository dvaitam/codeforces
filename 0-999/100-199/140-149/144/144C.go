package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   sLine, err := reader.ReadString('\n')
   if err != nil && len(sLine) == 0 {
       return
   }
   s := strings.TrimSpace(sLine)
   pLine, err := reader.ReadString('\n')
   if err != nil && len(pLine) == 0 {
       return
   }
   p := strings.TrimSpace(pLine)
   n := len(s)
   m := len(p)
   if m > n {
       fmt.Println(0)
       return
   }
   // frequency of p
   var countP [26]int
   for i := 0; i < m; i++ {
       countP[p[i]-'a']++
   }
   // sliding window frequency of s
   var countS [26]int
   question := 0
   for i := 0; i < m; i++ {
       ch := s[i]
       if ch == '?' {
           question++
       } else {
           countS[ch-'a']++
       }
   }
   ans := 0
   // check function
   check := func() bool {
       totalDef := 0
       for c := 0; c < 26; c++ {
           if countS[c] > countP[c] {
               return false
           }
           totalDef += countP[c] - countS[c]
       }
       return totalDef == question
   }
   if check() {
       ans++
   }
   // slide
   for i := m; i < n; i++ {
       // remove s[i-m]
       old := s[i-m]
       if old == '?' {
           question--
       } else {
           countS[old-'a']--
       }
       // add s[i]
       newc := s[i]
       if newc == '?' {
           question++
       } else {
           countS[newc-'a']++
       }
       if check() {
           ans++
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans)
}
