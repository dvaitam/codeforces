package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // Read guesses
   guesses := make([]string, n)
   bulls := make([]int, n)
   cows := make([]int, n)
   for i := 0; i < n; i++ {
       var s string
       if _, err := fmt.Fscan(reader, &s, &bulls[i], &cows[i]); err != nil {
           return
       }
       // ensure 4 chars (leading zeros handled as input strings)
       guesses[i] = s
   }

   var candidates []string
   // generate all 4-digit numbers with distinct digits
   for d1 := '0'; d1 <= '9'; d1++ {
       for d2 := '0'; d2 <= '9'; d2++ {
           if d2 == d1 {
               continue
           }
           for d3 := '0'; d3 <= '9'; d3++ {
               if d3 == d1 || d3 == d2 {
                   continue
               }
               for d4 := '0'; d4 <= '9'; d4++ {
                   if d4 == d1 || d4 == d2 || d4 == d3 {
                       continue
                   }
                   cand := string([]rune{d1, d2, d3, d4})
                   ok := true
                   // check against all guesses
                   for i := 0; i < n; i++ {
                       g := guesses[i]
                       // compute bulls
                       bcnt := 0
                       for k := 0; k < 4; k++ {
                           if cand[k] == g[k] {
                               bcnt++
                           }
                       }
                       if bcnt != bulls[i] {
                           ok = false
                           break
                       }
                       // compute total matches
                       match := 0
                       // use array of 10 bools for cand
                       var seen [10]bool
                       for k := 0; k < 4; k++ {
                           seen[cand[k]-'0'] = true
                       }
                       for k := 0; k < 4; k++ {
                           if seen[g[k]-'0'] {
                               match++
                           }
                       }
                       ccnt := match - bcnt
                       if ccnt != cows[i] {
                           ok = false
                           break
                       }
                   }
                   if ok {
                       candidates = append(candidates, cand)
                       // early exit if too many
                       if len(candidates) > 1 {
                           // still need to know if 0, >1 or =1, so we can break later
                       }
                   }
               }
           }
       }
   }

   switch len(candidates) {
   case 0:
       fmt.Println("Incorrect data")
   case 1:
       fmt.Println(candidates[0])
   default:
       fmt.Println("Need more data")
   }
}
