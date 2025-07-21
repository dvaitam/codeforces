package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var nStr string
   if _, err := fmt.Fscan(reader, &nStr); err != nil {
       return
   }
   n := new(big.Int)
   n.SetString(nStr, 10)

   results := make(map[uint64]struct{})
   two := big.NewInt(2)
   eight := big.NewInt(8)
   // iterate over k such that 2^k up to reasonable limit
   for k := 0; k <= 60; k++ {
       // p = 1 << k
       p := new(big.Int).Lsh(big.NewInt(1), uint(k))
       // B = 2^(k+1) - 3
       B := new(big.Int).Lsh(big.NewInt(1), uint(k+1))
       B.Sub(B, big.NewInt(3))

       // D = B*B + 8*n
       D := new(big.Int).Mul(B, B)
       tmp := new(big.Int).Mul(eight, n)
       D.Add(D, tmp)
       // sqrtD
       sqrtD := new(big.Int).Sqrt(D)
       if new(big.Int).Mul(sqrtD, sqrtD).Cmp(D) != 0 {
           continue
       }
       // r = (sqrtD - B) / 2
       r := new(big.Int).Sub(sqrtD, B)
       if r.Sign() <= 0 || r.Bit(0) == 0 {
           continue
       }
       r.Div(r, two)
       if r.Sign() <= 0 {
           continue
       }
       // T = r * p
       T := new(big.Int).Mul(r, p)
       // ensure T fits in uint64
       if T.BitLen() > 64 {
           continue
       }
       t64 := T.Uint64()
       if t64 > 0 {
           results[t64] = struct{}{}
       }
   }

   // collect and sort
   var out []uint64
   for t := range results {
       out = append(out, t)
   }
   sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   if len(out) == 0 {
       fmt.Fprintln(writer, -1)
       return
   }
   for _, t := range out {
       fmt.Fprintln(writer, t)
   }
}
