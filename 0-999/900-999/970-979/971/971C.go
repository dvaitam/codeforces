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
   r := []rune(s)
   for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
       r[i], r[j] = r[j], r[i]
   }
   fmt.Println(string(r))
}
