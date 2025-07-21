package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, u, r int
   b, k []int
   p []int
   maxScore int64
   arr [][]int
)

func dfs(depth int, prevOp1 bool) {
   if depth == u {
       // compute score
       var s int64
       cur := arr[depth]
       for i := 0; i < n; i++ {
           s += int64(cur[i]) * int64(k[i])
       }
       if s > maxScore {
           maxScore = s
       }
       return
   }
   cur := arr[depth]
   nxt := arr[depth+1]
   // op1: xor, skip if previous was op1
   if !prevOp1 {
       for i := 0; i < n; i++ {
           nxt[i] = cur[i] ^ b[i]
       }
       dfs(depth+1, true)
   }
   // op2: permute and add r
   for i := 0; i < n; i++ {
       nxt[i] = cur[p[i]] + r
   }
   dfs(depth+1, false)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   // read input
   fmt.Fscan(in, &n, &u, &r)
   a := make([]int, n)
   b = make([]int, n)
   k = make([]int, n)
   p = make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &b[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &k[i])
   }
   for i := 0; i < n; i++ {
       var pi int
       fmt.Fscan(in, &pi)
       p[i] = pi - 1
   }
   // allocate arrays
   arr = make([][]int, u+1)
   for i := 0; i <= u; i++ {
       arr[i] = make([]int, n)
   }
   copy(arr[0], a)
   maxScore = -1 << 60
   dfs(0, false)
   // output result
   fmt.Println(maxScore)
}
