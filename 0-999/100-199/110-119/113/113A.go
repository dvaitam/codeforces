package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   scanner.Split(bufio.ScanWords)
   var words []string
   for scanner.Scan() {
       words = append(words, scanner.Text())
   }
   if len(words) == 0 {
       return
   }
   // suffix, type: 0=adjective,1=noun,2=verb, gender: 0=male,1=female
   suffixes := []struct{ suf string; typ, gender int }{
       {"lios", 0, 0},
       {"liala", 0, 1},
       {"etr", 1, 0},
       {"etra", 1, 1},
       {"initis", 2, 0},
       {"inites", 2, 1},
   }
   types := make([]int, len(words))
   var gender *int
   nounCount := 0
   for i, w := range words {
       matched := false
       for _, s := range suffixes {
           if strings.HasSuffix(w, s.suf) {
               matched = true
               types[i] = s.typ
               if gender == nil {
                   g := s.gender
                   gender = &g
               } else if *gender != s.gender {
                   fmt.Println("NO")
                   return
               }
               if s.typ == 1 {
                   nounCount++
               }
               break
           }
       }
       if !matched {
           fmt.Println("NO")
           return
       }
   }
   if len(words) == 1 {
       fmt.Println("YES")
       return
   }
   if nounCount != 1 {
       fmt.Println("NO")
       return
   }
   // check order: adjectives (0)*, noun (1), verbs (2)*
   state := 0
   for _, t := range types {
       switch state {
       case 0:
           if t == 0 {
               continue
           } else if t == 1 {
               state = 1
           } else {
               fmt.Println("NO")
               return
           }
       case 1:
           if t == 1 {
               continue
           } else if t == 2 {
               state = 2
           } else {
               fmt.Println("NO")
               return
           }
       case 2:
           if t == 2 {
               continue
           }
           fmt.Println("NO")
           return
       }
   }
   fmt.Println("YES")
}
