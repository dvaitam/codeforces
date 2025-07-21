package main

import (
   "bufio"
   "fmt"
   "os"
)

// Node represents a trie node
type Node struct {
   next  [26]*Node
   win   bool
   lose  bool
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   var k int64
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   root := &Node{}
   // build trie
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       cur := root
       for j := 0; j < len(s); j++ {
           c := s[j] - 'a'
           if cur.next[c] == nil {
               cur.next[c] = &Node{}
           }
           cur = cur.next[c]
       }
   }
   // dfs to compute win/lose
   var dfs func(u *Node)
   dfs = func(u *Node) {
       hasChild := false
       u.win = false
       u.lose = false
       for c := 0; c < 26; c++ {
           v := u.next[c]
           if v == nil {
               continue
           }
           hasChild = true
           dfs(v)
           if !v.win {
               u.win = true
           }
           if !v.lose {
               u.lose = true
           }
       }
       if !hasChild {
           // leaf
           u.win = false
           u.lose = true
       }
   }
   dfs(root)
   // decide result
   if !root.win {
       fmt.Fprint(writer, "Second")
   } else if root.lose {
       fmt.Fprint(writer, "First")
   } else {
       if k%2 == 1 {
           fmt.Fprint(writer, "First")
       } else {
           fmt.Fprint(writer, "Second")
       }
   }
}
