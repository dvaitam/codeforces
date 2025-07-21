package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

type MaxCnt struct {
   val int
   cnt int64
}

func maxCnt(a MaxCnt, b MaxCnt) MaxCnt {
   if a.val > b.val {
       return a
   } else if b.val > a.val {
       return b
   }
   // equal vals
   return MaxCnt{a.val, a.cnt + b.cnt}
}

func maxCnt3(a, b, c MaxCnt) MaxCnt {
   // merge a and b, then merge with c, taking care of double-count for overlap of a and b not needed as a and b from disjoint sources
   m := maxCnt(a, b)
   return maxCnt(m, c)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, r int
   if _, err := fmt.Fscan(in, &n, &m, &r); err != nil {
       return
   }
   a := make([][]int, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(in, &a[i][j])
       }
   }
   H := n - 2*r
   W := m - 2*r
   if H <= 0 || W <= 0 {
       fmt.Println("0 0")
       return
   }
   // row prefix sums
   rp := make([][]int, n)
   for i := 0; i < n; i++ {
       rp[i] = make([]int, m+1)
       for j := 0; j < m; j++ {
           rp[i][j+1] = rp[i][j] + a[i][j]
       }
   }
   // widths for circle
   wds := make([]int, 2*r+1)
   for di := -r; di <= r; di++ {
       wds[di+r] = int(math.Floor(math.Sqrt(float64(r*r - di*di))))
   }
   // compute weights
   wgrid := make([][]int, H)
   for ci := 0; ci < H; ci++ {
       wgrid[ci] = make([]int, W)
       for cj := 0; cj < W; cj++ {
           sum := 0
           for di := -r; di <= r; di++ {
               wdx := wds[di+r]
               row := ci + r + di
               left := cj + r - wdx
               right := cj + r + wdx
               sum += rp[row][right+1] - rp[row][left]
           }
           wgrid[ci][cj] = sum
       }
   }
   // case1: row separation
   rowMax := 0
   var rowCnt int64 = 0
   // best_up for rows
   bestUp := make([]MaxCnt, H)
   // compute max in row
   for i := 0; i < H; i++ {
       // max in this row
       mval := 0
       var mcnt int64 = 0
       for j := 0; j < W; j++ {
           v := wgrid[i][j]
           if v > mval {
               mval = v; mcnt = 1
           } else if v == mval {
               mcnt++
           }
       }
       cur := MaxCnt{mval, mcnt}
       if i == 0 {
           bestUp[i] = cur
       } else {
           bestUp[i] = maxCnt(bestUp[i-1], cur)
       }
   }
   sep := 2*r + 1
   for i2 := 0; i2 < H; i2++ {
       i1lim := i2 - sep
       if i1lim < 0 {
           continue
       }
       best1 := bestUp[i1lim]
       if best1.cnt == 0 {
           continue
       }
       for j2 := 0; j2 < W; j2++ {
           sum := wgrid[i2][j2] + best1.val
           if sum > rowMax {
               rowMax = sum; rowCnt = best1.cnt
           } else if sum == rowMax {
               rowCnt += best1.cnt
           }
       }
   }
   // case2: column separation
   colMax := 0
   var colCnt int64 = 0
   bestL := make([]MaxCnt, W)
   for j := 0; j < W; j++ {
       mval := 0; var mcnt int64 = 0
       for i := 0; i < H; i++ {
           v := wgrid[i][j]
           if v > mval { mval = v; mcnt = 1 } else if v == mval { mcnt++ }
       }
       cur := MaxCnt{mval, mcnt}
       if j == 0 {
           bestL[j] = cur
       } else {
           bestL[j] = maxCnt(bestL[j-1], cur)
       }
   }
   for j2 := 0; j2 < W; j2++ {
       j1lim := j2 - sep
       if j1lim < 0 {
           continue
       }
       best1 := bestL[j1lim]
       if best1.cnt == 0 {
           continue
       }
       for i2 := 0; i2 < H; i2++ {
           sum := wgrid[i2][j2] + best1.val
           if sum > colMax {
               colMax = sum; colCnt = best1.cnt
           } else if sum == colMax {
               colCnt += best1.cnt
           }
       }
   }
   // case3: top-left to bottom-right
   diag1Max := 0; var diag1Cnt int64 = 0
   dpTL := make([][]MaxCnt, H)
   for i := 0; i < H; i++ {
       dpTL[i] = make([]MaxCnt, W)
       for j := 0; j < W; j++ {
           cur := MaxCnt{wgrid[i][j], 1}
           if i > 0 {
               cur = mergeRect(cur, dpTL[i-1][j])
           }
           if j > 0 {
               cur = mergeRect(cur, dpTL[i][j-1])
           }
           if i > 0 && j > 0 {
               // subtract overlap if needed
               // if dpTL[i-1][j].val==dpTL[i][j-1].val and ==cur.val, subtract dpTL[i-1][j-1].cnt
               if dpTL[i-1][j].val == dpTL[i][j-1].val && dpTL[i-1][j].val == cur.val {
                   cur.cnt -= dpTL[i-1][j-1].cnt
               }
           }
           dpTL[i][j] = cur
       }
   }
   for i2 := 0; i2 < H; i2++ {
       for j2 := 0; j2 < W; j2++ {
           i1lim := i2 - sep; j1lim := j2 - sep
           if i1lim >= 0 && j1lim >= 0 {
               best1 := dpTL[i1lim][j1lim]
               sum := wgrid[i2][j2] + best1.val
               if sum > diag1Max {
                   diag1Max = sum; diag1Cnt = best1.cnt
               } else if sum == diag1Max {
                   diag1Cnt += best1.cnt
               }
           }
       }
   }
   // case4: top-right to bottom-left
   diag2Max := 0; var diag2Cnt int64 = 0
   dpTR := make([][]MaxCnt, H)
   for i := 0; i < H; i++ {
       dpTR[i] = make([]MaxCnt, W)
       for j := W-1; j >= 0; j-- {
           cur := MaxCnt{wgrid[i][j], 1}
           if i > 0 {
               cur = mergeRect(cur, dpTR[i-1][j])
           }
           if j+1 < W {
               cur = mergeRect(cur, dpTR[i][j+1])
           }
           if i > 0 && j+1 < W {
               if dpTR[i-1][j].val == dpTR[i][j+1].val && dpTR[i-1][j].val == cur.val {
                   cur.cnt -= dpTR[i-1][j+1].cnt
               }
           }
           dpTR[i][j] = cur
       }
   }
   for i2 := 0; i2 < H; i2++ {
       for j2 := 0; j2 < W; j2++ {
           i1lim := i2 - sep; j1lim := j2 + sep
           if i1lim >= 0 && j1lim < W {
               best1 := dpTR[i1lim][j1lim]
               sum := wgrid[i2][j2] + best1.val
               if sum > diag2Max {
                   diag2Max = sum; diag2Cnt = best1.cnt
               } else if sum == diag2Max {
                   diag2Cnt += best1.cnt
               }
           }
       }
   }
   // overall
   S := rowMax; cnt := rowCnt
   if colMax > S { S = colMax; cnt = colCnt } else if colMax == S { cnt += colCnt }
   if diag1Max > S { S = diag1Max; cnt = diag1Cnt } else if diag1Max == S { cnt += diag1Cnt }
   if diag2Max > S { S = diag2Max; cnt = diag2Cnt } else if diag2Max == S { cnt += diag2Cnt }
   if S <= 0 {
       fmt.Println("0 0")
   } else {
       fmt.Printf("%d %d\n", S, cnt)
   }
}

// mergeRect merges candidate cur with other dp value, simple max+count
func mergeRect(cur, other MaxCnt) MaxCnt {
   if other.val > cur.val {
       return other
   } else if other.val == cur.val {
       cur.cnt += other.cnt
   }
   return cur
}
