package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var A string
   if _, err := fmt.Fscan(reader, &A); err != nil {
       return
   }
   n := len(A)
   // collect unique substrings
   subs := make(map[string]struct{})
   for i := 0; i < n; i++ {
       for j := i + 1; j <= n; j++ {
           subs[A[i:j]] = struct{}{}
       }
   }
   bestCount := 0
   bestLen := 0
   bestS := ""
   // evaluate each substring
   for s := range subs {
       m := len(s)
       // count occurrences (including overlapping)
       cnt := 0
       for i := 0; i <= n-m; i++ {
           if A[i:i+m] == s {
               cnt++
           }
       }
       if cnt > bestCount || (cnt == bestCount && (m > bestLen || (m == bestLen && s > bestS))) {
           bestCount = cnt
           bestLen = m
           bestS = s
       }
   }
   // output answer
   fmt.Fprintln(os.Stdout, bestS)
}
