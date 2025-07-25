package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var s string
   fmt.Fscan(reader, &s)
   // length must be even
   if n%2 != 0 {
       fmt.Fprintln(writer, ":(")
       return
   }
   // invalid fixed endpoints
   if s[0] == ')' || s[n-1] == '(' {
       fmt.Fprintln(writer, ":(")
       return
   }
   // convert to byte slice for mutation
   b := []byte(s)
   // first must be '(', last must be ')'
   b[0] = '('
   b[n-1] = ')'
   // count existing
   totalOpen := 0
   totalClose := 0
   for i := 0; i < n; i++ {
       if b[i] == '(' {
           totalOpen++
       } else if b[i] == ')' {
           totalClose++
       }
   }
   need := n / 2
   remOpen := need - totalOpen
   remClose := need - totalClose
   if remOpen < 0 || remClose < 0 {
       fmt.Fprintln(writer, ":(")
       return
   }
   // assign remaining
   for i := 1; i < n-1; i++ {
       if b[i] == '?' {
           if remOpen > 0 {
               b[i] = '('
               remOpen--
           } else {
               b[i] = ')'
               remClose--
           }
       }
   }
   // validate
   bal := 0
   for i := 0; i < n; i++ {
       if b[i] == '(' {
           bal++
       } else {
           bal--
       }
       // before last, balance must be >=1
       if i < n-1 {
           if bal <= 0 {
               fmt.Fprintln(writer, ":(")
               return
           }
       } else {
           // at end must be zero
           if bal != 0 {
               fmt.Fprintln(writer, ":(")
               return
           }
       }
   }
   // all good
   fmt.Fprintln(writer, string(b))
}
