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
   // Check length requirement
   if len(s) < 5 {
       fmt.Println("Too weak")
       return
   }
   // Flags for required character types
   hasUpper, hasLower, hasDigit := false, false, false
   for _, c := range s {
       if c >= 'A' && c <= 'Z' {
           hasUpper = true
       } else if c >= 'a' && c <= 'z' {
           hasLower = true
       } else if c >= '0' && c <= '9' {
           hasDigit = true
       }
   }
   // Final check
   if hasUpper && hasLower && hasDigit {
       fmt.Println("Correct")
   } else {
       fmt.Println("Too weak")
   }
}
