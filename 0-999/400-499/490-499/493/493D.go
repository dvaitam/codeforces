package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int64
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   if n == 2 {
       // White can capture black queen immediately
       fmt.Println("white")
       fmt.Println(1, 2)
   } else {
       fmt.Println("black")
   }
}
