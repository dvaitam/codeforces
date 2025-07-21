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
   a, b := 1, 1
   for i := 2; i <= n; i++ {
       a, b = b, a+b
   }
   fmt.Fprintln(os.Stdout, b)
}
