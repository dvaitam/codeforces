package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, k int64
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   k-- // convert to zero-based index
   row := k % n
   col := k / n
   ans := row*m + col + 1
   fmt.Println(ans)
}
