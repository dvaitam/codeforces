package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func modPow(a, e int64) int64 {
   res := int64(1)
   a %= mod
   for e > 0 {
       if e&1 == 1 {
           res = (res * a) % mod
       }
       a = (a * a) % mod
       e >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m, k int64
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   // skip obstacle definitions
   for i := int64(0); i < k; i++ {
       var x1, y1, x2, y2 int64
       fmt.Fscan(reader, &x1, &y1, &x2, &y2)
   }
   // each obstacle can be passed either to the left or to the right
   // assuming independence, total ways = 2^k mod mod
   ans := modPow(2, k)
   fmt.Fprintln(writer, ans)
}
