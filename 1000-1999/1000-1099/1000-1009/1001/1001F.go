package main

import (
   "bufio"
   "fmt"
   "os"
)

// F: Distinguish which of two basis states the qubits are in.
// Classical simulation: read N, two reference bitstrings bits0 and bits1,
// and an unknown state bitstring. Output 0 if it matches bits0, else 1.
func main() {
   reader := bufio.NewReader(os.Stdin)
   var N int
   if _, err := fmt.Fscan(reader, &N); err != nil {
       return
   }
   var b0, b1, s string
   if _, err := fmt.Fscan(reader, &b0); err != nil {
       return
   }
   if _, err := fmt.Fscan(reader, &b1); err != nil {
       return
   }
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   // Compare the unknown state to bits0
   if s == b0 {
       fmt.Println(0)
   } else {
       fmt.Println(1)
   }
}
