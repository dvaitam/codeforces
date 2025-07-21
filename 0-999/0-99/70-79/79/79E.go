package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func sumRange(L, R int64) int64 {
   if L > R {
       return 0
   }
   m := R - L + 1
   return (L + R) * m / 2
}

// sum of |k - u| for k in [L..R]
func sumAbsRange(L, R, u int64) int64 {
   if L > R {
       return 0
   }
   if u < L {
       return sumRange(L, R) - u*(R-L+1)
   }
   if u > R {
       return u*(R-L+1) - sumRange(L, R)
   }
   // L <= u <= R
   // left part [L..u-1]
   m1 := u - L
   sum1 := 0
   if m1 > 0 {
       sum1 = m1*u - sumRange(L, u-1)
   }
   // right part [u+1..R]
   m2 := R - u
   sum2 := 0
   if m2 > 0 {
       sum2 = sumRange(u+1, R) - m2*u
   }
   return sum1 + sum2
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, a, b, c int
   var t int64
   if _, err := fmt.Fscan(in, &n, &t, &a, &b, &c); err != nil {
       return
   }
   // Option A: R*(a-1), U*(n-1), R*(n-a)
   // compute sum_x for u = a and a+c-1
   nn := int64(n)
   a64 := int64(a)
   b64 := int64(b)
   c64 := int64(c)
   // sum_x_A(u)
   sumxA := func(u int64) int64 {
       // part1: k from 2 to a-1
       s1 := sumAbsRange(2, a64-1, u)
       // part2: n * |a - u|
       s2 := nn * abs(a64-u)
       // part3: k from a+1 to n
       s3 := sumAbsRange(a64+1, nn, u)
       return s1 + s2 + s3
   }
   // sum_y_A(v)
   sumyA := func(v int64) int64 {
       // y=1 for a-1 times
       s1 := int64(a-1) * abs(1-v)
       // k from 2 to n-1
       s2 := sumAbsRange(2, nn-1, v)
       // y=n for n-a+1 times
       s3 := int64(n-a+1) * abs(nn-v)
       return s1 + s2 + s3
   }
   // Option B: U*(b-1), R*(n-1), U*(n-b)
   // sum_x_B(u)
   sumxB := func(u int64) int64 {
       // x=1 for b-1 times
       s1 := int64(b-1) * abs(1-u)
       // k from 2 to n-1
       s2 := sumAbsRange(2, nn-1, u)
       // x=n for n-b+1 times
       s3 := int64(n-b+1) * abs(nn-u)
       return s1 + s2 + s3
   }
   // sum_y_B(v)
   sumyB := func(v int64) int64 {
       // k from 2 to b
       s1 := sumAbsRange(2, b64, v)
       // (n-1) * |b - v|
       s2 := nn-1 * abs(b64-v)
       // k from b+1 to n
       s3 := sumAbsRange(b64+1, nn, v)
       return s1 + s2 + s3
   }
   var okA, okB bool
   // endpoints for sensors
   u1 := a64
   u2 := a64 + c64 - 1
   v1 := b64
   v2 := b64 + c64 - 1
   // check A
   maxSxA := sumxA(u1)
   if v := sumxA(u2); v > maxSxA {
       maxSxA = v
   }
   maxSyA := sumyA(v1)
   if v := sumyA(v2); v > maxSyA {
       maxSyA = v
   }
   if maxSxA+maxSyA <= t {
       okA = true
   }
   // check B
   maxSxB := sumxB(u1)
   if v := sumxB(u2); v > maxSxB {
       maxSxB = v
   }
   maxSyB := sumyB(v1)
   if v := sumyB(v2); v > maxSyB {
       maxSyB = v
   }
   if maxSxB+maxSyB <= t {
       okB = true
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   if !okA && !okB {
       fmt.Fprintln(writer, "Impossible")
       return
   }
   // build paths
   var pathA, pathB string
   if okA {
       var sb strings.Builder
       sb.Grow(2*n - 2)
       if a > 1 {
           sb.WriteString(strings.Repeat("R", a-1))
       }
       sb.WriteString(strings.Repeat("U", n-1))
       if n-a > 0 {
           sb.WriteString(strings.Repeat("R", n-a))
       }
       pathA = sb.String()
   }
   if okB {
       var sb strings.Builder
       sb.Grow(2*n - 2)
       if b > 1 {
           sb.WriteString(strings.Repeat("U", b-1))
       }
       sb.WriteString(strings.Repeat("R", n-1))
       if n-b > 0 {
           sb.WriteString(strings.Repeat("U", n-b))
       }
       pathB = sb.String()
   }
   // choose lexicographically smallest
   if okA && okB {
       if pathA <= pathB {
           fmt.Fprintln(writer, pathA)
       } else {
           fmt.Fprintln(writer, pathB)
       }
   } else if okA {
       fmt.Fprintln(writer, pathA)
   } else {
       fmt.Fprintln(writer, pathB)
   }
}

func abs(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}
