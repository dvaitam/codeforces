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
   // dp[i][usedAt] = best string to reach position i with usedAt (0/1)
   type state struct { ok bool; str string }
   dp := make([][2]state, n+1)
   dp[0][0] = state{ok: true, str: ""}
   // helper to update state
   update := func(i, used int, cand string) {
       st := &dp[i][used]
       if !st.ok || len(cand) < len(st.str) || (len(cand) == len(st.str) && cand < st.str) {
           st.ok = true
           st.str = cand
       }
   }
   for i := 0; i < n; i++ {
       for used := 0; used < 2; used++ {
           cur := dp[i][used]
           if !cur.ok {
               continue
           }
           // take letter
           c := s[i]
           // always a lowercase letter
           update(i+1, used, cur.str+string(c))
           // replace "at" with '@'
           if used == 0 && i > 0 && i+2 < n {
               if s[i] == 'a' && s[i+1] == 't' {
                   update(i+2, 1, cur.str+"@")
               }
           }
           // replace "dot" with '.'
           if i > 0 && i+3 < n {
               if s[i] == 'd' && s[i+1] == 'o' && s[i+2] == 't' {
                   update(i+3, used, cur.str+".")
               }
           }
       }
   }
   if dp[n][1].ok {
       fmt.Print(dp[n][1].str)
   }
}
