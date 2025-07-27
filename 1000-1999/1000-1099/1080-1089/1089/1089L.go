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
   vowels := map[rune]bool{
       'a': true,
       'e': true,
       'i': true,
       'o': true,
       'u': true,
   }
   var result []rune
   for _, ch := range s {
       if !vowels[ch] {
           result = append(result, ch)
       }
   }
   fmt.Println(string(result))
}
