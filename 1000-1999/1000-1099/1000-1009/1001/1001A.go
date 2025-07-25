package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

// Simulate preparing |+> or |-> state from |0> based on sign.
// Reads integer sign (1 or -1) and outputs the amplitudes of |0> and |1>.
func main() {
   reader := bufio.NewReader(os.Stdin)
   var sign int
   if _, err := fmt.Fscan(reader, &sign); err != nil {
       return
   }
   // Hadamard maps |0> to (|0> + |1>)/√2; apply Z if sign == -1
   inv := 1 / math.Sqrt2
   a := inv
   b := inv * float64(sign)
   // output amplitudes for |0> and |1>
   fmt.Printf("%.10f %.10f", a, b)
}
