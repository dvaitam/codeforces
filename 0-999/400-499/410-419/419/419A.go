package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   _, err := fmt.Fscan(reader, &s)
   if err != nil {
       return
   }
   // Set of mirror-symmetric uppercase letters
   sym := map[rune]bool{
       'A': true, 'H': true, 'I': true, 'M': true,
       'O': true, 'T': true, 'U': true, 'V': true,
       'W': true, 'X': true, 'Y': true,
   }
   n := len(s)
   ok := true
   for i, r := range s {
       // Check symmetric character
       if !sym[r] {
           ok = false
           break
       }
       // Check palindrome property
       if r != rune(s[n-1-i]) {
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
