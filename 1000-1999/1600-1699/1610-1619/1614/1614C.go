package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod int64 = 1000000007

func pow2(exp int64) int64 {
   result := int64(1)
   base := int64(2)
   for exp > 0 {
      if exp&1 == 1 {
         result = result * base % mod
      }
      base = base * base % mod
      exp >>= 1
   }
   return result
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
      return
   }
   for ; t > 0; t-- {
      var n, m int64
      fmt.Fscan(reader, &n, &m)
      var total int64
      for i := int64(0); i < m; i++ {
         var l, r, x int64
         fmt.Fscan(reader, &l, &r, &x)
         total |= x
      }
      ans := total % mod * pow2(n-1) % mod
      fmt.Fprintln(writer, ans)
   }
}
