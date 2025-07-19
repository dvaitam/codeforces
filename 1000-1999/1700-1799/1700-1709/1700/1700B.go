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
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(in, &n)
       var s string
       fmt.Fscan(in, &s)
       v := make([]int, n)
       for i := 0; i < n; i++ {
           v[i] = int(s[i] - '0')
       }
       for i := 0; i < n; i++ {
           v[i] = 9 - v[i]
       }
       if v[0] == 0 {
           carry := 0
           d := v[n-1] + 2
           carry = d / 10
           v[n-1] = d % 10
           for i := n - 2; i >= 0; i-- {
               d = v[i] + 3 + carry
               carry = d / 10
               v[i] = d % 10
           }
       }
       for i := 0; i < n; i++ {
           out.WriteByte(byte(v[i] + '0'))
       }
       out.WriteByte('\n')
   }
}
