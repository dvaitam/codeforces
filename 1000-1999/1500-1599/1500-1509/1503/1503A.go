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
       fmt.Fscan(in, &n)
       var s string
       fmt.Fscan(in, &s)
       ones := 0
       for i := 0; i < n; i++ {
           if s[i] == '1' {
               ones++
           }
       }
       a := make([]byte, n)
       b := make([]byte, n)
       half := ones / 2
       sw := true
       for i := 0; i < n; i++ {
           if s[i] == '1' {
               if half > 0 {
                   a[i], b[i] = '(', '('
               } else {
                   a[i], b[i] = ')', ')'
               }
               half--
           } else {
               if sw {
                   a[i], b[i] = '(', ')'
               } else {
                   a[i], b[i] = ')', '('
               }
               sw = !sw
           }
       }
       valid := true
       balA, balB := 0, 0
       for i := 0; i < n; i++ {
           if a[i] == '(' {
               balA++
           } else {
               balA--
           }
           if b[i] == '(' {
               balB++
           } else {
               balB--
           }
           if balA < 0 || balB < 0 {
               valid = false
               break
           }
       }
       if !valid {
           fmt.Fprintln(out, "NO")
           continue
       }
       fmt.Fprintln(out, "YES")
       out.Write(a)
       fmt.Fprintln(out)
       out.Write(b)
       fmt.Fprintln(out)
   }
}
