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
   var n, q int
   fmt.Fscan(in, &n, &q)
   a := make([]int, n+1)
   maxVal := 0
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
       if a[i] > maxVal {
           maxVal = a[i]
       }
   }
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       adj[x] = append(adj[x], y)
       adj[y] = append(adj[y], x)
   }
   type Query struct{ typ, v, w int }
   queries := make([]Query, q)
   // collect update values
   updateCount := 0
   for i := 0; i < q; i++ {
       fmt.Fscan(in, &queries[i].typ)
       if queries[i].typ == 1 {
           fmt.Fscan(in, &queries[i].v)
       } else {
           fmt.Fscan(in, &queries[i].v, &queries[i].w)
           if queries[i].w > maxVal {
               maxVal = queries[i].w
           }
           updateCount++
       }
   }
   // sieve spf
   spf := make([]int, maxVal+1)
   for i := 2; i <= maxVal; i++ {
       if spf[i] == 0 {
           for j := i; j <= maxVal; j += i {
               if spf[j] == 0 {
                   spf[j] = i
               }
           }
       }
   }
   // factor initial values and update values, collect primes
   primesNodes := make([][]int, n+1)
   primeSet := make(map[int]struct{})
   // initial
   for i := 1; i <= n; i++ {
       x := a[i]
       var ps []int
       for x > 1 {
           p := spf[x]
           ps = append(ps, p)
           primeSet[p] = struct{}{}
           for x%p == 0 {
               x /= p
           }
       }
       primesNodes[i] = ps
   }
   // updates
   type Upd struct{ idx, v int; primes []int }
   updates := make([]Upd, 0, updateCount)
   for i, qu := range queries {
       if qu.typ == 2 {
           var ps []int
           x := qu.w
           for x > 1 {
               p := spf[x]
               ps = append(ps, p)
               primeSet[p] = struct{}{}
               for x%p == 0 {
                   x /= p
               }
           }
           updates = append(updates, Upd{i, qu.v, ps})
       }
   }
   // compress primes
   primeList := make([]int, 0, len(primeSet))
   for p := range primeSet {
       primeList = append(primeList, p)
   }
   // map to index
   pid := make(map[int]int, len(primeList))
   for i, p := range primeList {
       pid[p] = i
   }
   P := len(primeList)
   // build compressed primes for nodes
   primesIdx := make([][]int, n+1)
   for i := 1; i <= n; i++ {
       ps := primesNodes[i]
       idxs := make([]int, len(ps))
       for j, p := range ps {
           idxs[j] = pid[p]
       }
       primesIdx[i] = idxs
   }
   // compressed primes for updates
   for k := range updates {
       ps := updates[k].primes
       idxs := make([]int, len(ps))
       for j, p := range ps {
           idxs[j] = pid[p]
       }
       updates[k].primes = idxs
   }
   // build rooted tree children and depth
   children := make([][]int, n+1)
   depth := make([]int, n+1)
   parent := make([]int, n+1)
   // BFS
   queue := make([]int, 0, n)
   queue = append(queue, 1)
   parent[1] = 0
   depth[1] = 0
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           parent[v] = u
           depth[v] = depth[u] + 1
           children[u] = append(children[u], v)
           queue = append(queue, v)
       }
   }
   // prepare answer array
   ans := make([]int, q)
   // split segments
   type Seg struct{ l, r int }
   segs := make([]Seg, 0, len(updates)+1)
   start := 0
   updIdx := 0
   for i := 0; i < q; i++ {
       if queries[i].typ == 2 {
           if start <= i-1 {
               segs = append(segs, Seg{start, i - 1})
           }
           start = i + 1
           updIdx++
       }
   }
   if start <= q-1 {
       segs = append(segs, Seg{start, q - 1})
   }
   // stacks for primes
   stacks := make([][]int, P)
   // process segments
   for si, seg := range segs {
       // map node to its queries in seg
       qmap := make([][]int, n+1)
       for i := seg.l; i <= seg.r; i++ {
           if queries[i].typ == 1 {
               v := queries[i].v
               qmap[v] = append(qmap[v], i)
           }
       }
       // DFS
       var dfs func(u int)
       dfs = func(u int) {
           // answer queries at u
           for _, qi := range qmap[u] {
               best := -1
               bd := -1
               // for each prime in current value of u
               for _, pidx := range primesIdx[u] {
                   st := stacks[pidx]
                   if len(st) > 0 {
                       v := st[len(st)-1]
                       if depth[v] > bd {
                           bd = depth[v]
                           best = v
                       }
                   }
               }
               if best < 0 {
                   ans[qi] = -1
               } else {
                   ans[qi] = best
               }
           }
           // push u
           for _, pidx := range primesIdx[u] {
               stacks[pidx] = append(stacks[pidx], u)
           }
           // children
           for _, v := range children[u] {
               dfs(v)
           }
           // pop u
           for _, pidx := range primesIdx[u] {
               stacks[pidx] = stacks[pidx][:len(stacks[pidx])-1]
           }
       }
       dfs(1)
       // apply update if exists
       if si < len(updates) {
           up := updates[si]
           // update primesIdx for node
           primesIdx[up.v] = up.primes
       }
   }
   // output answers
   for i := 0; i < q; i++ {
       if queries[i].typ == 1 {
           fmt.Fprintln(out, ans[i])
       }
   }
}
