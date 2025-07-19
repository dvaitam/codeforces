package main

import (
   "bufio"
   "os"
   "sort"
)

const (
   MAXN   = 300005
   MAXLOG = 19
)

var LG [MAXN]int
var up [MAXLOG+1][MAXN]int

func init() {
   LG[1] = 0
   for i := 2; i < MAXN; i++ {
       LG[i] = LG[i>>1] + 1
   }
}

func readInt(rd *bufio.Reader) int {
   sign, x := 1, 0
   ch, err := rd.ReadByte()
   if err != nil {
       return 0
   }
   for (ch < '0' || ch > '9') && ch != '-' {
       ch, err = rd.ReadByte()
       if err != nil {
           return 0
       }
   }
   if ch == '-' {
       sign = -1
       ch, _ = rd.ReadByte()
   }
   for ch >= '0' && ch <= '9' {
       x = x*10 + int(ch-'0')
       ch, err = rd.ReadByte()
       if err != nil {
           break
       }
   }
   return x * sign
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

// ask maximum of up[...][*] in range [l, r]
func ask(l, r int) int {
   if l > r {
       return 0
   }
   k := LG[r-l+1]
   a := up[k][l]
   b := up[k][r-(1<<k)+1]
   if a > b {
       return a
   }
   return b
}

func main() {
   rd := bufio.NewReader(os.Stdin)
   wr := bufio.NewWriter(os.Stdout)
   defer wr.Flush()

   t := readInt(rd)
   for tt := 0; tt < t; tt++ {
       n := readInt(rd)
       a := make([]int, n+1)
       for i := 1; i <= n; i++ {
           a[i] = readInt(rd)
       }
       // build sparse table on a[i]+i
       for i := 1; i <= n; i++ {
           up[0][i] = a[i] + i
       }
       for k := 1; k <= MAXLOG; k++ {
           shift := 1 << (k - 1)
           for i := 1; i+ (1<<k) -1 <= n; i++ {
               x := up[k-1][i]
               y := up[k-1][i+shift]
               if y > x {
                   x = y
               }
               up[k][i] = x
           }
       }
       f := make([]int, n+1)
       pre := make([]int, n+1)
       ans := make([]byte, n+1)
       for i := 1; i <= n; i++ {
           ans[i] = 'R'
       }
       // dp
       for i := 1; i <= n; i++ {
           f[i] = f[i-1]
           pre[i] = i - 1
           if f[i] >= i {
               v := a[i] + i
               if v > f[i] {
                   f[i] = v
                   pre[i] = i - 1
               }
           }
           // binary search for first it in [0, i) with f[it] >= i - a[i] - 1
           target := i - a[i] - 1
           it := sort.Search(i, func(j int) bool { return f[j] >= target })
           if it < i {
               val := i - 1
               if f[it] > val {
                   val = f[it]
               }
               if it+1 <= i-1 {
                   t2 := ask(it+1, i-1)
                   if t2 > val {
                       val = t2
                   }
               }
               if val > f[i] {
                   f[i] = val
                   pre[i] = -it
               }
           }
       }
       if f[n] < n {
           wr.WriteString("NO\n")
           continue
       }
       wr.WriteString("YES\n")
       // reconstruct
       x := n
       for x > 0 {
           if pre[x] <= 0 {
               ans[x] = 'L'
               x = -pre[x]
           } else {
               x = pre[x]
           }
       }
       wr.Write(ans[1:])
       wr.WriteByte('\n')
   }
}
