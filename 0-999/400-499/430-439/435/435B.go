package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a string
   var k int
   _, _ = fmt.Fscan(reader, &a, &k)
   s := []byte(a)
   n := len(s)
   for i := 0; i < n && k > 0; i++ {
       maxDigit := s[i]
       pos := i
       limit := i + k
       if limit >= n {
           limit = n - 1
       }
       for j := i + 1; j <= limit; j++ {
           if s[j] > maxDigit {
               maxDigit = s[j]
               pos = j
           }
       }
       if pos != i {
           for j := pos; j > i; j-- {
               s[j], s[j-1] = s[j-1], s[j]
           }
           k -= pos - i
       }
   }
   fmt.Println(string(s))
}
