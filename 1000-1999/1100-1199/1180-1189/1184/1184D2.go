package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k, m int
   if _, err := fmt.Fscan(reader, &n, &k, &m); err != nil {
       return
   }
   // TODO: implement solution
   fmt.Println(0)
}
