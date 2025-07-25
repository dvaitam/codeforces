package main

import (
   "bufio"
   "fmt"
   "os"
)

// This program prints a checkerboard pattern of size 2^N x 2^N
// where each block is 2x2 cells, alternating X (non-zero) and . (zero),
// with the top-left 2x2 block filled with X.
func main() {
   reader := bufio.NewReader(os.Stdin)
   var N int
   if _, err := fmt.Fscan(reader, &N); err != nil {
       return
   }
   size := 1 << N
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < size; i++ {
       for j := 0; j < size; j++ {
           // Determine block indices by integer division by 2
           if ((i/2)+(j/2))%2 == 0 {
               writer.WriteByte('X')
           } else {
               writer.WriteByte('.')
           }
       }
       writer.WriteByte('\n')
   }
}
