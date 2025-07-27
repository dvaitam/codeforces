package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const MOD = 998244353

func add(a, b int) int {
   s := a + b
   if s >= MOD {
       s -= MOD
   }
   return s
}

func mul(a, b int) int {
   return int((int64(a) * int64(b)) % MOD)
}

func powmod(a, e int) int {
   res := 1
   base := a
   for e > 0 {
       if e&1 == 1 {
           res = mul(res, base)
       }
       base = mul(base, base)
       e >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   fmt.Fscan(reader, &n, &k)
   events := make([][2]int, 0, 2*n)
   for i := 0; i < n; i++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       // start at l, end at r+1
       events = append(events, [2]int{l, 1})
       events = append(events, [2]int{r + 1, -1})
   }
   sort.Slice(events, func(i, j int) bool {
       if events[i][0] != events[j][0] {
           return events[i][0] < events[j][0]
       }
       // for same time, process end events before start events
       return events[i][1] < events[j][1]
   })

   // precompute factorials
   maxN := n
   fact := make([]int, maxN+1)
   ifact := make([]int, maxN+1)
   fact[0] = 1
   for i := 1; i <= maxN; i++ {
       fact[i] = mul(fact[i-1], i)
   }
   ifact[maxN] = powmod(fact[maxN], MOD-2)
   for i := maxN; i > 0; i-- {
       ifact[i-1] = mul(ifact[i], i)
   }

   comb := func(a, b int) int {
       if b < 0 || b > a {
           return 0
       }
       return mul(fact[a], mul(ifact[b], ifact[a-b]))
   }

   cur := 0
   ans := 0
   for _, ev := range events {
       t, typ := ev[0], ev[1]
       if typ == 1 {
           // choose this and k-1 from cur
           if k-1 <= cur {
               ans = add(ans, comb(cur, k-1))
           }
           cur++
       } else {
           cur--
       }
       _ = t // unused
   }
   fmt.Fprintln(writer, ans)
}
