package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   _, err := fmt.Fscan(reader, &s)
   if err != nil {
       fmt.Fprintln(os.Stderr, "failed to read input:", err)
       return
   }
   n := new(big.Int)
   // parse integer in base 10
   if _, ok := n.SetString(s, 10); !ok {
       fmt.Fprintln(os.Stderr, "invalid integer input")
       return
   }
   // zero and one are not prime
   one := big.NewInt(1)
   if n.Cmp(one) <= 0 {
       fmt.Println("composite")
       return
   }
   // use Miller-Rabin primality test with sufficient rounds
   // ProbablyPrime reports whether n is prime
   // 20 rounds gives very high accuracy for large numbers
   if n.ProbablyPrime(20) {
       fmt.Println("prime")
   } else {
       fmt.Println("composite")
   }
}
