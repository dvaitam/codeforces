package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   res := 0
   for i := 0; i < n; {
       j := i + 1
       for j < n && s[j] == s[i] {
           j++
       }
       if (j-i)%2 == 0 {
           res++
       }
       i = j
   }
   fmt.Println(res)
}
