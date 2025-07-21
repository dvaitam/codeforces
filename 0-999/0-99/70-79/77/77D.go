package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

const MOD = 1000000007

// maskKey represents a bitmask of up to 256 bits
type maskKey [4]uint64

func (k maskKey) get(pos int) bool {
   w := pos >> 6
   b := pos & 63
   return (k[w]>>b)&1 != 0
}

func (k *maskKey) set(pos int) {
   w := pos >> 6
   b := pos & 63
   k[w] |= 1 << b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   // read grid of patterns
   totalR := 4*n + 1
   grid := make([]string, totalR)
   for i := 0; i < totalR; i++ {
       line, _ := reader.ReadString('\n')
       grid[i] = strings.TrimRight(line, "\r\n")
   }
   // determine allowed orientations
   allowH := make([][]bool, n)
   allowV := make([][]bool, n)
   for i := 0; i < n; i++ {
       allowH[i] = make([]bool, m)
       allowV[i] = make([]bool, m)
       for j := 0; j < m; j++ {
           // extract 3x3 block
           P := make([]byte, 9)
           for di := 0; di < 3; di++ {
               for dj := 0; dj < 3; dj++ {
                   P[di*3+dj] = grid[1+4*i+di][1+4*j+dj]
               }
           }
           // rotate P by 90deg: R[ri][ci] = P[2-ci][ri]
           R := make([]byte, 9)
           for ri := 0; ri < 3; ri++ {
               for ci := 0; ci < 3; ci++ {
                   R[ri*3+ci] = P[(2-ci)*3+ri]
               }
           }
           sP := string(P)
           sR := string(R)
           if sP == sR {
               allowH[i][j] = true
               allowV[i][j] = true
           } else if sP < sR {
               // P is base -> horizontal
               allowH[i][j] = true
           } else {
               // P is rotated -> vertical
               allowV[i][j] = true
           }
       }
   }
   // DP over rows
   dpCur := make(map[maskKey]int)
   var zeroKey maskKey
   dpCur[zeroKey] = 1
   for i := 0; i < n; i++ {
       dpNext := make(map[maskKey]int)
       // for each incoming mask
       for inMask, cnt := range dpCur {
           // dfs to fill row i
           var dfs func(pos int, nextMask maskKey)
           dfs = func(pos int, nextMask maskKey) {
               if pos >= m {
                   dpNext[nextMask] = (dpNext[nextMask] + cnt) % MOD
                   return
               }
               if inMask.get(pos) {
                   dfs(pos+1, nextMask)
                   return
               }
               // try vertical
               if allowV[i][pos] && i+1 < n && allowV[i+1][pos] {
                   var nm2 = nextMask
                   nm2.set(pos)
                   dfs(pos+1, nm2)
               }
               // try horizontal
               if pos+1 < m && !inMask.get(pos+1) && allowH[i][pos] && allowH[i][pos+1] {
                   dfs(pos+2, nextMask)
               }
           }
           dfs(0, zeroKey)
       }
       dpCur = dpNext
   }
   // answer is dpCur[0]
   fmt.Println(dpCur[zeroKey])
}
