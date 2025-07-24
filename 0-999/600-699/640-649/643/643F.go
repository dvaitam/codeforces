package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n uint64
   var p, q int
   if _, err := fmt.Fscan(reader, &n, &p, &q); err != nil {
       return
   }
   // m = min(p, n-1)
   var m int
   if n == 0 {
       m = 0
   } else {
       nm1 := int(n - 1)
       if p < nm1 {
           m = p
       } else {
           m = nm1
       }
   }
   // compute S = sum_{k=1..m} C(n, k)
   Smod := uint64(0)
   if m > 0 {
       sum := new(big.Int)
       cur := big.NewInt(1)
       nn := new(big.Int).SetUint64(n)
       for k := 1; k <= m; k++ {
           // cur = cur * (n - k + 1) / k
           term := new(big.Int).Sub(nn, big.NewInt(int64(k-1)))
           cur.Mul(cur, term)
           cur.Div(cur, big.NewInt(int64(k)))
           sum.Add(sum, cur)
       }
       // mod 2^32
       Smod = sum.Uint64() & 0xffffffff
   }
   const mod = uint64(1) << 32
   var ans uint32
   for i := 1; i <= q; i++ {
       s := uint64(i)
       // X_i = i + S * i^2  (mod 2^32)
       t := (s * s) & (mod - 1)
       t = (t * Smod) & (mod - 1)
       x := (s + t) & (mod - 1)
       ans ^= uint32(x)
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprint(writer, ans)
}
