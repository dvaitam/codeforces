// Package main implements the solution for Problem D of contest 309.
// Given an integer n, compute the n-th Fibonacci number (with F(0)=0, F(1)=1).
package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   switch {
   case n <= 0:
       fmt.Println(0)
       return
   case n == 1:
       fmt.Println(1)
       return
   }
   a := big.NewInt(0)
   b := big.NewInt(1)
   for i := 2; i <= n; i++ {
       tmp := new(big.Int).Add(a, b)
       a.Set(b)
       b.Set(tmp)
   }
   fmt.Println(b)
}
