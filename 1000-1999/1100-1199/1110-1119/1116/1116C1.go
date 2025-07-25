package main

import (
   "bufio"
   "fmt"
   "os"
)

// main reads an integer N and a binary string of length N, and outputs 1
// if the bits in the string alternate (i.e., no two adjacent bits are equal),
// otherwise outputs 0.
func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   ok := true
   for i := 1; i < n; i++ {
       if s[i] == s[i-1] {
           ok = false
           break
       }
   }
   if ok {
       fmt.Println(1)
   } else {
       fmt.Println(0)
   }
}
