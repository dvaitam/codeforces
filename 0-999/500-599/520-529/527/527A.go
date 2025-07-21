package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a, b int64
   if _, err := fmt.Fscan(reader, &a, &b); err != nil {
       return
   }
   var ships int64
   for b != 0 {
       ships += a / b
       a, b = b, a % b
   }
   fmt.Println(ships)
}
