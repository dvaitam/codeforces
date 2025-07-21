package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   if scanner.Scan() {
       s := scanner.Text()
       for _, c := range s {
           if c == 'H' || c == 'Q' || c == '9' {
               fmt.Println("YES")
               return
           }
       }
   }
   fmt.Println("NO")
}
