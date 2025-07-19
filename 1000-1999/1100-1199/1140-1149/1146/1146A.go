package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   if !scanner.Scan() {
       return
   }
   s := scanner.Text()
   cnt := 0
   for i := 0; i < len(s); i++ {
       if s[i] == 'a' {
           cnt++
       }
   }
   length := len(s)
   maxKeep := 2*cnt - 1
   if maxKeep >= length {
       fmt.Println(length)
   } else {
       fmt.Println(maxKeep)
   }
}
