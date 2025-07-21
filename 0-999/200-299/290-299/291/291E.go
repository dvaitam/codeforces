package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   children := make([][]int, n+1)
   strs := make([][]byte, n+1)
   // read edges: for nodes 2..n
   for v := 2; v <= n; v++ {
       var p int
       var s string
       fmt.Fscan(reader, &p, &s)
       children[p] = append(children[p], v)
       strs[v] = []byte(s)
   }
   var tstr string
   fmt.Fscan(reader, &tstr)
   t := []byte(tstr)
   m := len(t)
   // build prefix function
   pi := make([]int, m)
   for i := 1; i < m; i++ {
       j := pi[i-1]
       for j > 0 && t[j] != t[i] {
           j = pi[j-1]
       }
       if t[j] == t[i] {
           j++
       }
       pi[i] = j
   }
   var ans int64
   // iterative DFS stack
   type frame struct {
       u          int
       curBefore  int
       curAfter   int
       childIndex int
   }
   stack := make([]frame, 0, n)
   // start at root 1
   stack = append(stack, frame{u: 1, curBefore: 0, childIndex: -1})
   for len(stack) > 0 {
       fr := &stack[len(stack)-1]
       if fr.childIndex == -1 {
           // first time visit: process its string
           cur := fr.curBefore
           if fr.u > 1 {
               for _, c := range strs[fr.u] {
                   // KMP step
                   for cur > 0 && (cur == m || t[cur] != c) {
                       cur = pi[cur-1]
                   }
                   if cur < m && t[cur] == c {
                       cur++
                   }
                   if cur == m {
                       ans++
                   }
               }
           }
           fr.curAfter = cur
           fr.childIndex = 0
       } else if fr.childIndex < len(children[fr.u]) {
           v := children[fr.u][fr.childIndex]
           fr.childIndex++
           // push child
           stack = append(stack, frame{u: v, curBefore: fr.curAfter, childIndex: -1})
       } else {
           // done with this node
           stack = stack[:len(stack)-1]
       }
   }
   fmt.Fprint(writer, ans)
}
