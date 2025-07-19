package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   rdr := bufio.NewReader(os.Stdin)
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()

   var n, m int
   fmt.Fscan(rdr, &n, &m)
   var a, b string
   fmt.Fscan(rdr, &a, &b)

   const mod int64 = 998244353
   s := make([]int64, m+1)
   for i := 1; i <= m; i++ {
      s[i] = s[i-1] + int64(b[i-1]-'0')
   }

   var ans, mul int64 = 0, 1
   for i := n - 1; i >= 0; i-- {
      pos := n - 1 - i
      if a[i] == '1' && pos < m {
         ans = (ans + mul*s[m-pos]) % mod
      }
      mul = (mul * 2) % mod
   }

   fmt.Fprint(w, ans)
}
