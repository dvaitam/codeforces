package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

// can returns true if a occurs in s and b occurs in s after non-overlapping a.
func can(s, a, b string) bool {
   i := strings.Index(s, a)
   if i == -1 {
       return false
   }
   // search for b after the end of a
   j := strings.Index(s[i+len(a):], b)
   if j == -1 {
       return false
   }
   return true
}

// reverse returns the reversed string of s.
func reverse(s string) string {
   b := []rune(s)
   for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
       b[i], b[j] = b[j], b[i]
   }
   return string(b)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s, a, b string
   fmt.Fscan(reader, &s, &a, &b)
   forward := can(s, a, b)
   backward := can(reverse(s), a, b)
   switch {
   case forward && backward:
       fmt.Println("both")
   case forward:
       fmt.Println("forward")
   case backward:
       fmt.Println("backward")
   default:
       fmt.Println("fantasy")
   }
}
