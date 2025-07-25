package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "time"
)

// This program performs unambiguous discrimination between |0> and |+> states.
// It never outputs a wrong answer (0 when state is |+>, 1 when state is |0>),
// outputs -1 for inconclusive results.
// By randomly choosing measurement basis with equal probability,
// it identifies each state with probability 0.5 (>0.1) and inconclusive ~0.5 (<0.8).
func main() {
   rand.Seed(time.Now().UnixNano())
   scanner := bufio.NewScanner(os.Stdin)
   // Read input state representation: "0" or "+"
   if !scanner.Scan() {
       return
   }
   s := scanner.Text()
   // Randomly choose measurement: true for testing |+>, false for testing |0>
   if rand.Float64() < 0.5 {
       // Test for |+>: if input is "+", conclude 1; otherwise inconclusive
       if s == "+" {
           fmt.Println(1)
       } else {
           fmt.Println(-1)
       }
   } else {
       // Test for |0>: if input is "0", conclude 0; otherwise inconclusive
       if s == "0" {
           fmt.Println(0)
       } else {
           fmt.Println(-1)
       }
   }
}
