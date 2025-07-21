package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   fmt.Fscan(in, &n, &k)
   sBytes := make([]byte, n)
   fmt.Fscan(in, &sBytes)
   s := string(sBytes)
   if k > n {
       fmt.Println(0)
       return
   }
   // dp0: state before B-run; dp1: after B-run before W-run; dp2: done
   // use ring buffers for dp0 and dp1 of size k
   dp0 := make([]int, k)
   t0 := make([]int, k)
   dp1 := make([]int, k)
   t1 := make([]int, k)
   head0, head1 := 0, 0
   ts0, ts1 := 1, 1
   // initial dp0[0]=1
   dp0[0], t0[0] = 1, ts0
   sum0 := 1
   sum1 := 0
   dp2 := 0
   for i := 0; i < n; i++ {
       c := s[i]
       allowB := c != 'W'
       allowW := c != 'B'
       // fetch q0_last and q1_last old
       idx0 := (head0 + k - 1) % k
       q0 := 0
       if t0[idx0] == ts0 {
           q0 = dp0[idx0]
       }
       idx1 := (head1 + k - 1) % k
       q1 := 0
       if t1[idx1] == ts1 {
           q1 = dp1[idx1]
       }
       // dp2 update
       opts := 0
       if allowB {
           opts++
       }
       if allowW {
           opts++
       }
       if opts == 2 {
           dp2 = (dp2 * 2) % MOD
       }
       // opts can be 1 or 0; dp2*1 or dp2*0
       if opts == 1 {
           // multiply by 1, no change
       }
       if opts == 0 {
           dp2 = 0
       }
       if allowW {
           dp2 = (dp2 + q1) % MOD
       }
       // update dp1
       oldSum1 := sum1
       if allowW {
           // shift
           head1 = (head1 + 1) % k
           // write new dp1[0]
           v1 := 0
           if allowB {
               v1 = oldSum1
           }
           dp1[head1] = v1
           t1[head1] = ts1
           // update sum1
           // sum1 = (allowB?oldSum1:0) + (oldSum1 - q1)
           s1 := 0
           if allowB {
               s1 = oldSum1
           }
           s1 += (oldSum1 - q1 + MOD)
           sum1 = s1 % MOD
       } else {
           // wipe
           ts1++
           head1 = head1 // unchanged
           v1 := 0
           if allowB {
               v1 = oldSum1
           }
           dp1[head1] = v1
           t1[head1] = ts1
           sum1 = v1 % MOD
       }
       // update dp0
       oldSum0 := sum0
       if allowB {
           // shift
           head0 = (head0 + 1) % k
           v0 := 0
           if allowW {
               v0 = oldSum0
           }
           dp0[head0] = v0
           t0[head0] = ts0
           // sum0 = (allowW?oldSum0:0) + (oldSum0 - q0)
           s0 := 0
           if allowW {
               s0 = oldSum0
           }
           s0 += (oldSum0 - q0 + MOD)
           sum0 = s0 % MOD
       } else {
           // wipe
           ts0++
           head0 = head0
           v0 := 0
           if allowW {
               v0 = oldSum0
           }
           dp0[head0] = v0
           t0[head0] = ts0
           sum0 = v0 % MOD
       }
   }
   fmt.Println(dp2)
}
