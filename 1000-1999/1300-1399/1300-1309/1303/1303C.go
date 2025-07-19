package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   fmt.Fscan(in, &t)
   for t > 0 {
       t--
       var s string
       fmt.Fscan(in, &s)
       // build adjacency
       adj := make([][]bool, 26)
       for i := range adj {
           adj[i] = make([]bool, 26)
       }
       n := len(s)
       for i := 0; i < n; i++ {
           cur := s[i] - 'a'
           if i > 0 {
               prev := s[i-1] - 'a'
               adj[cur][prev] = true
               adj[prev][cur] = true
           }
           if i < n-1 {
               next := s[i+1] - 'a'
               adj[cur][next] = true
               adj[next][cur] = true
           }
       }
       visited := make([]bool, 26)
       connections := make([]int, 26)
       poss := true
       for i := 0; i < 26; i++ {
           cnt := 0
           for j := 0; j < 26; j++ {
               if adj[i][j] {
                   cnt++
               }
           }
           if cnt > 2 {
               poss = false
           }
           connections[i] = cnt
       }
       var ans []byte
       // traverse each component starting at nodes with degree <=1
       for i := 0; i < 26; i++ {
           if !visited[i] && connections[i] <= 1 {
               // stack holds pairs [node, parent]
               stack := [][2]int{{i, -1}}
               for len(stack) > 0 {
                   top := stack[len(stack)-1]
                   stack = stack[:len(stack)-1]
                   cur, parent := top[0], top[1]
                   if visited[cur] {
                       continue
                   }
                   visited[cur] = true
                   ans = append(ans, byte(cur)+'a')
                   for c := 0; c < 26; c++ {
                       if adj[cur][c] {
                           if c == parent {
                               continue
                           }
                           if visited[c] {
                               poss = false
                           } else {
                               stack = append(stack, [2]int{c, cur})
                           }
                       }
                   }
               }
           }
       }
       if !poss || len(ans) != 26 {
           fmt.Fprintln(out, "NO")
       } else {
           fmt.Fprintln(out, "YES")
           fmt.Fprintln(out, string(ans))
       }
   }
}
