package main

import (
   "bufio"
   "fmt"
   "io"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   s, err := reader.ReadString('\n')
   if err != nil && err != io.EOF {
       return
   }
   s = strings.TrimSpace(s)
   m := len(s)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // dp[j][d][p] reachable: j flips used, d direction (0=+1,1=-1), p index [0..2m]
   offset := m
   maxP := 2*m + 1
   // initialize dp arrays
   dpPrev := make([][][]bool, n+1)
   dpCurr := make([][][]bool, n+1)
   for j := 0; j <= n; j++ {
       dpPrev[j] = make([][]bool, 2)
       dpCurr[j] = make([][]bool, 2)
       for d := 0; d < 2; d++ {
           dpPrev[j][d] = make([]bool, maxP)
           dpCurr[j][d] = make([]bool, maxP)
       }
   }
   // start at position 0, direction +
   dpPrev[0][0][offset] = true
   // process each command
   for i := 0; i < m; i++ {
       // clear dpCurr
       for j := 0; j <= n; j++ {
           for d := 0; d < 2; d++ {
               for p := 0; p < maxP; p++ {
                   dpCurr[j][d][p] = false
               }
           }
       }
       // original command
       orig := s[i]
       for j := 0; j <= n; j++ {
           for d := 0; d < 2; d++ {
               for p := 0; p < maxP; p++ {
                   if !dpPrev[j][d][p] {
                       continue
                   }
                   // two options: no flip or flip (if j+1 <= n)
                   for flip := 0; flip < 2; flip++ {
                       nj := j + flip
                       if nj > n {
                           continue
                       }
                       // resulting command
                       cmd := orig
                       if flip == 1 {
                           if orig == 'T' {
                               cmd = 'F'
                           } else {
                               cmd = 'T'
                           }
                       }
                       nd, np := d, p
                       if cmd == 'T' {
                           // turn
                           nd = 1 - d
                       } else {
                           // forward
                           if d == 0 {
                               np = p + 1
                           } else {
                               np = p - 1
                           }
                       }
                       // check bounds
                       if np >= 0 && np < maxP {
                           dpCurr[nj][nd][np] = true
                       }
                   }
               }
           }
       }
       // swap dpPrev and dpCurr
       dpPrev, dpCurr = dpCurr, dpPrev
   }
   // compute answer: max |pos| reachable with exactly n flips (allow waste flips parity)
   ans := 0
   for j := 0; j <= n; j++ {
       if (n-j)%2 != 0 {
           continue
       }
       for d := 0; d < 2; d++ {
           for p := 0; p < maxP; p++ {
               if dpPrev[j][d][p] {
                   pos := p - offset
                   if pos < 0 {
                       pos = -pos
                   }
                   if pos > ans {
                       ans = pos
                   }
               }
           }
       }
   }
   fmt.Println(ans)
}
