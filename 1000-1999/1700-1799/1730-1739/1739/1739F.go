package main

import (
   "bufio"
   "fmt"
   "os"
)

const ALPH = 12
const S = (1 << ALPH) - 1

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   weights := make([]int32, n)
   ss := make([]string, n)
   var sumLen int
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &weights[i], &ss[i])
       sumLen += len(ss[i])
   }
   // estimate trie size
   maxStates := sumLen*2 + 5
   // trie structures
   nxt := make([][]int, maxStates)
   for i := range nxt {
       nxt[i] = make([]int, ALPH+1)
   }
   ed := make([]int32, maxStates)
   fail := make([]int, maxStates)
   tot := 1
   // process each string
   for i := 0; i < n; i++ {
       s := ss[i]
       w := weights[i]
       // build graph of letters
       var seen [ALPH + 1]bool
       var mp [ALPH + 1][ALPH + 1]bool
       deg := make([]int, ALPH+1)
       v := make([][]int, ALPH+1)
       var vs, es int
       // scan string
       for j := 0; j < len(s); j++ {
           c := int(s[j]-'a') + 1
           if !seen[c] {
               seen[c] = true
               vs++
           }
           if j > 0 {
               p := int(s[j-1]-'a') + 1
               if !mp[p][c] {
                   mp[p][c] = true
                   mp[c][p] = true
                   deg[p]++
                   deg[c]++
                   v[p] = append(v[p], c)
                   v[c] = append(v[c], p)
                   es++
               }
           }
       }
       // check if path
       flg := (vs == es+1)
       rt := -1
       for j := 1; j <= ALPH; j++ {
           if deg[j] == 1 {
               rt = j
           }
           if deg[j] > 2 {
               flg = false
           }
       }
       // reset seen for dfs
       for j := 1; j <= ALPH; j++ {
           seen[j] = false
       }
       if flg && rt != -1 {
           // extract path via dfs
           var now []int
           var dfs func(int)
           dfs = func(x int) {
               seen[x] = true
               now = append(now, x)
               for _, y := range v[x] {
                   if !seen[y] {
                       dfs(y)
                   }
               }
           }
           dfs(rt)
           // insert forward and reverse
           insertPath := func(path []int) {
               p := 1
               for _, x := range path {
                   if nxt[p][x] == 0 {
                       tot++
                       nxt[p][x] = tot
                   }
                   p = nxt[p][x]
               }
               ed[p] += w
           }
           insertPath(now)
           // reverse
           for l, r := 0, len(now)-1; l < r; l, r = l+1, r-1 {
               now[l], now[r] = now[r], now[l]
           }
           insertPath(now)
       }
   }
   // build Aho-Corasick
   queue := make([]int, 0, tot)
   // initialize
   for i := 1; i <= ALPH; i++ {
       if nxt[1][i] != 0 {
           fail[nxt[1][i]] = 1
           queue = append(queue, nxt[1][i])
       } else {
           nxt[1][i] = 1
       }
   }
   // bfs
   for qi := 0; qi < len(queue); qi++ {
       x := queue[qi]
       ed[x] += ed[fail[x]]
       for i := 1; i <= ALPH; i++ {
           if nxt[x][i] != 0 {
               fail[nxt[x][i]] = nxt[fail[x]][i]
               queue = append(queue, nxt[x][i])
           } else {
               nxt[x][i] = nxt[fail[x]][i]
           }
       }
   }
   // DP
   finf := int32(-1e9)
   // f[mask][state]
   f := make([][]int32, S+1)
   recMask := make([][]int16, S+1)
   recState := make([][]int16, S+1)
   for mask := 0; mask <= S; mask++ {
       f[mask] = make([]int32, tot+1)
       recMask[mask] = make([]int16, tot+1)
       recState[mask] = make([]int16, tot+1)
       for j := 1; j <= tot; j++ {
           f[mask][j] = finf
       }
   }
   f[0][1] = 0
   // transitions
   for mask := 0; mask <= S; mask++ {
       for state := 1; state <= tot; state++ {
           if f[mask][state] == finf {
               continue
           }
           base := f[mask][state]
           for c := 1; c <= ALPH; c++ {
               bit := 1 << (c - 1)
               if mask&bit != 0 {
                   continue
               }
               nm := mask | bit
               ns := nxt[state][c]
               val := base + ed[ns]
               if val > f[nm][ns] {
                   f[nm][ns] = val
                   recMask[nm][ns] = int16(mask)
                   recState[nm][ns] = int16(state)
               }
           }
       }
   }
   // find answer
   full := S
   ans := finf
   bestState := 1
   for st := 1; st <= tot; st++ {
       if f[full][st] > ans {
           ans = f[full][st]
           bestState = st
       }
   }
   // reconstruct
   res := make([]int, ALPH)
   mask := full
   state := bestState
   pos := ALPH - 1
   for mask != 0 {
       pm := int(recMask[mask][state])
       ps := int(recState[mask][state])
       diff := mask ^ pm
       var letter int
       for i := 0; i < ALPH; i++ {
           if diff&(1<<i) != 0 {
               letter = i + 1
               break
           }
       }
       res[pos] = letter
       pos--
       mask = pm
       state = ps
   }
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < ALPH; i++ {
       writer.WriteByte(byte('a' + res[i] - 1))
   }
}
