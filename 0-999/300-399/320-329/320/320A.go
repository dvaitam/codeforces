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
   ok := true
   for i := 0; i < len(s); {
       if i+3 <= len(s) && s[i:i+3] == "144" {
           i += 3
       } else if i+2 <= len(s) && s[i:i+2] == "14" {
           i += 2
       } else if s[i] == '1' {
           i++
       } else {
           ok = false
           break
       }
   }
   if ok {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
