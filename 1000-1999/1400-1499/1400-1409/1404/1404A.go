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
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n, k int
       var s string
       fmt.Fscan(in, &n, &k)
       fmt.Fscan(in, &s)
       // tarr holds forced values for each position mod k
       tarr := make([]byte, k)
       for i := 0; i < k; i++ {
           tarr[i] = '?'
       }
       ok := true
       // propagate constraints
       for i := 0; i < n; i++ {
           c := s[i]
           if c == '?' {
               continue
           }
           pos := i % k
           if tarr[pos] == '?' {
               tarr[pos] = c
           } else if tarr[pos] != c {
               ok = false
               break
           }
       }
       if !ok {
           fmt.Fprintln(out, "NO")
           continue
       }
       // count fixed zeros and ones
       cnt0, cnt1 := 0, 0
       for i := 0; i < k; i++ {
           if tarr[i] == '0' {
               cnt0++
           } else if tarr[i] == '1' {
               cnt1++
           }
       }
       // must be able to place remaining to balance
       if cnt0 > k/2 || cnt1 > k/2 {
           fmt.Fprintln(out, "NO")
       } else {
           fmt.Fprintln(out, "YES")
       }
   }
}
