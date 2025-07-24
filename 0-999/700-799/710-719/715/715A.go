package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // xprev holds the current number on screen before '+' press
   xprev := big.NewInt(2)
   one := big.NewInt(1)
   for k := 1; k <= n; k++ {
       // compute ceil(sqrt(xprev))
       sqrtFloor := new(big.Int).Sqrt(xprev)
       sqrtCeil := new(big.Int).Set(sqrtFloor)
       tmp := new(big.Int).Mul(sqrtFloor, sqrtFloor)
       if tmp.Cmp(xprev) < 0 {
           sqrtCeil.Add(sqrtCeil, one)
       }
       // m0 = ceil(sqrtCeil / (k+1))
       kp1 := big.NewInt(int64(k + 1))
       m0 := new(big.Int).Add(sqrtCeil, new(big.Int).Sub(kp1, one))
       m0.Div(m0, kp1)
       // m = ceil(m0 / k) * k
       bk := big.NewInt(int64(k))
       m := new(big.Int).Add(m0, new(big.Int).Sub(bk, one))
       m.Div(m, bk)
       m.Mul(m, bk)
       // t = m * (k+1)
       t := new(big.Int).Mul(m, kp1)
       // t2 = t * t
       t2 := new(big.Int).Mul(t, t)
       // diff = t2 - xprev
       diff := new(big.Int).Sub(t2, xprev)
       // a = diff / k
       a := new(big.Int).Div(diff, bk)
       // print a
       fmt.Fprintln(writer, a)
       // update xprev = t
       xprev.Set(t)
   }
}
