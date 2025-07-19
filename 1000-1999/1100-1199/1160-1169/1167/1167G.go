package main

import (
   "bufio"
   "fmt"
   "math"
   "math/bits"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, d int
   fmt.Fscan(reader, &n, &d)
   locs := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &locs[i])
   }
   var m int
   fmt.Fscan(reader, &m)
   queries := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &queries[i])
   }

   const kMaxRange = 20000
   // prepare dot locations
   dot_locs := make([]int, 0, 2*n+2)
   for _, x := range locs {
       dot_locs = append(dot_locs, x, x+1)
   }
   dot_locs = append(dot_locs, -10000000)
   sort.Ints(dot_locs)
   // unique
   uniq := dot_locs[:1]
   for i := 1; i < len(dot_locs); i++ {
       if dot_locs[i] != uniq[len(uniq)-1] {
           uniq = append(uniq, dot_locs[i])
       }
   }
   dot_locs = uniq
   dot_locs = append(dot_locs, 2000000000)

   // initialize bitsets
   wlen := (kMaxRange + 63) / 64
   pref := make([]uint64, wlen)
   suf := make([]uint64, wlen)
   for _, p := range dot_locs {
       if p >= 0 && p < kMaxRange {
           idx := p / 64
           off := p % 64
           suf[idx] |= 1 << off
       }
   }
   lastX := 0
   dotptr := 0
   pi := math.Pi
   for _, x := range queries {
       // previous and next loc indices
       posNext := sort.SearchInts(locs, x)
       posPrev := posNext - 1
       // advance dot pointer
       for dotptr < len(dot_locs) && dot_locs[dotptr] < x {
           dotptr++
       }
       xDiff := x - lastX
       lastX = x
       if xDiff >= kMaxRange {
           for i := range pref {
               pref[i], suf[i] = 0, 0
           }
       } else {
           shiftLeft(pref, xDiff, kMaxRange)
           shiftRight(suf, xDiff)
       }
       // update pref
       for dl := dotptr - 1; dl >= 0; dl-- {
           dx := x - dot_locs[dl]
           if dx >= kMaxRange || dx > xDiff {
               break
           }
           idx := dx / 64
           off := dx % 64
           pref[idx] |= 1 << off
       }
       // update suf
       upto := x + kMaxRange
       dr := sort.Search(len(dot_locs), func(i int) bool { return dot_locs[i] > upto }) - 1
       for ; dr >= dotptr; dr-- {
           dx := dot_locs[dr] - x
           if dx < 0 || dx >= kMaxRange || dx < kMaxRange-xDiff {
               break
           }
           idx := dx / 64
           off := dx % 64
           suf[idx] |= 1 << off
       }
       // compute answer
       ans := 0.0
       if posPrev >= 0 {
           if locs[posPrev] == x-1 {
               ans = math.Max(ans, pi/2)
           } else {
               ans = math.Max(ans, math.Atan(1.0/float64(x-locs[posPrev]-1)))
           }
       }
       if posNext < n {
           if locs[posNext] == x {
               ans = math.Max(ans, pi/2)
           } else {
               ans = math.Max(ans, math.Atan(1.0/float64(locs[posNext]-x)))
           }
       }
       if posPrev < 0 || posNext >= n {
           fmt.Fprintln(writer, ans)
           continue
       }
       if locs[posPrev] == locs[posNext]-1 {
           ans = math.Max(ans, pi)
           fmt.Fprintln(writer, ans)
           continue
       }
       // intersection of pref and suf
       dist := findFirstIntersect(pref, suf)
       if dist >= 0 {
           ans = math.Max(ans, 2*math.Atan(1.0/float64(dist)))
       }
       fmt.Fprintln(writer, ans)
   }
}

// shiftLeft shifts bitset a left by k bits, clearing overflow beyond kMaxRange
func shiftLeft(a []uint64, k, kMaxRange int) {
   w := k / 64
   o := k % 64
   n := len(a)
   if w > 0 {
       for i := n - 1; i >= 0; i-- {
           if i-w >= 0 {
               a[i] = a[i-w]
           } else {
               a[i] = 0
           }
       }
   }
   if o > 0 {
       for i := n - 1; i > 0; i-- {
           a[i] = (a[i] << o) | (a[i-1] >> (64 - o))
       }
       a[0] <<= o
   }
   // clear bits beyond kMaxRange
   totalBits := n * 64
   extra := totalBits - kMaxRange
   if extra > 0 {
       mask := ^uint64(0) >> extra
       a[n-1] &= mask
   }
}

// shiftRight shifts bitset a right by k bits
func shiftRight(a []uint64, k int) {
   w := k / 64
   o := k % 64
   n := len(a)
   if w > 0 {
       for i := 0; i < n; i++ {
           if i+w < n {
               a[i] = a[i+w]
           } else {
               a[i] = 0
           }
       }
   }
   if o > 0 {
       for i := 0; i < n-1; i++ {
           a[i] = (a[i] >> o) | (a[i+1] << (64 - o))
       }
       a[n-1] >>= o
   }
}

// findFirstIntersect returns the smallest bit index where both a and b have a bit set, or -1
func findFirstIntersect(a, b []uint64) int {
   for i := range a {
       m := a[i] & b[i]
       if m != 0 {
           idx := bits.TrailingZeros64(m)
           return i*64 + idx
       }
   }
   return -1
}
