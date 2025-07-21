package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var pStr, qStr string
   if _, err := fmt.Fscan(reader, &pStr, &qStr); err != nil {
       return
   }
   p := new(big.Int)
   if _, ok := p.SetString(pStr, 10); !ok {
       return
   }
   q := new(big.Int)
   if _, ok := q.SetString(qStr, 10); !ok {
       return
   }
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]*big.Int, n)
   for i := 0; i < n; i++ {
       var aiStr string
       if _, err := fmt.Fscan(reader, &aiStr); err != nil {
           return
       }
       ai := new(big.Int)
       if _, ok := ai.SetString(aiStr, 10); !ok {
           return
       }
       a[i] = ai
   }
   // Build continued fraction value: numerator/denominator
   num := new(big.Int).Set(a[n-1])
   den := big.NewInt(1)
   for i := n - 2; i >= 0; i-- {
       newNum := new(big.Int).Mul(a[i], num)
       newNum.Add(newNum, den)
       newDen := new(big.Int).Set(num)
       num, den = newNum, newDen
   }
   // Compare p/q and num/den: p*den == q*num
   lhs := new(big.Int).Mul(p, den)
   rhs := new(big.Int).Mul(q, num)
   if lhs.Cmp(rhs) == 0 {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
