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

   var n, p int
   if _, err := fmt.Fscan(in, &n, &p); err != nil {
       return
   }
   var s string
   fmt.Fscan(in, &s)
   a := []byte(s)

   possible := false
   idx := -1
   for i := 0; i < n-p; i++ {
       x, y := a[i], a[i+p]
       if x == '.' || y == '.' || (x != y) {
           possible = true
           idx = i
           break
       }
   }
   if !possible {
       fmt.Fprintln(out, "No")
       return
   }
   // Make s[idx] and s[idx+p] a mismatch
   x, y := a[idx], a[idx+p]
   if x == '.' && y == '.' {
       a[idx] = '0'
       a[idx+p] = '1'
   } else if x == '.' {
       // y is '0' or '1'
       if y == '0' {
           a[idx] = '1'
       } else {
           a[idx] = '0'
       }
   } else if y == '.' {
       // x is '0' or '1'
       if x == '0' {
           a[idx+p] = '1'
       } else {
           a[idx+p] = '0'
       }
   }
   // Fill remaining '.' with '0'
   for i := 0; i < n; i++ {
       if a[i] == '.' {
           a[i] = '0'
       }
   }
   fmt.Fprintln(out, string(a))
}
