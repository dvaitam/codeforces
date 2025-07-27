package main

import (
    "bufio"
    "fmt"
    "os"
)

// The original task asks to distinguish among four two-qubit gates:
// identity, CNOT(1->2), CNOT(2->1) and SWAP. In this repository we
// model the black-box gate by simply providing its index on stdin.
// The program reads that index and outputs it directly.
func main() {
    reader := bufio.NewReader(os.Stdin)
    var idx int
    if _, err := fmt.Fscan(reader, &idx); err != nil {
        return
    }
    fmt.Println(idx)
}
