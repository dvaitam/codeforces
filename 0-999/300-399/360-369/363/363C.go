package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   input, _ := reader.ReadString('\n')
   s := strings.TrimSpace(input)
   n := len(s)
   ans := make([]byte, 0, n)
   for i := 0; i < n; i++ {
       c := s[i]
       m := len(ans)
       if m >= 2 && ans[m-1] == c && ans[m-2] == c {
           continue
       }
       if m >= 3 && ans[m-3] == ans[m-2] && ans[m-1] == c {
           continue
       }
       ans = append(ans, c)
   }
   fmt.Println(string(ans))
}
