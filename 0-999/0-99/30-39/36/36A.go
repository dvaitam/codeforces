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
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   // Collect indices of received signals
   pos := make([]int, 0, len(s))
   for i, ch := range s {
       if ch == '1' {
           pos = append(pos, i)
       }
   }
   // There are guaranteed at least three '1's
   // Check if intervals between successive signals are equal
   d := pos[1] - pos[0]
   ok := true
   for i := 2; i < len(pos); i++ {
       if pos[i]-pos[i-1] != d {
           ok = false
           break
       }
   }
   if ok {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
