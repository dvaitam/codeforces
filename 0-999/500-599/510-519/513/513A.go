package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n1, n2, k1, k2 int
   if _, err := fmt.Fscan(reader, &n1, &n2, &k1, &k2); err != nil {
       return
   }
   if n1 > n2 {
       fmt.Println("First")
   } else {
       fmt.Println("Second")
   }
}
