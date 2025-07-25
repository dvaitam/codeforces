package main

import (
   "bufio"
   "fmt"
   "os"
)

// E: Distinguish which Bell state the two qubits are in.
// Classical simulation: read the Bell-state index (0 to 3) and output it.
func main() {
   reader := bufio.NewReader(os.Stdin)
   var idx int
   if _, err := fmt.Fscan(reader, &idx); err != nil {
       return
   }
   fmt.Println(idx)
}
