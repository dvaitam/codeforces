package main

import (
   "bufio"
   "fmt"
   "os"
)

type pair struct {
   val int
   pos int
}

func minPair(a, b pair) pair {
   if a.val <= b.val {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var N, Q int
   fmt.Fscan(in, &N, &Q)
   A := make([]int, N+1)
   for i := 1; i <= N; i++ {
       fmt.Fscan(in, &A[i])
   }
   adj := make([][]int, N+1)
   for i := 2; i <= N; i++ {
       var p int
       fmt.Fscan(in, &p)
       adj[p] = append(adj[p], i)
   }
   inTime := make([]int, N+1)
   outTime := make([]int, N+1)
   depth := make([]int, N+1)
   nodeAt := make([]int, N)
   t := 0
   var dfs func(u int)
   dfs = func(u int) {
       inTime[u] = t
       nodeAt[t] = u
       t++
       for _, v := range adj[u] {
           depth[v] = depth[u] + 1
           dfs(v)
       }
       outTime[u] = t
   }
   dfs(1)
   // build segment tree for R values
   n2 := 1
   for n2 < N {
       n2 <<= 1
   }
   INF := 1<<60
   tree := make([]pair, 2*n2)
   // init leaves
   for i := 0; i < n2; i++ {
       if i < N {
           tree[n2+i] = pair{1, i}
       } else {
           tree[n2+i] = pair{INF, i}
       }
   }
   // build
   for i := n2 - 1; i >= 1; i-- {
       tree[i] = minPair(tree[2*i], tree[2*i+1])
   }
   // function to update position p to value v
   update := func(p, v int) {
       idx := p + n2
       tree[idx].val = v
       for idx >>= 1; idx >= 1; idx >>= 1 {
           tree[idx] = minPair(tree[2*idx], tree[2*idx+1])
       }
   }
   // query min in [l, r)
   queryMin := func(l, r int) pair {
       l += n2; r += n2
       res := pair{INF, -1}
       for l < r {
           if l&1 == 1 {
               res = minPair(res, tree[l])
               l++
           }
           if r&1 == 1 {
               r--
               res = minPair(res, tree[r])
           }
           l >>= 1; r >>= 1
       }
       return res
   }
   // process queries
   for day := 1; day <= Q; day++ {
       var X int
       fmt.Fscan(in, &X)
       l := inTime[X]
       r := outTime[X]
       cnt := 0
       sumDist := 0
       for {
           mn := queryMin(l, r)
           if mn.val > day {
               break
           }
           idx := mn.pos
           u := nodeAt[idx]
           cnt++
           sumDist += depth[u] - depth[X]
           // next ready time
           update(idx, day + A[u])
       }
       fmt.Fprintln(out, sumDist, cnt)
   }
}
