package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var d int
   var s string
   if _, err := fmt.Fscan(in, &d); err != nil {
       return
   }
   fmt.Fscan(in, &s)
   n := len(s)
   // No palindrome substring of length >= d allowed
   if d == 1 {
       fmt.Fprintln(out, "Impossible")
       return
   }
   orig := []byte(s)
   t := make([]byte, n)
   copy(t, orig)

   // check if placing t[pos] is valid
   ok := func(pos int) bool {
       if pos - d + 1 >= 0 && t[pos] == t[pos - d + 1] {
           return false
       }
       if pos - d >= 0 && t[pos] == t[pos - d] {
           return false
       }
       return true
   }

   // try to find position to increment
   for i := n - 1; i >= 0; i-- {
       // try next characters
       for c := orig[i] + 1; c <= 'z'; c++ {
           t[i] = c
           if !ok(i) {
               continue
           }
           // fill suffix minimally
           viable := true
           for j := i + 1; j < n; j++ {
               placed := false
               for c2 := byte('a'); c2 <= 'z'; c2++ {
                   t[j] = c2
                   if ok(j) {
                       placed = true
                       break
                   }
               }
               if !placed {
                   viable = false
                   break
               }
           }
           if viable {
               fmt.Fprintln(out, string(t))
               return
           }
       }
       // restore t[i] to original before moving on (prefix untouched beyond i)
       t[i] = orig[i]
   }
   fmt.Fprintln(out, "Impossible")
}
