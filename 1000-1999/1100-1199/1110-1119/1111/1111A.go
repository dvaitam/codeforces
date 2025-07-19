package main

import (
   "bufio"
   "fmt"
   "os"
)

// isVowel returns true if c is a vowel (a, u, e, i, o)
func isVowel(c byte) bool {
   switch c {
   case 'a', 'u', 'e', 'i', 'o':
       return true
   default:
       return false
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var s, t string
   // read two strings
   if _, err := fmt.Fscan(reader, &s, &t); err != nil {
       return
   }
   // if lengths differ, answer is No
   if len(s) != len(t) {
       fmt.Fprint(writer, "No")
       return
   }
   // check each character
   for i := 0; i < len(s); i++ {
       if isVowel(s[i]) != isVowel(t[i]) {
           fmt.Fprint(writer, "No")
           return
       }
   }
   fmt.Fprint(writer, "Yes")
}
