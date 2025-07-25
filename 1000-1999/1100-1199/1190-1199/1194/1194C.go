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
   var q int
   fmt.Fscan(in, &q)
   for ; q > 0; q-- {
       var s, t, p string
       fmt.Fscan(in, &s, &t, &p)
       if canTransform(s, t, p) {
           fmt.Fprintln(out, "YES")
       } else {
           fmt.Fprintln(out, "NO")
       }
   }
}

// canTransform checks if string s can be transformed into t by inserting characters from p
func canTransform(s, t, p string) bool {
   i := 0
   // need[c] is number of extra c's required from p
   var need [26]int
   for j := 0; j < len(t); j++ {
       if i < len(s) && s[i] == t[j] {
           i++
       } else {
           need[t[j]-'a']++
       }
   }
   // s must be subsequence of t
   if i != len(s) {
       return false
   }
   // count available characters in p
   var have [26]int
   for j := 0; j < len(p); j++ {
       have[p[j]-'a']++
   }
   // check if available meets needed
   for c := 0; c < 26; c++ {
       if need[c] > have[c] {
           return false
       }
   }
   return true
}
