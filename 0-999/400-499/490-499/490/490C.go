package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   var a, b int64
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   fmt.Fscan(reader, &a, &b)
   n := len(s)
   prefix := make([]int64, n)
   for i := 0; i < n; i++ {
       d := int64(s[i] - '0')
       if i == 0 {
           prefix[i] = d % a
       } else {
           prefix[i] = (prefix[i-1]*10 + d) % a
       }
   }
   suffix := make([]int64, n+1)
   var mult int64 = 1
   for i := n - 1; i >= 0; i-- {
       d := int64(s[i] - '0')
       suffix[i] = (d*mult + suffix[i+1]) % b
       mult = (mult * 10) % b
   }
   for i := 1; i < n; i++ {
       if prefix[i-1] == 0 && suffix[i] == 0 && s[i] != '0' {
           fmt.Println("YES")
           fmt.Println(s[:i])
           fmt.Println(s[i:])
           return
       }
   }
   fmt.Println("NO")
}
