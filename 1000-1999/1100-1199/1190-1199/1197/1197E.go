package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const MOD = 1000000007
const INF = int64(4e18)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   dolls := make([]struct{out, inV int}, n)
   outs := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &dolls[i].out, &dolls[i].inV)
       outs[i] = dolls[i].out
   }
   // find global min out and max in
   Omin := outs[0]
   Imax := dolls[0].inV
   for i := 1; i < n; i++ {
       if dolls[i].out < Omin {
           Omin = dolls[i].out
       }
       if dolls[i].inV > Imax {
           Imax = dolls[i].inV
       }
   }
   // compress out values
   sort.Ints(outs)
   uouts := outs[:1]
   for i := 1; i < n; i++ {
       if outs[i] != outs[i-1] {
           uouts = append(uouts, outs[i])
       }
   }
   m := len(uouts)
   // map out to index
   // we will use 0..m-1
   // tree arrays
   fcoord := make([]int64, m)
   ccoord := make([]int, m)
   for i := 0; i < m; i++ {
       fcoord[i] = INF
       ccoord[i] = 0
   }
   // sort dolls by inV asc
   sort.Slice(dolls, func(i, j int) bool {
       if dolls[i].inV != dolls[j].inV {
           return dolls[i].inV < dolls[j].inV
       }
       return dolls[i].out < dolls[j].out
   })
   // dp arrays per doll
   dp := make([]int64, n)
   cnt := make([]int, n)
   // process
   for i, d := range dolls {
       dp[i] = INF
       cnt[i] = 0
       // base if start
       if d.out == Omin {
           dp[i] = int64(d.inV)
           cnt[i] = 1
       }
       // find max out <= inV
       r := sort.Search(len(uouts), func(j int) bool { return uouts[j] > d.inV })
       // r is first index with out > inV, so coordinates [0..r-1]
       if r > 0 {
           // M = uouts[r-1]
           if fcoord[r-1] < INF {
               val := fcoord[r-1] + int64(d.inV)
               if val < dp[i] {
                   dp[i] = val
                   cnt[i] = ccoord[r-1]
               } else if val == dp[i] {
                   cnt[i] = (cnt[i] + ccoord[r-1]) % MOD
               }
           }
       }
       if dp[i] < INF {
           // update at out index
           pos := sort.Search(len(uouts), func(j int) bool { return uouts[j] >= d.out })
           // pos < m and uouts[pos] == d.out
           fnew := dp[i] - int64(d.out)
           if fnew < fcoord[pos] {
               fcoord[pos] = fnew
               ccoord[pos] = cnt[i]
           } else if fnew == fcoord[pos] {
               ccoord[pos] = (ccoord[pos] + cnt[i]) % MOD
           }
       }
   }
   // collect ends with inV == Imax
   best := INF
   ways := 0
   for i, d := range dolls {
       if d.inV == Imax && dp[i] < INF {
           if dp[i] < best {
               best = dp[i]
               ways = cnt[i]
           } else if dp[i] == best {
               ways = (ways + cnt[i]) % MOD
           }
       }
   }
   fmt.Fprint(out, ways)
}
