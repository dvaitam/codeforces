package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   line, err := reader.ReadString('\n')
   if err != nil && len(line) == 0 {
       fmt.Println("NO")
       return
   }
   s := strings.TrimRight(line, "\r\n")
   if !isValidJabberID(s) {
       fmt.Println("NO")
   } else {
       fmt.Println("YES")
   }
}

func isValidJabberID(s string) bool {
   // Split username and rest
   if strings.Count(s, "@") != 1 {
       return false
   }
   parts := strings.SplitN(s, "@", 2)
   user := parts[0]
   rest := parts[1]
   if !isValidToken(user, 1, 16) {
       return false
   }
   // Resource optional
   if strings.Count(rest, "/") > 1 {
       return false
   }
   var host, res string
   if idx := strings.Index(rest, "/"); idx >= 0 {
       host = rest[:idx]
       res = rest[idx+1:]
       if !isValidToken(res, 1, 16) {
           return false
       }
   } else {
       host = rest
   }
   // Validate host
   if len(host) < 1 || len(host) > 32 {
       return false
   }
   // Split host by dots
   labels := strings.Split(host, ".")
   for _, lbl := range labels {
       if !isValidToken(lbl, 1, 16) {
           return false
       }
   }
   return true
}

func isValidToken(s string, minLen, maxLen int) bool {
   n := len(s)
   if n < minLen || n > maxLen {
       return false
   }
   for i := 0; i < n; i++ {
       c := s[i]
       if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_' {
           continue
       }
       return false
   }
   return true
}
