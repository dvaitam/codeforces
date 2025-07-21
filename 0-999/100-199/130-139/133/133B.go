package main

import (
   "fmt"
)

func main() {
   const mod = 1000003
   var p string
   if _, err := fmt.Scan(&p); err != nil {
       return
   }
   // mapping of Brainfuck command to 4-bit code
   codes := map[rune]string{
       '>': "1000",
       '<': "1001",
       '+': "1010",
       '-': "1011",
       '.': "1100",
       ',': "1101",
       '[': "1110",
       ']': "1111",
   }
   result := 0
   for _, ch := range p {
       code, ok := codes[ch]
       if !ok {
           continue
       }
       for _, bit := range code {
           result = (result*2 + int(bit-'0')) % mod
       }
   }
   fmt.Println(result)
}
