package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   s, err := reader.ReadString('\n')
   if err != nil && len(s) == 0 {
       return
   }
   // Remove trailing newline or carriage return
   if s[len(s)-1] == '\n' {
       s = s[:len(s)-1]
   }
   if len(s) > 0 && s[len(s)-1] == '\r' {
       s = s[:len(s)-1]
   }
   ans := uint64(1)
   n := len(s)
   for i := 0; i < n; {
       if s[i] == '9' {
           i++
           continue
       }
       // check if start of a sum-to-9 run
       if i+1 < n && (s[i]-'0')+(s[i+1]-'0') == 9 {
           // find end of run
           j := i
           for j+1 < n && s[j] != '9' && (s[j]-'0')+(s[j+1]-'0') == 9 {
               j++
           }
           length := j - i + 1 // number of digits in run
           if length%2 == 1 {
               // odd length => number of max matchings = floor(length/2)+1
               ans *= uint64(length/2 + 1)
           }
           // even length contributes factor 1
           i = j + 1
       } else {
           i++
       }
   }
   fmt.Println(ans)
}
