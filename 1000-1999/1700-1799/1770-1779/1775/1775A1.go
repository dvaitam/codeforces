package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var T int
   fmt.Fscan(reader, &T)
   for t := 0; t < T; t++ {
       var s string
       fmt.Fscan(reader, &s)
       solve(s)
   }
}

// solve processes one test case
func solve(s string) {
   n := len(s)
   a := s[0:1]
   c := s[n-1 : n]
   l, r := 1, n-1
   for l <= r && r < n {
       b := s[l:r]
       if (b >= a && b >= c) || (b <= a && b <= c) {
           fmt.Println(a, b, c)
           return
       } else if b >= a && b <= c {
           if len(b) <= 1 {
               fmt.Println(":(")
               return
           }
           l++
       } else {
           // b <= a && b >= c
           if len(b) <= 1 {
               fmt.Println(":(")
               return
           }
           r--
       }
   }
   fmt.Println(":(")
}
