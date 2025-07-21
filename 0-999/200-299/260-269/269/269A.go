package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const INF = int64(1e18)

func mulCap(a, b int64) int64 {
   if a == 0 || b == 0 {
       return 0
   }
   if a > INF/b {
       return INF
   }
   return a * b
}

// pow4Cap returns min(4^d, INF)
func pow4Cap(d int64) int64 {
   res := int64(1)
   base := int64(4)
   for d > 0 {
       if d&1 == 1 {
           res = mulCap(res, base)
       }
       base = mulCap(base, base)
       d >>= 1
   }
   return res
}

// valid checks if a container of size 2^p can hold all boxes
func valid(p int64, ks []int64, as []int64) bool {
   slots := int64(1)
   last := p
   for i, k := range ks {
       if k > p {
           continue
       }
       // scale down from level last to k
       d := last - k
       if d > 0 {
           slots = mulCap(slots, pow4Cap(d))
       }
       // need to place as[i] boxes of size k
       if slots < as[i] {
           return false
       }
       slots -= as[i]
       last = k
   }
   // finally scale to level 0
   if last > 0 {
       slots = mulCap(slots, pow4Cap(last))
   }
   // all boxes placed (a0 accounted in ks/as if k=0)
   return true
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   m := n
   ks := make([]int64, 0, m+1)
   as := make([]int64, 0, m+1)
   have0 := false
   maxk := int64(0)
   for i := 0; i < n; i++ {
       var k, a int64
       fmt.Fscan(in, &k, &a)
       ks = append(ks, k)
       as = append(as, a)
       if k == 0 {
           have0 = true
       }
       if k > maxk {
           maxk = k
       }
   }
   if !have0 {
       ks = append(ks, 0)
       as = append(as, 0)
   }
   // sort descending by k
   idx := make([]int, len(ks))
   for i := range idx {
       idx[i] = i
   }
   sort.Slice(idx, func(i, j int) bool {
       return ks[idx[i]] > ks[idx[j]]
   })
   k2 := make([]int64, len(ks))
   a2 := make([]int64, len(as))
   for i, id := range idx {
       k2[i] = ks[id]
       a2[i] = as[id]
   }
   // binary search p
   lo := maxk
   hi := maxk + 200
   for lo < hi {
       mid := lo + (hi-lo)/2
       if valid(mid, k2, a2) {
           hi = mid
       } else {
           lo = mid + 1
       }
   }
   fmt.Println(lo)
}
