package main

import (
   "fmt"
   "math/big"
)

func main() {
   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   a := new(big.Int)
   a.SetString(s, 10)

   // now = 1e17 % a
   ten17 := new(big.Int)
   ten17.SetString("100000000000000000", 10)
   now := new(big.Int).Mod(ten17, a)

   // P = sum_{i=1..9} now*i mod a
   P := new(big.Int)
   for i := int64(1); i <= 9; i++ {
       term := new(big.Int).Mul(now, big.NewInt(i))
       term.Mod(term, a)
       P.Add(P, term)
       P.Mod(P, a)
   }

   // now = 18*P mod a
   now.SetInt64(0)
   for i := 1; i <= 18; i++ {
       now.Add(now, P)
       now.Mod(now, a)
   }
   // now = (now + 1) mod a
   now.Add(now, big.NewInt(1))
   now.Mod(now, a)

   // diff = (a - now) mod a
   diff := new(big.Int).Sub(a, now)
   diff.Mod(diff, a)

   // base = 1e18
   base := new(big.Int)
   base.SetString("1000000000000000000", 10)
   // l = base + diff, r = base + diff
   l := new(big.Int).Add(base, diff)
   r := new(big.Int).Add(base, diff)

   fmt.Printf("%s %s", l.String(), r.String())
}
