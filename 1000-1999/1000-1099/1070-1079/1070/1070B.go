package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
   "strings"
   "math/bits"
)

type interval struct { l, r uint64 }
type subnet struct { ip uint64; prefix int }

func parseSubnet(s string) (interval, error) {
   // s like a.b.c.d or a.b.c.d/x
   var maskLen int
   addr := s
   if i := strings.Index(s, "/"); i >= 0 {
       addr = s[:i]
       ml, err := strconv.Atoi(s[i+1:])
       if err != nil {
           return interval{}, err
       }
       maskLen = ml
   } else {
       maskLen = 32
   }
   parts := strings.Split(addr, ".")
   var ip uint64
   for _, p := range parts {
       v, err := strconv.Atoi(p)
       if err != nil {
           return interval{}, err
       }
       ip = (ip<<8) | uint64(v)
   }
   // compute range
   const MAX = (1<<32) - 1
   // mask bits: high maskLen bits one
   var mask uint64
   if maskLen == 0 {
       mask = 0
   } else {
       mask = MAX ^ ((uint64(1) << (32 - maskLen)) - 1)
   }
   l := ip & mask
   r := l | (MAX ^ mask)
   return interval{l, r}, nil
}

func mergeIntervals(a []interval) []interval {
   if len(a) == 0 {
       return a
   }
   sort.Slice(a, func(i, j int) bool { return a[i].l < a[j].l })
   res := make([]interval, 0, len(a))
   cur := a[0]
   for _, iv := range a[1:] {
       if iv.l <= cur.r+1 {
           if iv.r > cur.r {
               cur.r = iv.r
           }
       } else {
           res = append(res, cur)
           cur = iv
       }
   }
   res = append(res, cur)
   return res
}

func subtractIntervals(b, w []interval) []interval {
   // return b \\ w; both lists sorted, non-overlapping
   res := make([]interval, 0)
   j := 0
   for _, bi := range b {
       start := bi.l
       // advance j to first w that may overlap bi
       for j < len(w) && w[j].r < bi.l {
           j++
       }
       idx := j
       for idx < len(w) && w[idx].l <= bi.r {
           if w[idx].l > start {
               res = append(res, interval{start, w[idx].l - 1})
           }
           if w[idx].r+1 > start {
               start = w[idx].r + 1
           }
           idx++
       }
       if start <= bi.r {
           res = append(res, interval{start, bi.r})
       }
   }
   return res
}

func splitToSubnets(iv interval) []subnet {
   var out []subnet
   cur := iv.l
   for cur <= iv.r {
       // alignment via trailing zeros
       tz := bits.TrailingZeros32(uint32(cur))
       if tz > 32 {
           tz = 32
       }
       block := uint64(1) << tz
       rem := iv.r - cur + 1
       for block > rem {
           block >>= 1
           tz--
       }
       prefix := 32 - int(tz)
       out = append(out, subnet{cur, prefix})
       cur += block
   }
   return out
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   blacks := make([]interval, 0, n)
   whites := make([]interval, 0, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(in, &s)
       sign := s[0]
       sub := s[1:]
       iv, err := parseSubnet(sub)
       if err != nil {
           fmt.Fprintln(out, "-1")
           return
       }
       if sign == '-' {
           blacks = append(blacks, iv)
       } else {
           whites = append(whites, iv)
       }
   }
   // check intersection
   sort.Slice(blacks, func(i, j int) bool { return blacks[i].l < blacks[j].l })
   sort.Slice(whites, func(i, j int) bool { return whites[i].l < whites[j].l })
   i, j := 0, 0
   for i < len(blacks) && j < len(whites) {
       b := blacks[i]
       w := whites[j]
       if b.r < w.l {
           i++
       } else if w.r < b.l {
           j++
       } else {
           fmt.Fprintln(out, "-1")
           return
       }
   }
   // merge intervals
   mb := mergeIntervals(blacks)
   mw := mergeIntervals(whites)
   // subtract
   remain := subtractIntervals(mb, mw)
   // split to subnets
   var ans []subnet
   for _, iv := range remain {
       ss := splitToSubnets(iv)
       ans = append(ans, ss...)
   }
   // output
   fmt.Fprintln(out, len(ans))
   for _, s := range ans {
       a := (s.ip >> 24) & 0xFF
       b := (s.ip >> 16) & 0xFF
       c := (s.ip >> 8) & 0xFF
       d := s.ip & 0xFF
       fmt.Fprintf(out, "%d.%d.%d.%d/%d\n", a, b, c, d, s.prefix)
   }
}
