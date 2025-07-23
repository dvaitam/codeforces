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
   if err != nil && err.Error() != "EOF" {
       // unexpected error
       return
   }
   s = strings.TrimSpace(s)

   // check AB then BA
   if containsNonOverlap(s, "AB", "BA") || containsNonOverlap(s, "BA", "AB") {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}

// containsNonOverlap checks if pattern1 occurs and then pattern2 occurs non-overlapping
func containsNonOverlap(s, p1, p2 string) bool {
   n := len(s)
   // find first p1
   for i := 0; i+len(p1) <= n; i++ {
       if s[i:i+len(p1)] == p1 {
           // search for p2 starting from i+len(p1)
           for j := i + len(p1); j+len(p2) <= n; j++ {
               if s[j:j+len(p2)] == p2 {
                   return true
               }
           }
           break
       }
   }
   return false
}
