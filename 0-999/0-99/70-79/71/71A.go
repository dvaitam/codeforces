package main

import (
   "bufio"
   "fmt"
   "os"
)

// Solution for Codeforces Problem A: Way Too Long Words
func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(in, &s)
       if len(s) > 10 {
           // abbreviate: first letter + count of middle letters + last letter
           fmt.Fprintf(out, "%c%d%c\n", s[0], len(s)-2, s[len(s)-1])
       } else {
           fmt.Fprintln(out, s)
       }
   }
}
