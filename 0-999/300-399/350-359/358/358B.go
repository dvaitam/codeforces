package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   words := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &words[i])
   }
   var s string
   fmt.Fscan(reader, &s)

   // Build the required pattern: <3word1<3word2...<3wordn<3
   var patternBuilder strings.Builder
   for _, w := range words {
       patternBuilder.WriteString("<3")
       patternBuilder.WriteString(w)
   }
   patternBuilder.WriteString("<3")
   pattern := patternBuilder.String()

   // Check if pattern is a subsequence of s
   j := 0
   for i := 0; i < len(s) && j < len(pattern); i++ {
       if s[i] == pattern[j] {
           j++
       }
   }
   if j == len(pattern) {
       fmt.Println("yes")
   } else {
       fmt.Println("no")
   }
}
