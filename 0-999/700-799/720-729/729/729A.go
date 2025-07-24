package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var s string
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   var result []rune
   for i := 0; i < len(s); {
       // check for "ogo"
       if i+3 <= len(s) && s[i:i+3] == "ogo" {
           j := i + 3
           // consume "go" repeats
           for j+1 < len(s) && s[j:j+2] == "go" {
               j += 2
           }
           // replace with ***
           result = append(result, '*', '*', '*')
           i = j
       } else {
           result = append(result, rune(s[i]))
           i++
       }
   }
   fmt.Println(string(result))
}
