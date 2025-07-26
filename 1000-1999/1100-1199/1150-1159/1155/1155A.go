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
   minChar := make([]byte, n)
   pos := make([]int, n)
   minChar[n-1] = s[n-1]
   pos[n-1] = n-1
   for i := n - 2; i >= 0; i-- {
       if s[i+1] < minChar[i+1] {
           minChar[i] = s[i+1]
           pos[i] = i + 1
       } else {
           minChar[i] = minChar[i+1]
           pos[i] = pos[i+1]
       }
   }
   for i := 0; i < n-1; i++ {
       if minChar[i] < s[i] {
           fmt.Println("YES")
           fmt.Println(i+1, pos[i]+1)
           return
       }
   }
   fmt.Println("NO")
}
