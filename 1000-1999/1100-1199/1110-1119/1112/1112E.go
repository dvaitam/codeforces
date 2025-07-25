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

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]byte, n)
   b := make([]byte, n)
   // read strings
   var sa, sb string
   fmt.Fscan(in, &sa)
   fmt.Fscan(in, &sb)
   for i := 0; i < n; i++ {
       a[i] = sa[i] - '0'
       b[i] = sb[i] - '0'
   }
   // compute differences D and ops
   D := make([]int, n)
   for i := 0; i < n; i++ {
       D[i] = int(b[i]) - int(a[i])
   }
   ops := make([]int, n-1)
   if n >= 1 {
       ops[0] = D[0]
   }
   for i := 1; i < n-1; i++ {
       ops[i] = D[i] - ops[i-1]
   }
   // check last
   if n > 1 && ops[n-2] != D[n-1] {
       fmt.Fprintln(out, -1)
       return
   }
   // sum moves
   var total int64 = 0
   for i := 0; i < n-1; i++ {
       if ops[i] < 0 {
           total -= int64(ops[i])
       } else {
           total += int64(ops[i])
       }
   }
   fmt.Fprintln(out, total)
   // output moves up to limit
   const maxPrint = 100000
   printed := 0
   for i := 0; i < n-1 && printed < maxPrint; i++ {
       cnt := ops[i]
       var s int
       if cnt > 0 {
           s = 1
       } else {
           s = -1
       }
       for cnt != 0 && printed < maxPrint {
           fmt.Fprintf(out, "%d %d\n", i+1, s)
           printed++
           if s > 0 {
               cnt--
           } else {
               cnt++
           }
       }
   }
}
