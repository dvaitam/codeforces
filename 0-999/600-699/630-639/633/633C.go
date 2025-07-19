package main

import (
   "bufio"
   "fmt"
   "os"
)

// Node represents a trie node
type Node struct {
   next [26]int32
   vis  int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   var s string
   fmt.Fscan(reader, &n, &s, &m)
   // reverse s
   sb := []byte(s)
   for i, j := 0, len(sb)-1; i < j; i, j = i+1, j-1 {
       sb[i], sb[j] = sb[j], sb[i]
   }
   // read words
   words := make([]string, m+1)
   // build trie
   nodes := make([]Node, 1)
   for i := 1; i <= m; i++ {
       var w string
       fmt.Fscan(reader, &w)
       words[i] = w
       // insert word into trie (lowercased)
       x := int32(0)
       for j := 0; j < len(w); j++ {
           c := w[j]
           // to lowercase
           if c >= 'A' && c <= 'Z' {
               c |= 32
           }
           idx := c - 'a'
           if nodes[x].next[idx] == 0 {
               nodes = append(nodes, Node{})
               nodes[x].next[idx] = int32(len(nodes) - 1)
           }
           x = nodes[x].next[idx]
       }
       nodes[x].vis = i
   }
   // dp: -1 unknown, 0 false, 1 true
   dp := make([]int8, n+1)
   for i := 0; i <= n; i++ {
       dp[i] = -1
   }
   var ans []int
   var calc func(int) bool
   calc = func(i int) bool {
       if i == n {
           return true
       }
       if dp[i] != -1 {
           return dp[i] == 1
       }
       x := int32(0)
       // traverse trie along sb
       for j := i; j < n; j++ {
           c := sb[j]
           idx := c - 'a'
           if idx < 0 || idx >= 26 || nodes[x].next[idx] == 0 {
               dp[i] = 0
               return false
           }
           x = nodes[x].next[idx]
           if nodes[x].vis != 0 {
               if calc(j + 1) {
                   ans = append(ans, nodes[x].vis)
                   dp[i] = 1
                   return true
               }
           }
       }
       dp[i] = 0
       return false
   }
   calc(0)
   // output in order
   for i := 0; i < len(ans); i++ {
       fmt.Fprint(writer, words[ans[i]])
       if i+1 < len(ans) {
           writer.WriteByte(' ')
       }
   }
   writer.WriteByte('\n')
}
