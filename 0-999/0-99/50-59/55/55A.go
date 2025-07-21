package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // The flea visits all hassocks iff n is a power of two
   if n > 0 && (n&(n-1)) == 0 {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
