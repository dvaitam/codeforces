package main

import (
   "bufio"
   "fmt"
   "os"
)

const (
   LOG = 16
   M   = 50004
   N   = 100010
)

var (
   sparse  [LOG][M]int
   a       [N]int64
   pref    [N]int64
   prefSum [N]int64
   n       int
)

func gcd(x, y int) int {
   if y == 0 {
       return x
   }
   return gcd(y, x%y)
}

func solveTriang(A, B, C int64) int64 {
   if C < 0 {
       return 0
   }
   if A > B {
       A, B = B, A
   }
   p := C / B
   k := B / A
   d := (C - p*B) / A
   return solveTriang(B-k*A, A, C-A*(k*p+d+1)) + (p+1)*(d+1) + k*p*(p+1)/2
}

func solve2(A, CA, B, CB, C int64) int64 {
   if C >= A*CA+B*CB {
       return CA * CB
   }
   if C < 0 {
       return 0
   }
   return solveTriang(A, B, C) - solveTriang(A, B, C-A*CA) - solveTriang(A, B, C-B*CB) + solveTriang(A, B, C-A*CA-B*CB)
}

func solve(B int64) int64 {
   var ans int64
   // single index contributions
   for i := 1; i < n; i++ {
       if a[i] > 0 {
           L := a[i]
           maxL := B / int64(i)
           if L > maxL {
               L = maxL
           }
           ans += L*(a[i]+1) - L*(L+1)/2
       }
   }
   r := 1
   for l := 1; l < n; l++ {
       if a[l] == 0 {
           continue
       }
       if r > l {
           ans += a[l] * (pref[r-1] - pref[l])
       }
       for r < n && prefSum[r]-prefSum[l] <= B {
           if r > l && a[r] > 0 {
               ans += solve2(int64(l), a[l], int64(r), a[r], B-(prefSum[r-1]-prefSum[l])-int64(l)-int64(r))
           }
           r++
       }
       if r == n {
           continue
       }
       ans += solve2(int64(l), a[l], int64(r), a[r], B-(prefSum[r-1]-prefSum[l])-int64(l)-int64(r))
   }
   return ans
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   fmt.Fscan(reader, &n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &sparse[0][i])
   }
   // build sparse table
   for k := 0; k < LOG-1; k++ {
       for i := 0; i < n; i++ {
           v := sparse[k][i]
           j := i + (1 << k)
           if j < n {
               v = gcd(v, sparse[k][j])
           }
           sparse[k+1][i] = v
       }
   }
   // count subarrays by gcd
   for l := 0; l < n; l++ {
       cur := 0
       r := l
       for r < n {
           cur = gcd(cur, sparse[0][r])
           rr := r
           for k := LOG - 1; k >= 0; k-- {
               if rr < n && sparse[k][rr]%cur == 0 {
                   rr += 1 << k
               }
           }
           if rr > n {
               rr = n
           }
           a[cur] += int64(rr - r)
           r = rr
       }
   }
   // total count m
   in := int64(n)
   m := in * (in + 1) / 2
   m = m * (m + 1) / 2
   m = (m + 1) / 2
   // prepare prefix sums over gcd values
   n = N
   for i := 1; i < n; i++ {
       pref[i] = pref[i-1] + a[i]
       prefSum[i] = prefSum[i-1] + a[i]*int64(i)
   }
   // binary search answer
   var l int64 = -1
   var r64 = prefSum[n-1]
   for r64-l > 1 {
       mid := (l + r64) / 2
       if solve(mid) >= m {
           r64 = mid
       } else {
           l = mid
       }
   }
   fmt.Println(r64)
}
