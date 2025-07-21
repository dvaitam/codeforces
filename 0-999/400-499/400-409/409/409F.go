package main

import (
   "fmt"
   "math/big"
   "os"
)

// Reads an integer a (1 ≤ a ≤ 64) and prints a! (factorial of a) as a big integer.
func main() {
   var a int
   if _, err := fmt.Fscan(os.Stdin, &a); err != nil {
       return
   }
   result := big.NewInt(1)
   for i := 2; i <= a; i++ {
       result.Mul(result, big.NewInt(int64(i)))
   }
   fmt.Println(result)
}
