package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var s, t string
   fmt.Fscan(reader, &s)
   fmt.Fscan(reader, &t)
   n := len(s)
   // counts
   Cs := make([]int, 26)
   Ct := make([]int, 26)
   for i := 0; i < n; i++ {
       Cs[s[i]-'A']++
       Ct[t[i]-'A']++
   }
   // maximum matches
   bestMatches := 0
   for c := 0; c < 26; c++ {
       bestMatches += min(Cs[c], Ct[c])
   }
   // remaining counts
   remCt := make([]int, 26)
   remCs := make([]int, 26)
   copy(remCt, Ct)
   copy(remCs, Cs)
   matched := 0
   u := make([]byte, n)
   // build u
   for i := 0; i < n; i++ {
       si := int(s[i] - 'A')
       // try each candidate letter
       for c := 0; c < 26; c++ {
           if remCt[c] == 0 {
               continue
           }
           // simulate assign c at pos i
           // matches if c == si
           inc := 0
           if c == si {
               inc = 1
           }
           // compute potential matches in remaining positions
           // remCt2 and remCs2
           sum := 0
           for x := 0; x < 26; x++ {
               rc := remCt[x]
               rcs := remCs[x]
               if x == c {
                   rc--
               }
               if x == si {
                   rcs--
               }
               if rc < 0 {
                   rc = 0
               }
               if rcs < 0 {
                   rcs = 0
               }
               sum += min(rc, rcs)
           }
           if matched+inc+sum >= bestMatches {
               // choose c
               u[i] = byte('A' + c)
               if c == si {
                   matched++
               }
               remCt[c]--
               remCs[si]--
               break
           }
       }
   }
   // minimal replacements
   z := n - bestMatches
   fmt.Fprintln(writer, z)
   fmt.Fprintln(writer, string(u))
}
