package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   var k int
   if _, err := fmt.Fscan(reader, &s, &k); err != nil {
       return
   }
   N := len(s)
   total := N + k
   bs := []byte(s)
   maxLen := 0
   for i := 0; i < total; i++ {
       // try half-length n
       for n := 1; i+2*n <= total; n++ {
           ok := true
           for j := 0; j < n; j++ {
               p1 := i + j
               p2 := i + n + j
               if p1 < N && p2 < N && bs[p1] != bs[p2] {
                   ok = false
                   break
               }
           }
           if ok {
               l := 2 * n
               if l > maxLen {
                   maxLen = l
               }
           }
       }
   }
   fmt.Println(maxLen)
}
