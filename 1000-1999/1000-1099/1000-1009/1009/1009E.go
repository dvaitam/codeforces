package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 998244353

var rdr = bufio.NewReader(os.Stdin)
var wtr = bufio.NewWriter(os.Stdout)

// readInt reads the next integer from input
func readInt() (int, error) {
   var c byte
   var err error
   // skip non-digit characters
   for {
       c, err = rdr.ReadByte()
       if err != nil {
           return 0, err
       }
       if (c >= '0' && c <= '9') || c == '-' {
           break
       }
   }
   neg := false
   if c == '-' {
       neg = true
       c, err = rdr.ReadByte()
       if err != nil {
           return 0, err
       }
   }
   var x int
   for ; c >= '0' && c <= '9'; c, err = rdr.ReadByte() {
       if err != nil {
           break
       }
       x = x*10 + int(c-'0')
   }
   if neg {
       x = -x
   }
   return x, nil
}

func main() {
   defer wtr.Flush()
   nInt, err := readInt()
   if err != nil {
       return
   }
   n := nInt
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       ai, err := readInt()
       if err != nil {
           return
       }
       a[i] = int64(ai)
   }
   var ans, pow int64 = 0, 1
   // compute contributions from a[1] to a[n-1]
   for i := n - 1; i >= 1; i-- {
       term := a[i] * int64(n-i+2) % mod * pow % mod
       ans += term
       if ans >= mod {
           ans -= mod
       }
       pow = pow * 2 % mod
   }
   ans += a[n]
   ans %= mod
   fmt.Fprint(wtr, ans)
}
