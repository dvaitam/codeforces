package main

import (
   "bufio"
   "fmt"
   "io"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil && err != io.EOF {
       panic(err)
   }
   // dp[r]: max damage with total cards played %10 == r
   const INF = int64(9e18)
   dp := make([]int64, 10)
   for i := 1; i < 10; i++ {
       dp[i] = -INF
   }
   for t := 0; t < n; t++ {
       var k int
       fmt.Fscan(reader, &k)
       // track top3 for cost1, top1 for cost2 and cost3
       var d1 [3]int64
       var d2, d3 int64
       for i := 0; i < 3; i++ {
           d1[i] = -1
       }
       d2, d3 = -1, -1
       for i := 0; i < k; i++ {
           var c int
           var d int64
           fmt.Fscan(reader, &c, &d)
           switch c {
           case 1:
               // insert d into d1 sorted desc
               if d > d1[0] {
                   d1[2] = d1[1]
                   d1[1] = d1[0]
                   d1[0] = d
               } else if d > d1[1] {
                   d1[2] = d1[1]
                   d1[1] = d
               } else if d > d1[2] {
                   d1[2] = d
               }
           case 2:
               if d > d2 {
                   d2 = d
               }
           case 3:
               if d > d3 {
                   d3 = d
               }
           }
       }
       // generate options: cnt, base sum, max damage in selection
       type opt struct{ cnt int; base, maxd int64 }
       opts := make([]opt, 0, 7)
       opts = append(opts, opt{0, 0, 0})
       // cost=1 options
       if d1[0] >= 0 {
           opts = append(opts, opt{1, d1[0], d1[0]})
       }
       if d2 >= 0 {
           opts = append(opts, opt{1, d2, d2})
       }
       if d3 >= 0 {
           opts = append(opts, opt{1, d3, d3})
       }
       // cost=2
       if d1[1] >= 0 {
           opts = append(opts, opt{2, d1[0] + d1[1], d1[0]})
       }
       // cost=3: c2 + c1
       if d2 >= 0 && d1[0] >= 0 {
           maxd := d1[0]
           if d2 > maxd {
               maxd = d2
           }
           opts = append(opts, opt{2, d1[0] + d2, maxd})
       }
       // cost=3: three c1
       if d1[2] >= 0 {
           opts = append(opts, opt{3, d1[0] + d1[1] + d1[2], d1[0]})
       }

       // new dp
       newdp := make([]int64, 10)
       for i := 0; i < 10; i++ {
           newdp[i] = -INF
       }
       for r := 0; r < 10; r++ {
           prev := dp[r]
           if prev < -INF/2 {
               continue
           }
           for _, o := range opts {
               nr := (r + o.cnt) % 10
               bonus := int64(0)
               // position where double occurs: pos == (10 - r%10)%10
               off := (10 - r%10) % 10
               if off > 0 && off <= o.cnt {
                   bonus = o.maxd
               }
               val := prev + o.base + bonus
               if val > newdp[nr] {
                   newdp[nr] = val
               }
           }
       }
       dp = newdp
   }
   // answer is max dp[r]
   ans := int64(0)
   for _, v := range dp {
       if v > ans {
           ans = v
       }
   }
   fmt.Println(ans)
}
