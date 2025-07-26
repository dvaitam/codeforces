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
   var res uint64 = 1
   for i := 2; i <= n; i++ {
       res *= uint64(i)
   }
   fmt.Println(res)
}
