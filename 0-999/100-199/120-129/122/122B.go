package main

import "fmt"

func main() {
   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   n := len(s)
   // collect all lucky substrings
   candidates := make(map[string]struct{})
   for i := 0; i < n; i++ {
       if s[i] != '4' && s[i] != '7' {
           continue
       }
       for j := i; j < n && (s[j] == '4' || s[j] == '7'); j++ {
           substr := s[i : j+1]
           candidates[substr] = struct{}{}
       }
   }
   best := ""
   bestCount := 0
   // for each candidate, count occurrences
   for substr := range candidates {
       cnt := 0
       m := len(substr)
       for k := 0; k+m <= n; k++ {
           if s[k:k+m] == substr {
               cnt++
           }
       }
       if cnt > bestCount || (cnt == bestCount && (best == "" || substr < best)) {
           best = substr
           bestCount = cnt
       }
   }
   if bestCount == 0 {
       fmt.Println("-1")
   } else {
       fmt.Println(best)
   }
}
