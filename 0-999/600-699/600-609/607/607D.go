package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var v1 int64
   var q int
   if _, err := fmt.Fscan(in, &v1, &q); err != nil {
       return
   }
   maxN := q + 5
   parent := make([]int, maxN)
   cnt := make([]int64, maxN)
   sum := make([]int64, maxN)
   // initialize root
   parent[1] = 0
   cnt[1] = 1
   sum[1] = v1 % mod

   nextID := 2
   for i := 0; i < q; i++ {
       var t int
       fmt.Fscan(in, &t)
       if t == 1 {
           var p int
           var v int64
           fmt.Fscan(in, &p, &v)
           id := nextID
           nextID++
           parent[id] = p
           cnt[id] = 1
           sum[id] = v % mod
           // update at p
           // old values
           oldCnt := cnt[p]
           oldSum := sum[p]
           // update p
           sum[p] = (sum[p] + sum[id]) % mod
           cnt[p]++
           // compute delta of power at p: new cnt*new sum - old cnt*old sum
           newP := cnt[p] * sum[p] % mod
           oldP := oldCnt * oldSum % mod
           delta := newP - oldP
           if delta < 0 {
               delta += mod
           }
           // propagate to ancestors
           u := parent[p]
           for u != 0 && delta != 0 {
               // sum[u] increases by delta
               sum[u] = (sum[u] + delta) % mod
               // new delta is cnt[u] * delta
               delta = delta * cnt[u] % mod
               u = parent[u]
           }
       } else if t == 2 {
           var u int
           fmt.Fscan(in, &u)
           ans := cnt[u] * sum[u] % mod
           fmt.Fprintln(out, ans)
       }
   }
}
