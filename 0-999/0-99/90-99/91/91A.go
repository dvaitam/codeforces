package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // Read s1
   s1line, err := reader.ReadString('\n')
   if err != nil && len(s1line) == 0 {
       return
   }
   s1 := strings.TrimSpace(s1line)
   // Read s2
   s2line, err := reader.ReadString('\n')
   if err != nil && len(s2line) == 0 {
       return
   }
   s2 := strings.TrimSpace(s2line)

   n := len(s1)
   const K = 26
   // next[i][c]: next occurrence of c at or after i in s1, or -1
   next := make([][K]int, n+1)
   for c := 0; c < K; c++ {
       next[n][c] = -1
   }
   for i := n - 1; i >= 0; i-- {
       next[i] = next[i+1]
       next[i][s1[i]-'a'] = i
   }

   res := 1
   pos := 0
   for i := 0; i < len(s2); i++ {
       c := s2[i] - 'a'
       // If the character can be matched in current copy
       if pos <= n && next[pos][c] != -1 {
           pos = next[pos][c] + 1
       } else {
           // Start a new copy of s1
           if next[0][c] == -1 {
               fmt.Println(-1)
               return
           }
           res++
           pos = next[0][c] + 1
       }
   }
   fmt.Println(res)
