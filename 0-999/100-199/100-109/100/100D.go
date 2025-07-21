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
   var s, t string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   ls, lt := len(s), len(t)
   // find common prefix length
   i, max := 0, ls
   if lt < max {
       max = lt
   }
   for i < max && s[i] == t[i] {
       i++
   }
   // minimal operations: remove unmatched suffix of s and add unmatched suffix of t
   ops := (ls - i) + (lt - i)
   if ops <= n {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
