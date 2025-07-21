package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   grid := make([]byte, n*m)
   // read grid
   for i := 0; i < n; i++ {
       line, _ := reader.ReadBytes('\n')
       if len(line) < m {
           // no newline split
           // try reading full line
           line = append(line, make([]byte, m-len(line))...)
       }
       for j := 0; j < m; j++ {
           c := line[j]
           if c == '#' {
               grid[i*m+j] = 1
           }
       }
   }
   // start neighbors and end neighbors
   // zero-based: S1=(0,1), S2=(1,0), T1=(n-2,m-1), T2=(n-1,m-2)
   // compute dp from S1 and S2
   a11, a12 := computeDP(n, m, grid, 0, 1)
   a21, a22 := computeDP(n, m, grid, 1, 0)
   // determinant a11*a22 - a12*a21
   res := (int64(a11)*int64(a22) - int64(a12)*int64(a21)) % mod
   if res < 0 {
       res += mod
   }
   fmt.Println(res)
}

// computeDP returns number of ways from start (sx,sy) to T1 and T2
func computeDP(n, m int, grid []byte, sx, sy int) (t1, t2 int32) {
   // targets
   t1x, t1y := n-2, m-1
   t2x, t2y := n-1, m-2
   // if start or targets invalid or blocked, we still compute dp and get zero
   total := n * m
   dp := make([]int32, total)
   // if start is obstacle, return zeroes
   if sx < 0 || sx >= n || sy < 0 || sy >= m || grid[sx*m+sy] != 0 {
       return 0, 0
   }
   dp[sx*m+sy] = 1
   for i := sx; i < n; i++ {
       base := i * m
       for j := sy; j < m; j++ {
           idx := base + j
           if i == sx && j == sy {
               continue
           }
           if grid[idx] != 0 {
               dp[idx] = 0
           } else {
               var v int32
               if i > sx {
                   v += dp[(i-1)*m+j]
               }
               if j > sy {
                   v += dp[i*m+j-1]
               }
               if v >= mod {
                   v %= mod
               }
               dp[idx] = v
           }
       }
   }
   // fetch targets if within dp region
   if t1x >= sx && t1y >= sy {
       t1 = dp[t1x*m+t1y]
   }
   if t2x >= sx && t2y >= sy {
       t2 = dp[t2x*m+t2y]
   }
   return
}
