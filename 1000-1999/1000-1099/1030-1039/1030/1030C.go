package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var s string
   if _, err := fmt.Fscan(reader, &n, &s); err != nil {
       return
   }
   // convert digits to ints
   a := make([]int, n)
   for i := 0; i < n; i++ {
       a[i] = int(s[i] - '0')
   }
   // try prefix cut at i: first segment [0..i] sum
   for i := 1; i < n; i++ {
       target := 0
       for j := 0; j < i; j++ {
           target += a[j]
       }
       // greedily check rest segments
       curr := 0
       count := 0
       ok := true
       for j := i; j < n; j++ {
           curr += a[j]
           if curr > target {
               ok = false
               break
           }
           if curr == target {
               curr = 0
               count++
           }
       }
       if ok && curr == 0 && count >= 1 {
           fmt.Println("YES")
           return
       }
   }
   fmt.Println("NO")
}
