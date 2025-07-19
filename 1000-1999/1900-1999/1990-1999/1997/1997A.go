package main

import (
   "bufio"
   "fmt"
   "os"
)

// helper computes the score of string s: starts with 2, adds 1 if current and previous chars match, else adds 2
func helper(s string) int {
   ans := 2
   for i := 1; i < len(s); i++ {
       if s[i] == s[i-1] {
           ans++
       } else {
           ans += 2
       }
   }
   return ans
}

// solve processes one test case: reads string s, inserts one lowercase letter to maximize helper score
func solve(reader *bufio.Reader, writer *bufio.Writer) {
   var s string
   fmt.Fscan(reader, &s)
   var ans string
   curr := 0
   n := len(s)
   for i := 0; i <= n; i++ {
       for ch := byte('a'); ch <= 'z'; ch++ {
           modified := s[:i] + string(ch) + s[i:]
           score := helper(modified)
           if score > curr {
               curr = score
               ans = modified
           }
       }
   }
   fmt.Fprintln(writer, ans)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for i := 0; i < t; i++ {
       solve(reader, writer)
   }
}
