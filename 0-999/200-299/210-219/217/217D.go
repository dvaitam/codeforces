package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, t int
   if _, err := fmt.Fscan(in, &n, &m, &t); err != nil {
       return
   }
   cnt := make([]int, m)
   for i := 0; i < t; i++ {
       var x int
       fmt.Fscan(in, &x)
       r := x % m
       cnt[r]++
   }
   // Build orbits of residues
   type orbit struct{ r, cr, cmr int }
   orbits := make([]orbit, 0, m/2+1)
   for r := 1; r*2 < m; r++ {
       cr := cnt[r]
       cmr := cnt[m-r]
       if cr+cmr > 0 {
           orbits = append(orbits, orbit{r, cr, cmr})
       }
   }
   // if m even, handle r = m/2
   if m%2 == 0 {
       r := m / 2
       if cnt[r] > 0 {
           orbits = append(orbits, orbit{r, cnt[r], 0})
       }
   }
   // dp: map mask string to weight
   // mask length m, '0' or '1'
   start := make([]byte, m)
   for i := range start {
       start[i] = '0'
   }
   start[0] = '1'
   dp := map[string]int{
       string(start): 1,
   }
   // process orbits
   for _, o := range orbits {
       dp2 := make(map[string]int)
       r := o.r
       tot := o.cr + o.cmr
       for maskStr, w := range dp {
           // no pick
           dp2[maskStr] = (dp2[maskStr] + w) % mod
           if tot == 0 {
               continue
           }
           mask := []byte(maskStr)
           // check safe: no existing sum == r or m-r
           if mask[r] == '1' || mask[(m-r)%m] == '1' {
               continue
           }
           // pick one from orbit, weight tot
           // expand mask
           newMask := make([]byte, m)
           copy(newMask, mask)
           for j := 0; j < m; j++ {
               if mask[j] == '1' {
                   // step +r
                   newMask[(j+r)%m] = '1'
                   // step -r
                   newMask[(j+m-r)%m] = '1'
               }
           }
           nmStr := string(newMask)
           dp2[nmStr] = (dp2[nmStr] + int((int64(w) * int64(tot)) % mod)) % mod
       }
       dp = dp2
   }
   // sum all weights
   ans := 0
   for _, v := range dp {
       ans = (ans + v) % mod
   }
   fmt.Println(ans)
}
