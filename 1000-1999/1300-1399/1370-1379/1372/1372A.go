package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   // Read and discard the first value (e.g., number of test cases)
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // Process remaining inputs
   for {
       if _, err := fmt.Fscan(reader, &n); err != nil {
           break
       }
       // Print '1 ' n times
       for i := 0; i < n; i++ {
           writer.WriteString("1 ")
       }
       // Match original output: a space then newline
       writer.WriteString(" \n")
   }
}
