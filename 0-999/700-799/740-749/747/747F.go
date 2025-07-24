package main

import (
   "bufio"
   "fmt"
   "os"
)

// find k-th smallest positive integer in hexadecimal notation with each digit <= t times
var (
   t int
   memo map[state]uint64
)

type state struct {
   rem  uint8
   defs [16]uint8
}

// count number of sequences of length rem given current deficits defs (counts used)
func count(rem int, defs [16]uint8) uint64 {
   if rem == 0 {
       return 1
   }
   st := state{uint8(rem), defs}
   if v, ok := memo[st]; ok {
       return v
   }
   var total uint64
   for d := 0; d < 16; d++ {
       if int(defs[d]) < t {
           defs[d]++
           total += count(rem-1, defs)
           defs[d]--
       }
   }
   memo[st] = total
   return total
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var k uint64
   if _, err := fmt.Fscan(in, &k, &t); err != nil {
       return
   }
   memo = make(map[state]uint64)
   // determine length
   var length int
   // max possible length is 16*t, but will break earlier
   for l := 1; l <= 16*t; l++ {
       var cnt uint64
       // first digit cannot be zero
       for d := 1; d < 16; d++ {
           if t > 0 {
               var defs [16]uint8
               defs[d] = 1
               cnt += count(l-1, defs)
           }
       }
       if cnt >= k {
           length = l
           break
       }
       k -= cnt
   }
   // build result
   var defs [16]uint8
   rem := length
   res := make([]byte, 0, length)
   for pos := 0; pos < length; pos++ {
       for d := 0; d < 16; d++ {
           if pos == 0 && d == 0 {
               continue
           }
           if int(defs[d]) >= t {
               continue
           }
           defs[d]++
           cnt := count(rem-1, defs)
           if cnt >= k {
               // choose d
               var ch byte
               if d < 10 {
                   ch = byte('0' + d)
               } else {
                   ch = byte('a' + (d - 10))
               }
               res = append(res, ch)
               rem--
               break
           } else {
               k -= cnt
               defs[d]--
           }
       }
   }
   // output
   fmt.Fprint(os.Stdout, string(res))
}
