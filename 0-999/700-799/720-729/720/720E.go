package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   // precompute powers of 10 up to 18
   const MAXN = 18
   pow10 := make([]uint64, MAXN+1)
   pow10[0] = 1
   for i := 1; i <= MAXN; i++ {
       pow10[i] = pow10[i-1] * 10
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(in, &n)
       var s string
       fmt.Fscan(in, &s)
       // read mapping, lines of 10 chars
       codes := make([]string, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &codes[i])
       }
       // parse N0 value (may have leading zeros)
       var N0 uint64
       for i := 0; i < n; i++ {
           N0 = N0*10 + uint64(s[i]-'0')
       }
       // wrap time: when beep sounds
       wrap := pow10[n] - N0
       // compute max code distinguish time
       var tcode uint64
       infinite := false
       for i := 0; i < n && !infinite; i++ {
           // initial digit
           d0 := int(s[i] - '0')
           ci := codes[i]
           // find all digits with same initial code
           var same []int
           for d := 0; d < 10; d++ {
               if ci[d] == ci[d0] {
                   same = append(same, d)
               }
           }
           if len(same) <= 1 {
               continue
           }
           // position weight: seconds between digit increments
           delta := pow10[n-1-i]
           // track max s for this position
           var smax uint64
           for _, d1 := range same {
               if d1 == d0 {
                   continue
               }
               // find minimal s>=1 s.t. codes differ, or mark identical
               found := false
               for sstep := 1; sstep < 10; sstep++ {
                   a := (d0 + sstep) % 10
                   b := (d1 + sstep) % 10
                   if ci[a] != ci[b] {
                       // distinguished at sstep
                       if uint64(sstep) > smax {
                           smax = uint64(sstep)
                       }
                       found = true
                       break
                   }
               }
               if !found {
                   // sequences identical, rely on beep
                   infinite = true
                   break
               }
           }
           if infinite {
               break
           }
           // time until this position distinguishes worst-case
           ti := smax * delta
           if ti > tcode {
               tcode = ti
           }
       }
       var ans uint64
       if infinite || tcode >= wrap {
           ans = wrap
       } else {
           ans = tcode
       }
       fmt.Fprintln(out, ans)
   }
}
