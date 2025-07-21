package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 777777777

type Node struct {
   dp [3][3]int64
}

func combine(a, b *Node, w *[3][3]int) Node {
   var c Node
   for i := 0; i < 3; i++ {
       for j := 0; j < 3; j++ {
           var s int64
           for li := 0; li < 3; li++ {
               if a.dp[i][li] == 0 {
                   continue
               }
               for rj := 0; rj < 3; rj++ {
                   if w[li][rj] == 0 || b.dp[rj][j] == 0 {
                       continue
                   }
                   s += a.dp[i][li] * b.dp[rj][j] % mod
               }
           }
           c.dp[i][j] = s % mod
       }
   }
   return c
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m int
   fmt.Fscan(in, &n, &m)
   var w [3][3]int
   for i := 0; i < 3; i++ {
       for j := 0; j < 3; j++ {
           fmt.Fscan(in, &w[i][j])
       }
   }
   // build segment tree
   size := 1
   for size < n {
       size <<= 1
   }
   tree := make([]Node, 2*size)
   // initialize leaves
   for i := 0; i < n; i++ {
       idx := size + i
       for x := 0; x < 3; x++ {
           tree[idx].dp[x][x] = 1
       }
   }
   // others leave default zero but unused
   // build inner nodes
   for i := size - 1; i > 0; i-- {
       tree[i] = combine(&tree[2*i], &tree[2*i+1], &w)
   }
   // process queries
   for qi := 0; qi < m; qi++ {
       var vi, ti int
       fmt.Fscan(in, &vi, &ti)
       pos := size + vi - 1
       // reset leaf
       for x := 0; x < 3; x++ {
           for y := 0; y < 3; y++ {
               tree[pos].dp[x][y] = 0
           }
       }
       if ti == 0 {
           // unassigned
           for x := 0; x < 3; x++ {
               tree[pos].dp[x][x] = 1
           }
       } else {
           t := ti - 1
           tree[pos].dp[t][t] = 1
       }
       // update up
       for p := pos >> 1; p > 0; p >>= 1 {
           tree[p] = combine(&tree[2*p], &tree[2*p+1], &w)
       }
       // answer is sum over root dp
       var ans int64
       for x := 0; x < 3; x++ {
           for y := 0; y < 3; y++ {
               ans += tree[1].dp[x][y]
           }
       }
       ans %= mod
       if qi > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, ans)
   }
   out.WriteByte('\n')
}
