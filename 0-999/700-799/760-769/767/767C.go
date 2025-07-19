package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   parent := make([]int, n+1)
   v := make([]int64, n+1)
   children := make([][]int, n+1)
   var root int
   var tot int64
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &parent[i], &v[i])
       tot += v[i]
   }
   for i := 1; i <= n; i++ {
       if parent[i] == 0 {
           root = i
       } else {
           children[parent[i]] = append(children[parent[i]], i)
       }
   }
   if tot%3 != 0 {
       fmt.Println(-1)
       return
   }
   target := tot / 3
   // post-order traversal to compute subtree sums
   sum := make([]int64, n+1)
   var ans1, ans2 int
   type frame struct{ node int; done bool }
   stack := make([]frame, 0, n*2)
   stack = append(stack, frame{node: root, done: false})
   for len(stack) > 0 && (ans1 == 0 || ans2 == 0) {
       f := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       if f.done {
           // compute sum for node
           sum[f.node] = v[f.node]
           for _, c := range children[f.node] {
               if sum[c] == target {
                   if ans1 == 0 {
                       ans1 = c
                   } else if ans2 == 0 {
                       ans2 = c
                   }
               } else {
                   sum[f.node] += sum[c]
               }
           }
       } else {
           // push node for post processing
           stack = append(stack, frame{node: f.node, done: true})
           // push children for processing
           for _, c := range children[f.node] {
               stack = append(stack, frame{node: c, done: false})
           }
       }
   }
   if ans1 == 0 || ans2 == 0 {
       fmt.Println(-1)
   } else {
       fmt.Printf("%d %d\n", ans1, ans2)
   }
}
