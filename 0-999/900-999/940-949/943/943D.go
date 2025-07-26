package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // Read input integer as string
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       fmt.Fprintln(os.Stderr, "Error reading input:", err)
       os.Exit(1)
   }
   // Parse big integer
   n := new(big.Int)
   if _, ok := n.SetString(s, 10); !ok {
       fmt.Fprintln(os.Stderr, "Invalid integer input")
       os.Exit(1)
   }
   // ProbablyPrime performs a Miller-Rabin test; 20 iterations is sufficient
   if n.ProbablyPrime(20) {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
