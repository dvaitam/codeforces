package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   parent     []int
   leftChild  []int
   rightChild []int
)

// find returns the representative of x with path compression.
func find(x int) int {
   if parent[x] != x {
       parent[x] = find(parent[x])
   }
   return parent[x]
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // Allocate slices for union-find and merge tree (1-based indexing).
   size := 2*n + 2
   parent = make([]int, size)
   leftChild = make([]int, size)
   rightChild = make([]int, size)
   for i := 1; i <= n; i++ {
       parent[i] = i
   }
   // Build merge tree: each union creates a new node n+i.
   for i := 1; i < n; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       fu := find(u)
       fv := find(v)
       newNode := n + i
       parent[fu] = newNode
       parent[fv] = newNode
       parent[newNode] = newNode
       leftChild[newNode] = fu
       rightChild[newNode] = fv
   }
   // The root of the merge tree is node 2*n-1.
   root := 2*n - 1
   // Iterative DFS to collect leaves in left-to-right order.
   res := make([]int, 0, n)
   stack := make([]int, 0, 2*n)
   stack = append(stack, root)
   for len(stack) > 0 {
       node := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       if node <= n {
           res = append(res, node)
       } else {
           // Push right then left so that left is processed first.
           stack = append(stack, rightChild[node])
           stack = append(stack, leftChild[node])
       }
   }
   // Output the initial arrangement.
   for i, v := range res {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, v)
   }
   out.WriteByte('\n')
}
