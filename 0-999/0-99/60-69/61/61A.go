package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a, b string
   fmt.Fscan(reader, &a)
   fmt.Fscan(reader, &b)
   n := len(a)
   res := make([]byte, n)
   for i := 0; i < n; i++ {
       if a[i] != b[i] {
           res[i] = '1'
       } else {
           res[i] = '0'
       }
   }
   fmt.Println(string(res))
}
