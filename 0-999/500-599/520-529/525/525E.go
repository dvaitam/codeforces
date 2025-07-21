package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   var S int64
   if _, err := fmt.Fscan(in, &n, &k, &S); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // precompute factorials up to 20 or until overflow > S
   maxF := int64(1)
   facts := make(map[int64]int64)
   facts[0] = 1
   facts[1] = 1
   for i := int64(2); ; i++ {
       maxF *= i
       if maxF > S {
           break
       }
       facts[i] = maxF
       if i >= 20 {
           break
       }
   }
   // split
   n1 := n / 2
   // map1: sum -> freq of ways with s stickers (0..n1)
   map1 := make(map[int64][]int64)
   // DFS for first half
   var dfs1 func(pos int, sum int64, used int)
   dfs1 = func(pos int, sum int64, used int) {
       if sum > S || used > k {
           return
       }
       if pos == n1 {
           arr, ok := map1[sum]
           if !ok {
               arr = make([]int64, n1+1)
               map1[sum] = arr
           }
           arr[used]++
           return
       }
       // skip
       dfs1(pos+1, sum, used)
       // choose without sticker
       s := sum + a[pos]
       if s <= S {
           dfs1(pos+1, s, used)
       }
       // choose with sticker if possible
       if used < k {
           if f, ok := facts[a[pos]]; ok {
               s2 := sum + f
               if s2 <= S {
                   dfs1(pos+1, s2, used+1)
               }
           }
       }
   }
   dfs1(0, 0, 0)
   // DFS for second half and combine
   ans := int64(0)
   var dfs2 func(pos int, sum int64, used int)
   dfs2 = func(pos int, sum int64, used int) {
       if sum > S || used > k {
           return
       }
       if pos == n {
           rem := S - sum
           if freq1, ok := map1[rem]; ok {
               // sum over s1 + used <= k
               maxS1 := k - used
               if maxS1 > n1 {
                   maxS1 = n1
               }
               for s1 := 0; s1 <= maxS1; s1++ {
                   ans += freq1[s1]
               }
           }
           return
       }
       // skip
       dfs2(pos+1, sum, used)
       // choose without sticker
       s := sum + a[pos]
       if s <= S {
           dfs2(pos+1, s, used)
       }
       // choose with sticker
       if used < k {
           if f, ok := facts[a[pos]]; ok {
               s2 := sum + f
               if s2 <= S {
                   dfs2(pos+1, s2, used+1)
               }
           }
       }
   }
   dfs2(n1, 0, 0)
   // output
   w := bufio.NewWriter(os.Stdout)
   fmt.Fprint(w, ans)
   w.Flush()
}
