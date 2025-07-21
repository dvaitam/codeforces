package main

import (
   "bufio"
   "fmt"
   "os"
)

// Node represents a trie node for reversed strings
type Node struct {
   ch  [26]int
   cnt int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   // Initialize trie with root node
   nodes := make([]Node, 1, 100000+5)
   ans := 1
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       p := 0
       // Insert reversed string into trie, counting passes
       for j := len(s) - 1; j >= 0; j-- {
           c := s[j] - 'a'
           if nodes[p].ch[c] == 0 {
               nodes[p].ch[c] = len(nodes)
               nodes = append(nodes, Node{})
           }
           p = nodes[p].ch[c]
           nodes[p].cnt++
           if nodes[p].cnt > ans {
               ans = nodes[p].cnt
           }
       }
   }
   if n == 0 {
       fmt.Fprintln(writer, 0)
   } else {
       fmt.Fprintln(writer, ans)
   }
}
