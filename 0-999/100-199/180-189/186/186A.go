package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   s, err := reader.ReadString('\n')
   if err != nil && len(s) == 0 {
       return
   }
   t, err := reader.ReadString('\n')
   if err != nil && len(t) == 0 {
       return
   }
   s = strings.TrimSpace(s)
   t = strings.TrimSpace(t)
   if len(s) != len(t) {
       fmt.Println("NO")
       return
   }
   var diffs []int
   for i := 0; i < len(s); i++ {
       if s[i] != t[i] {
           diffs = append(diffs, i)
           if len(diffs) > 2 {
               fmt.Println("NO")
               return
           }
       }
   }
   if len(diffs) != 2 {
       fmt.Println("NO")
       return
   }
   i, j := diffs[0], diffs[1]
   if s[i] == t[j] && s[j] == t[i] {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
