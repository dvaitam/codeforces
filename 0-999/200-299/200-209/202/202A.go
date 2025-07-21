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
   n := len(s)
   best := ""
   // Enumerate all non-empty subsequences
   for mask := 1; mask < (1 << n); mask++ {
       t := ""
       for i := 0; i < n; i++ {
           if mask&(1<<i) != 0 {
               t += string(s[i])
           }
       }
       // Check if palindrome
       ok := true
       for i, j := 0, len(t)-1; i < j; i, j = i+1, j-1 {
           if t[i] != t[j] {
               ok = false
               break
           }
       }
       if ok && t > best {
           best = t
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, best)
}
