package main

import (
   "fmt"
   "math"
)

// main prints the amplitudes of the quantum state (|00> + |01> + |10>) / sqrt(3)
// as real numbers in lex order: |00>, |01>, |10>, |11>.
func main() {
   amp := 1.0 / math.Sqrt(3.0)
   // Amplitudes: [amp, amp, amp, 0]
   fmt.Printf("%f %f %f %f\n", amp, amp, amp, 0.0)
}
