package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var aStr, bStr string
   if _, err := fmt.Fscan(reader, &aStr); err != nil {
       return
   }
   if _, err := fmt.Fscan(reader, &bStr); err != nil {
       return
   }
   a := new(big.Int)
   b := new(big.Int)
   a.SetString(aStr, 10)
   b.SetString(bStr, 10)
   sum := new(big.Int).Add(a, b)
   fmt.Println(sum.String())
}
