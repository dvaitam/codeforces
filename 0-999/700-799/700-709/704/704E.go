package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   // TODO: implement logic converted from solE.cpp
   // Placeholder: output -1 if no collisions
   fmt.Println(-1)
}
