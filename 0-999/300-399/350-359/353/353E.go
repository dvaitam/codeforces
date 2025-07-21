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
   n := len(s)
   hasIncoming := make([]bool, n)
   for i := 0; i < n; i++ {
       if s[i] == '0' {
           j := i + 1
           if j == n {
               j = 0
           }
           hasIncoming[j] = true
       } else {
           hasIncoming[i] = true
       }
   }
   count := 0
   for i := 0; i < n; i++ {
       if !hasIncoming[i] {
           count++
       }
   }
   fmt.Println(count)
}
