package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   var comps []string
   for i := 0; i < 3; i++ {
       if !scanner.Scan() {
           return
       }
       comps = append(comps, strings.TrimSpace(scanner.Text()))
   }
   perms := []string{"ABC", "ACB", "BAC", "BCA", "CAB", "CBA"}
   for _, perm := range perms {
       ok := true
       for _, comp := range comps {
           if len(comp) < 3 {
               ok = false
               break
           }
           lhs := comp[0]
           op := comp[1]
           rhs := comp[2]
           posL := strings.Index(perm, string(lhs))
           posR := strings.Index(perm, string(rhs))
           if op == '<' {
               if posL >= posR {
                   ok = false
                   break
               }
           } else if op == '>' {
               if posL <= posR {
                   ok = false
                   break
               }
           } else {
               ok = false
               break
           }
       }
       if ok {
           fmt.Println(perm)
           return
       }
   }
   fmt.Println("Impossible")
}
