package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   if n == 1 {
       fmt.Println("a a")
       return
   }
   a := make([]byte, 1, 100)
   for c := byte('a'); c <= byte('z'); c++ {
       a = append(a, c)
   }
   for c := byte('A'); c <= byte('Z'); c++ {
       a = append(a, c)
   }
   m := 0
   var s []byte
   for n > 1 {
       n--
       if n%2 == 0 {
           m++
           s = append(s, a[m], a[m])
           n /= 2
       } else {
           m++
           s = append(s, a[m])
       }
   }
   t := make([]byte, 0)
   for i := 1; i <= m; i++ {
       s = append(s, a[i])
       t = append(t, a[i])
   }
   fmt.Printf("%s %s", s, t)
}
