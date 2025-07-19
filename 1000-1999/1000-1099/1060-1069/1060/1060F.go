package main

import (
   "bufio"
   "fmt"
   "os"
)

const maxN = 55

var n int
var c [maxN][maxN]float64
var e [][]int
var vis []bool
var fa []int
var siz []int
var f [][]float64
var f1 []float64

func dfs(x int) {
   vis[x] = true
   siz[x] = 0
   // traverse children
   for _, y := range e[x] {
       if !vis[y] {
           fa[y] = x
           dfs(y)
           siz[x] += siz[y] + 1
       }
   }
   // reset f[x]
   for i := 0; i <= n; i++ {
       f[x][i] = 0
   }
   pre := 0
   f[x][0] = 1.0
   // combine DP from children
   for _, y := range e[x] {
       if fa[y] != x {
           continue
       }
       // clear f1 up to new size
       limit := pre + siz[y] + 1
       for i := 0; i <= limit; i++ {
           f1[i] = 0
       }
       // merge states
       for px := 1; px <= pre; px++ {
           for py := 1; py <= siz[y]+1; py++ {
               sumIdx := pre + siz[y] + 1
               // first w loop
               for w := px; w <= siz[y]+1+px; w++ {
                   if px-1+py-1 >= w-1 {
                       f1[w] += f[x][px] * f[y][py] * c[w-1][px-1] * c[sumIdx-w][pre-px]
                   }
               }
               // second w loop
               for w := py; w <= pre+py; w++ {
                   if px-1+py-1 >= w-1 {
                       f1[w] += f[x][px] * f[y][py] * c[w-1][py-1] * c[sumIdx-w][siz[y]+1-py]
                   }
               }
           }
       }
       // case f[y][0]
       f1[0] += f[x][0] * f[y][0] * c[pre+siz[y]+1][pre]
       for px := 1; px <= pre; px++ {
           for w := px; w <= siz[y]+1+px; w++ {
               f1[w] += f[x][px] * f[y][0] * c[w-1][px-1] * c[pre+siz[y]+1-w][pre-px]
           }
       }
       for py := 1; py <= siz[y]+1; py++ {
           for w := py; w <= pre+py; w++ {
               f1[w] += f[x][0] * f[y][py] * c[w-1][py-1] * c[pre+siz[y]+1-w][siz[y]+1-py]
           }
       }
       // copy back
       for i := 0; i <= limit; i++ {
           f[x][i] = f1[i]
       }
       pre += siz[y] + 1
   }
   // final accumulation
   for i := 0; i <= pre+1; i++ {
       f1[i] = 0
   }
   for px := 1; px <= pre; px++ {
       f1[px+1] += f[x][px] * float64(px)
       for w := 1; w <= px; w++ {
           f1[w] -= f[x][px] * 0.5
       }
   }
   f1[0] += f[x][0] * float64(pre+1)
   for w := 1; w <= pre+1; w++ {
       f1[w] -= f[x][0] * 0.5
   }
   for i := 0; i <= pre+1; i++ {
       f[x][i] = f1[i]
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n)
   // init combinatorics
   for i := 0; i <= n; i++ {
       c[i][0] = 1
       for j := 1; j <= i; j++ {
           c[i][j] = c[i-1][j-1] + c[i-1][j]
       }
   }
   // init globals
   e = make([][]int, n+1)
   vis = make([]bool, n+1)
   fa = make([]int, n+1)
   siz = make([]int, n+1)
   f = make([][]float64, n+1)
   for i := 0; i <= n; i++ {
       f[i] = make([]float64, n+2)
   }
   f1 = make([]float64, n+2)
   // read tree
   for i := 2; i <= n; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       e[x] = append(e[x], y)
       e[y] = append(e[y], x)
   }
   // for each root
   for st := 1; st <= n; st++ {
       for i := 1; i <= n; i++ {
           vis[i] = false
           fa[i] = 0
       }
       dfs(st)
       ans := 1.0
       // compute answer
       for _, x := range e[st] {
           // x is child since vis reset
           tmp := 0.0
           for i := 0; i <= siz[x]+1; i++ {
               tmp += f[x][i]
           }
           for i := 1; i <= siz[x]+1; i++ {
               tmp /= float64(i)
           }
           ans *= tmp
       }
       fmt.Fprintf(writer, "%.9f\n", ans)
   }
}
