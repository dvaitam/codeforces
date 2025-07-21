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
   ans := 0
   for i := 0; i < n; {
       j := i + 1
       for j < n && s[j] == s[i] {
           j++
       }
       length := j - i
       ans += length / 2
       i = j
   }
   fmt.Println(ans)
}
