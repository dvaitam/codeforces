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

   var t int
   fmt.Fscan(in, &t)
   for t > 0 {
       t--
       var n int
       var b string
       fmt.Fscan(in, &n, &b)
       // determine distinct sorted letters
       var present [26]bool
       for i := 0; i < n; i++ {
           present[b[i]-'a'] = true
       }
       var letters []byte
       for c := 0; c < 26; c++ {
           if present[c] {
               letters = append(letters, byte('a'+c))
           }
       }
       // build mapping: letters[i] -> letters[len-1-i]
       m := [26]byte{}
       L := len(letters)
       for i, c := range letters {
           m[c-'a'] = letters[L-1-i]
       }
       // decode
       res := make([]byte, n)
       for i := 0; i < n; i++ {
           res[i] = m[b[i]-'a']
       }
       fmt.Fprintln(out, string(res))
   }
}
