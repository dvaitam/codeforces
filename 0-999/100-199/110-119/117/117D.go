package main

import (
   "bufio"
   "fmt"
   "os"
)

var mod uint64
var lq, rq, uq, vq uint64

// dfs computes the sum of values in the transformed array b
// for the segment of the subtree of length len, starting at position pos,
// with values given by arithmetic progression starting at val with step spacing.
func dfs(len, pos, val, step uint64) uint64 {
   // no overlap on positions
   if pos > rq || pos+len-1 < lq {
       return 0
   }
   // compute last value in this subtree
   span := step * (len - 1)
   valEnd := val + span
   // no overlap on values
   if val > vq || valEnd < uq {
       return 0
   }
   // fully covered segment and value range
   if pos >= lq && pos+len-1 <= rq && val >= uq && valEnd <= vq {
       return sumAP(len, val, valEnd)
   }
   // single element case
   if len == 1 {
       // here pos and val overlap query ranges
       return val % mod
   }
   // split into odd and even subtrees
   lenOdd := (len + 1) / 2
   // odd positions: values val, val+2*step, ...
   sum := dfs(lenOdd, pos, val, step*2)
   // even positions: values val+step, val+3*step, ...
   lenEven := len / 2
   sum += dfs(lenEven, pos+lenOdd, val+step, step*2)
   return sum % mod
}

// sumAP returns sum of arithmetic progression first, first+step*(len-1)
// here given first and last, sum = len*(first+last)/2 mod mod
func sumAP(len, first, last uint64) uint64 {
   if len%2 == 0 {
       half := (len / 2) % mod
       return mulMod(half, (first+last)%mod)
   }
   halfSum := ((first + last) / 2) % mod
   return mulMod(len%mod, halfSum)
}

// mulMod performs (a*b) % mod safely assuming mod fits in uint64
func mulMod(a, b uint64) uint64 {
   return (a * b) % mod
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m uint64
   if _, err := fmt.Fscan(reader, &n, &m, &mod); err != nil {
       return
   }
   for i := uint64(0); i < m; i++ {
       fmt.Fscan(reader, &lq, &rq, &uq, &vq)
       res := dfs(n, 1, 1, 1)
       fmt.Fprintln(writer, res)
   }
}
