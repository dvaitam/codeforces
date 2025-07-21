package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   var s string
   fmt.Fscan(reader, &s)

   count := 0
   for i := 1; i < n; i++ {
       if s[i] == s[i-1] {
           count++
       }
   }
   fmt.Println(count)
}
