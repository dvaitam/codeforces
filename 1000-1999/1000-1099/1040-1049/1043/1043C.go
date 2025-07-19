package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   // convert string to mutable byte slice
   b := []byte(s)
   v := make([]int, n)
   // v[0] is 0 by default
   for i := 1; i < n; i++ {
       rev := false
       if i == n-1 {
           if b[i] == 'a' {
               rev = true
           }
       } else {
           if b[i] != b[i+1] {
               rev = true
           }
       }
       if rev {
           v[i] = 1
           // reverse prefix [0..i]
           for l, r := 0, i; l < r; l, r = l+1, r-1 {
               b[l], b[r] = b[r], b[l]
           }
       }
   }
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i, x := range v {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(fmt.Sprintf("%d", x))
   }
   writer.WriteByte('\n')
}
