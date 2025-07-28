package main

import (
   "bufio"
   "fmt"
   "os"
)

// bitset for DP
type bitset []uint64

func newBitset(nbits int) bitset {
   n := (nbits + 63) >> 6
   return make(bitset, n)
}
// set all bits [0..nbits)
func (b bitset) setAll(nbits int) {
   for i := range b {
       b[i] = ^uint64(0)
   }
   // clear bits beyond nbits
   tail := nbits & 63
   if tail != 0 {
       b[len(b)-1] &= (uint64(1)<<tail - 1)
   }
}
// check any bit set
func (b bitset) any() bool {
   for _, w := range b {
       if w != 0 {
           return true
       }
   }
   return false
}
// shift left by k bits, store in dst
func (dst bitset) shl(src bitset, k int) {
   n := len(src)
   wordShift := k >> 6
   offset := uint(k & 63)
   for i := n - 1; i >= 0; i-- {
       var v uint64
       j := i - wordShift
       if j >= 0 {
           v = src[j] << offset
           if offset != 0 && j-1 >= 0 {
               v |= src[j-1] >> (64 - offset)
           }
       }
       dst[i] = v
   }
}
// shift right by k bits, store in dst
func (dst bitset) shr(src bitset, k int) {
   n := len(src)
   wordShift := k >> 6
   offset := uint(k & 63)
   for i := 0; i < n; i++ {
       var v uint64
       j := i + wordShift
       if j < n {
           v = src[j] >> offset
           if offset != 0 && j+1 < n {
               v |= src[j+1] << (64 - offset)
           }
       }
       dst[i] = v
   }
}

// check feasibility for given max range d and segment lengths a
func feasible(a []int, d int) bool {
   // dp_prev bits for Q[k] in [0..d]
   bsize := d + 1
   dpPrev := newBitset(bsize)
   dpPrev.setAll(bsize)
   dpCur := newBitset(bsize)
   tmp := newBitset(bsize)
   for _, ai := range a {
       // dpCur = (dpPrev << ai) | (dpPrev >> ai)
       // shift left
       tmp.shl(dpPrev, ai)
       // mask beyond bsize
       // tmp bits beyond are automatically dropped by size
       // copy to dpCur
       for i := range dpCur {
           dpCur[i] = tmp[i]
       }
       // shift right into tmp
       tmp.shr(dpPrev, ai)
       // or
       for i := range dpCur {
           dpCur[i] |= tmp[i]
       }
       // clear bits beyond bsize
       tail := bsize & 63
       if tail != 0 {
           mask := uint64(1)<<uint(tail) - 1
           dpCur[len(dpCur)-1] &= mask
       }
       if !dpCur.any() {
           return false
       }
       // swap dpPrev, dpCur
       dpPrev, dpCur = dpCur, dpPrev
   }
   return true
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(in, &n)
       a := make([]int, n)
       maxa := 0
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &a[i])
           if a[i] > maxa {
               maxa = a[i]
           }
       }
       lo, hi := maxa, 2*maxa
       ans := hi
       for lo <= hi {
           mid := (lo + hi) >> 1
           if feasible(a, mid) {
               ans = mid
               hi = mid - 1
           } else {
               lo = mid + 1
           }
       }
       fmt.Fprintln(out, ans)
   }
}
