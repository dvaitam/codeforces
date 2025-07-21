package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a, s string
   if _, err := fmt.Fscan(reader, &a, &s); err != nil {
      return
   }
   sb := []byte(s)
   sort.Slice(sb, func(i, j int) bool { return sb[i] > sb[j] })
   ab := []byte(a)
   j := 0
   m := len(sb)
   for i := 0; i < len(ab) && j < m; i++ {
      if sb[j] > ab[i] {
         ab[i] = sb[j]
         j++
      }
   }
   fmt.Println(string(ab))
}
