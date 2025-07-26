package main

import (
   "fmt"
)

// main prints an 8x8 identity unitary matrix for N=3 qubits.
// TODO: update to match the specific pattern for problem D5.
func main() {
   const N = 3
   dim := 1 << N
   for i := 0; i < dim; i++ {
       for j := 0; j < dim; j++ {
           var x float64
           if i == j {
               x = 1.0
           } else {
               x = 0.0
           }
           if j > 0 {
               fmt.Print(" ")
           }
           fmt.Printf("%f", x)
       }
       fmt.Println()
   }
}
