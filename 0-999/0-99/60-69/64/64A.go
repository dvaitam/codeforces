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
   result := 1
   for i := 1; i <= n; i++ {
       result *= i
   }
   fmt.Println(result)
}
