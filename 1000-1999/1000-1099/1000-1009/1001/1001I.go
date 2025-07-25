package main

import (
   "bufio"
   "fmt"
   "os"
)

// A classical implementation simulating the Deutsch-Jozsa oracle test.
// Reads an integer N and 2^N function values (0 or 1) from stdin
// Outputs "true" if the function is constant (all values equal), else "false".
func main() {
   reader := bufio.NewReader(os.Stdin)
   var N int
   if _, err := fmt.Fscan(reader, &N); err != nil {
       return
   }
   total := 1 << N
   count := 0
   for i := 0; i < total; i++ {
       var v int
       if _, err := fmt.Fscan(reader, &v); err != nil {
           return
       }
       count += v
   }
   if count == 0 || count == total {
       fmt.Println("true")
   } else {
       fmt.Println("false")
   }
}
