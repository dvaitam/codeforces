package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, l, r, s, n1, n2 int64
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

// check1_t checks feasibility for given t, sum sVal, and flag
func check1_t(t, sVal, flag int64) bool {
   denom := t + n
   temp := sVal / denom
   temp2 := sVal % denom
   if temp2 >= n1+flag && temp2 <= n1+min(t, n1) {
       rem := t - (temp2 - n1)
       if rem >= 0 && rem <= n2 {
           return true
       }
   }
   if temp != 0 {
       temp2 += denom
       if temp2 >= n1+flag && temp2 <= n1+min(t, n1) {
           rem := t - (temp2 - n1)
           if rem >= 0 && rem <= n2 {
               return true
           }
       }
   }
   return false
}

// check1 tries both s and s+1 with flags 0 and 1
func check1(t int64) bool {
   return check1_t(t, s, 0) || check1_t(t, s+1, 1)
}

// solve1 brute forces t from n down to 0
func solve1() int64 {
   for x := n; x >= 0; x-- {
       if check1(x) {
           return x
       }
   }
   return -1
}

// div1 is floor division for signed ints
func div1(x, y int64) int64 {
   if x == 0 {
       return 0
   }
   if x > 0 {
       return x / y
   }
   x = -x
   return -((x-1)/y + 1)
}

// div2 is ceil division for signed ints
func div2(x, y int64) int64 {
   if x == 0 {
       return 0
   }
   if x > 0 {
       return (x-1)/y + 1
   }
   x = -x
   return -(x / y)
}

// calc2 computes maximum t for given k, sum sVal, and flag
func calc2(k, sVal, flag int64) int64 {
   sTemp := sVal - k*n - n1
   l1 := div2(sTemp, k+1)
   l2 := div2(sTemp-n1, k)
   l3 := div1(sTemp-flag, k)
   l4 := div1(sTemp+n2, k+1)
   l1 = max(l1, l2)
   l3 = min(l3, l4)
   if l1 <= l3 {
       return l3
   }
   return -1
}

// solve2 uses analytical approach
func solve2() int64 {
   var ans int64 = -1
   if s >= n1 && s <= n1*2 {
       ans = max(ans, s-n1+n2)
   }
   if s+1 > n1 && s+1 <= n1*2 {
       ans = max(ans, s+1-n1+n2)
   }
   for k := int64(1); k*n <= s+1; k++ {
       ans = max(ans, calc2(k, s, 0))
       ans = max(ans, calc2(k, s+1, 1))
   }
   return ans
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n, &l, &r, &s)
   n1 = r - l + 1
   if n1 <= 0 {
       n1 += n
   }
   n2 = n - n1
   var res int64
   if n < s/n {
       res = solve1()
   } else {
       res = solve2()
   }
   fmt.Fprintln(writer, res)
}
