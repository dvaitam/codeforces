package main

import (
   "bufio"
   "fmt"
   "io"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   s, err := reader.ReadString('\n')
   if err != nil && err != io.EOF {
       return
   }
   s = strings.TrimRight(s, "\r\n")
   open := 0
   ans := 0
   for i := 0; i < len(s); i++ {
       switch s[i] {
       case '(':
           open++
       case ')':
           if open > 0 {
               open--
               ans += 2
           }
       }
   }
   fmt.Println(ans)
}
