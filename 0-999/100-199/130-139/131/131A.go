package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   fmt.Fscan(reader, &s)
   if isCapsLockError(s) {
       // invert case
       r := []rune(s)
       for i, c := range r {
           if c >= 'a' && c <= 'z' {
               r[i] = c - ('a' - 'A')
           } else if c >= 'A' && c <= 'Z' {
               r[i] = c + ('a' - 'A')
           }
       }
       s = string(r)
   }
   fmt.Println(s)
}

// isCapsLockError returns true if either all letters are uppercase,
// or all letters except the first are uppercase.
func isCapsLockError(s string) bool {
   n := len(s)
   if n == 0 {
       return false
   }
   allUpper := true
   for _, c := range s {
       if c < 'A' || c > 'Z' {
           allUpper = false
           break
       }
   }
   if allUpper {
       return true
   }
   // check all except first
   for i, c := range s {
       if i == 0 {
           continue
       }
       if c < 'A' || c > 'Z' {
           return false
       }
   }
   return true
}
