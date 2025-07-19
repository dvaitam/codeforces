package main

import (
   "bufio"
   "fmt"
   "os"
)

// ispald checks whether b is a palindrome.
// Returns -1 if not a palindrome,
// a non-negative index x if palindrome and b[x]!=b[x-1] for some x>0 (last such x),
// or -2 if all characters are identical.
func ispald(b []byte) int {
   res := -2
   n := len(b)
   for i := 0; i <= n/2; i++ {
       if b[i] != b[n-i-1] {
           return -1
       }
       if i > 0 && b[i] != b[i-1] {
           res = i
       }
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var s string
       fmt.Fscan(reader, &s)
       b := []byte(s)
       x := ispald(b)
       if x == -1 {
           writer.WriteString(s)
           writer.WriteByte('\n')
       } else if x != -2 {
           b[x], b[x-1] = b[x-1], b[x]
           writer.Write(b)
           writer.WriteByte('\n')
       } else {
           writer.WriteString("-1\n")
       }
   }
}
