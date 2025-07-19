package main

import (
   "bufio"
   "fmt"
   "os"
)

const MX = 5032107 + 5

var (
   goArr [MX][8]int
   gl    [20]int
   bb    [MX]int
   primes []int
   vv    []int
   gst   int
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func run(x, cur, cnt int) {
   if x == len(vv) {
       for i := 0; i < 8; i++ {
           if gl[i+cnt] < goArr[cur][i] {
               gl[i+cnt] = goArr[cur][i]
           }
       }
       goArr[cur][cnt] = gst
       return
   }
   // include this prime (count removal)
   run(x+1, cur, cnt+1)
   // exclude this prime (multiply current)
   run(x+1, cur*vv[x], cnt)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   // sieve smallest prime factors
   for i := 1; i < MX; i++ {
       bb[i] = i
   }
   for i := 2; i < MX; i++ {
       if bb[i] == i {
           primes = append(primes, i)
       }
       for _, p := range primes {
           if i*p >= MX || bb[i] < p {
               break
           }
           bb[i*p] = p
       }
   }

   var n, q int
   fmt.Fscan(reader, &n, &q)
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   type pair struct{ l, idx int }
   gg := make([][]pair, n+1)
   for i := 0; i < q; i++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       if r <= n {
           gg[r] = append(gg[r], pair{l, i})
       }
   }
   ans := make([]int, q)
   // process
   for i := 1; i <= n; i++ {
       gst = i
       // factor to squarefree primes
       vv = vv[:0]
       now := a[i]
       for now != 1 {
           x := bb[now]
           cnt := 0
           for now%x == 0 {
               now /= x
               cnt ^= 1
           }
           if cnt == 1 {
               vv = append(vv, x)
           }
       }
       run(0, 1, 0)
       // answer queries ending at i
       for _, pr := range gg[i] {
           l := pr.l
           best := 20
           for j := 0; j < 20; j++ {
               if gl[j] >= l && j < best {
                   best = j
               }
           }
           ans[pr.idx] = best
       }
   }
   for i := 0; i < q; i++ {
       fmt.Fprintln(writer, ans[i])
   }
}
