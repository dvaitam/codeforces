package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a int
   if _, err := fmt.Fscan(reader, &a); err != nil {
       return
   }
   // compute factorial of a (1 ≤ a ≤ 40)
   result := big.NewInt(1)
   for i := 2; i <= a; i++ {
       result.Mul(result, big.NewInt(int64(i)))
   }
   // output as string
   fmt.Println(result.String())
}
