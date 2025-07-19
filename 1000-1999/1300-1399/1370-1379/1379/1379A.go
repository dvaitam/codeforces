package main

import (
   "bufio"
   "fmt"
   "os"
)

func countOcc(s string) int {
   cnt := 0
   pat := "abacaba"
   m := len(pat)
   for i := 0; i+m <= len(s); i++ {
       if s[i:i+m] == pat {
           cnt++
       }
   }
   return cnt
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   fmt.Fscan(in, &t)
   pat := "abacaba"
   for t > 0 {
       t--
       var n int
       var s string
       fmt.Fscan(in, &n)
       fmt.Fscan(in, &s)
       found := false
       var res string
       // Try every position to place pattern
       for i := 0; i+len(pat) <= n; i++ {
           ok := true
           for j := 0; j < len(pat); j++ {
               if s[i+j] != '?' && s[i+j] != pat[j] {
                   ok = false
                   break
               }
           }
           if !ok {
               continue
           }
           tmp := []byte(s)
           for j := 0; j < len(pat); j++ {
               tmp[i+j] = pat[j]
           }
           for k := 0; k < n; k++ {
               if tmp[k] == '?' {
                   tmp[k] = 'z'
               }
           }
           if countOcc(string(tmp)) == 1 {
               res = string(tmp)
               found = true
               break
           }
       }
       if found {
           fmt.Fprintln(out, "YES")
           fmt.Fprintln(out, res)
       } else {
           fmt.Fprintln(out, "NO")
       }
   }
}
