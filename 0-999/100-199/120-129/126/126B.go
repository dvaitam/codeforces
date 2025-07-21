package main

import (
   "bufio"
   "fmt"
   "io"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   s, err := reader.ReadString('\n')
   if err != nil && err != io.EOF {
       return
   }
   // remove trailing newline/carriage return
   s = strings.TrimRight(s, "\r\n")
   n := len(s)
   if n == 0 {
       fmt.Println("Just a legend")
       return
   }
   // compute prefix function (KMP)
   pi := make([]int, n)
   for i := 1; i < n; i++ {
       j := pi[i-1]
       for j > 0 && s[i] != s[j] {
           j = pi[j-1]
       }
       if s[i] == s[j] {
           j++
       }
       pi[i] = j
   }
   // find longest border that appears in the middle
   best := 0
   // maximum pi value excluding last position
   for i := 0; i < n-1; i++ {
       if pi[i] > best {
           best = pi[i]
       }
   }
   k := pi[n-1]
   for k > 0 {
       if best >= k {
           fmt.Println(s[:k])
           return
       }
       // fallback to next longest border
       k = pi[k-1]
   }
   fmt.Println("Just a legend")
}
