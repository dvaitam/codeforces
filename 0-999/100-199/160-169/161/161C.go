package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   lenK [31]int64
   memo map[string]int64
)

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}
func max(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func rec(k int, l1, r1, l2, r2 int64) int64 {
   if l1 > r1 || l2 > r2 || k <= 0 {
       return 0
   }
   key := fmt.Sprintf("%d_%d_%d_%d_%d", k, l1, r1, l2, r2)
   if v, ok := memo[key]; ok {
       return v
   }
   totalLen := lenK[k]
   // full cover
   if l1 <= 1 && r1 >= totalLen {
       return r2 - l2 + 1
   }
   if l2 <= 1 && r2 >= totalLen {
       return r1 - l1 + 1
   }
   // exact same
   if l1 == l2 && r1 == r2 {
       return r1 - l1 + 1
   }
   mid := (totalLen + 1) / 2
   var ans int64
   // cross mid at this level
   if l1 <= mid && mid <= r1 && l2 <= mid && mid <= r2 {
       left := min(mid-l1, mid-l2)
       right := min(r1-mid, r2-mid)
       cross := 1 + left + right
       ans = cross
   }
   // parts in left and right children (copies of S(k-1))
   l1l, r1l := l1, min(r1, mid-1)
   l2l, r2l := l2, min(r2, mid-1)
   l1r, r1r := max(1, l1-mid), max(int64(0), r1-mid)
   l2r, r2r := max(1, l2-mid), max(int64(0), r2-mid)
   // left1 vs left2
   if l1l <= r1l && l2l <= r2l {
       v := rec(k-1, l1l, r1l, l2l, r2l)
       if v > ans {
           ans = v
       }
   }
   // right1 vs right2
   if l1r <= r1r && l2r <= r2r {
       v := rec(k-1, l1r, r1r, l2r, r2r)
       if v > ans {
           ans = v
       }
   }
   // left1 vs right2
   if l1l <= r1l && l2r <= r2r {
       v := rec(k-1, l1l, r1l, l2r, r2r)
       if v > ans {
           ans = v
       }
   }
   // right1 vs left2
   if l1r <= r1r && l2l <= r2l {
       v := rec(k-1, l1r, r1r, l2l, r2l)
       if v > ans {
           ans = v
       }
   }
   memo[key] = ans
   return ans
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var l1, r1, l2, r2 int64
   if _, err := fmt.Fscan(reader, &l1, &r1, &l2, &r2); err != nil {
       return
   }
   // precompute lengths
   lenK[0] = 0
   lenK[1] = 1
   for i := 2; i <= 30; i++ {
       lenK[i] = lenK[i-1]*2 + 1
   }
   memo = make(map[string]int64)
   res := rec(30, l1, r1, l2, r2)
   fmt.Println(res)
}
