package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   sum := 0
   for _, ch := range s {
       sum += int(ch - '0')
   }
   fmt.Println(sum)
}
