package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   s, err := reader.ReadString('\n')
   if err != nil && len(s) == 0 {
       return
   }
   // remove trailing newline and carriage return
   end := len(s)
   if end > 0 && s[end-1] == '\n' {
       end--
   }
   if end > 0 && s[end-1] == '\r' {
       end--
   }
   str := s[:end]
   n := len(str)
   var result uint64
   var carry uint8
   for i := n - 1; i >= 1; i-- {
       bit := uint8(str[i] - '0')
       sum := bit + carry
       if sum == 0 {
           result++
       } else if sum == 1 {
           result += 2
           carry = 1
       } else { // sum == 2
           result++
           carry = 1
       }
   }
   // if there's a carry at the highest bit, one more division
   if carry == 1 {
       result++
   }
   fmt.Println(result)
}
