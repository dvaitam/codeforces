package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

// This program reads the (collapsed) state of N qubits, represented as
// a binary string of length N, and determines whether the original
// state was the GHZ state (all zeros or all ones) or the W state
// (exactly one qubit was |1>). It outputs 0 for GHZ state and 1 for W state.
func main() {
   scanner := bufio.NewScanner(os.Stdin)
   // Read N
   if !scanner.Scan() {
       return
   }
   n, err := strconv.Atoi(scanner.Text())
   if err != nil {
       fmt.Fprintln(os.Stderr, "invalid N")
       return
   }
   // Read state string
   if !scanner.Scan() {
       return
   }
   s := scanner.Text()
   if len(s) != n {
       fmt.Fprintln(os.Stderr, "state length mismatch")
       return
   }
   // Count number of '1's
   count := 0
   for _, c := range s {
       if c == '1' {
           count++
       }
   }
   // If exactly one '1', it's W state; otherwise GHZ state
   if count == 1 {
       fmt.Println(1)
   } else {
       fmt.Println(0)
   }
}
