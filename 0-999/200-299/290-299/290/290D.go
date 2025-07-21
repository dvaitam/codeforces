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
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   if _, err := fmt.Fscan(reader, &k); err != nil {
       return
   }
   res := make([]rune, len(s))
   for i, c := range s {
       if c >= 'a' && c <= 'z' {
           res[i] = 'a' + (c-'a'+rune(k))%26
       } else if c >= 'A' && c <= 'Z' {
           res[i] = 'A' + (c-'A'+rune(k))%26
       } else {
           res[i] = c
       }
   }
   fmt.Println(string(res))
}
