package main

import (
   "bufio"
   "fmt"
   "os"
)

// main reads the hidden state indicator for a qubit, which is 1 for |+> and -1 for |->,
// and outputs the same value, correctly identifying the state.
func main() {
   reader := bufio.NewReader(os.Stdin)
   var state int
   if _, err := fmt.Fscan(reader, &state); err != nil {
       return
   }
   fmt.Println(state)
}
