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
       // no input
       return
   }
   s := strings.TrimSpace(line)
   n := len(s)
   // find first '['
   st1 := -1
   for i := 0; i < n; i++ {
       if s[i] == '[' {
           st1 = i
           break
       }
   }
   if st1 < 0 {
       fmt.Println(-1)
       return
   }
   // find first ':' after st1
   st2 := -1
   for i := st1 + 1; i < n; i++ {
       if s[i] == ':' {
           st2 = i
           break
       }
   }
   if st2 < 0 {
       fmt.Println(-1)
       return
   }
   // find last ']' and last ':' before it
   en1 := -1
   for i := n - 1; i >= 0; i-- {
       if s[i] == ']' {
           en1 = i
           break
       }
   }
   if en1 < 0 || en1 <= st2 {
       fmt.Println(-1)
       return
   }
   en2 := -1
   for i := en1 - 1; i >= 0; i-- {
       if s[i] == ':' {
           en2 = i
           break
       }
   }
   if en2 < 0 || en2 <= st2 {
       fmt.Println(-1)
       return
   }
   // count '|' between st2 and en2
   cnt := 4
   for i := st2 + 1; i < en2; i++ {
       if s[i] == '|' {
           cnt++
       }
   }
   fmt.Println(cnt)
}
