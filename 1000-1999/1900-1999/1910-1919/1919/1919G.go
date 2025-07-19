package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "sort"
   "strconv"
)

// Info holds auxiliary flags for nodes
type Info struct {
   deci int
   val  int
}

// Edge represents a directed edge
type Edge struct{ u, v int }

var (
   n      int
   s      [][]byte
   id     []int
   has    []uint64
   hashed []uint64
   tothas []uint64
   p      []Info
   tem    []Info
   ans    []Edge
   sta    []int
   flag   bool
   all    uint64
   rnd    *rand.Rand
)

// solve recursively processes range [l,r]
func solve(l, r int) {
   if l == r {
       return
   }
   // detect a node with mismatched hash
   detectionPos := 0
   var tot uint64
   for i := l; i <= r; i++ {
       tot ^= has[id[i]]
   }
   for i := l; i <= r; i++ {
       u := id[i]
       if hashed[u] != tot && hashed[u] != 0 {
           detectionPos = i
           break
       }
   }
   if detectionPos == 0 {
       // all hashed[u] are either 0 or tot
       cnt := [3]int{}
       for i := l; i <= r; i++ {
           u := id[i]
           if p[u].deci != 0 {
               cnt[p[u].val]++
           } else if hashed[u] != 0 {
               cnt[1]++
               p[u].val = 1
           } else {
               cnt[0]++
               p[u].val = 0
           }
       }
       if cnt[1] == 0 && cnt[0] > 0 {
           flag = true
           return
       }
       if cnt[0] < cnt[1]+1 && cnt[1] != 0 {
           flag = true
           return
       }
       if cnt[1] == 0 {
           for i := l + 1; i <= r; i++ {
               ans = append(ans, Edge{id[i-1], id[i]})
           }
           return
       }
       // prepare temporary slice and sort by val
       for i := l; i <= r; i++ {
           tem[i] = p[id[i]]
       }
       subTem := tem[l : r+1]
       sort.Slice(subTem, func(i, j int) bool { return subTem[i].val < subTem[j].val })
       innerPos := l
       for innerPos <= r && tem[innerPos].val == 0 {
           innerPos++
       }
       tmp := innerPos
       // sort id by p[id].val
       subId := id[l : r+1]
       sort.Slice(subId, func(i, j int) bool { return p[subId[i]].val < p[subId[j]].val })
       pos := tmp
       pos1 := l + 1
       cnt2 := cnt[2]
       for pos <= r-cnt2 {
           ans = append(ans, Edge{subId[pos-l], subId[pos1-l]})
           ans = append(ans, Edge{subId[pos-l], subId[pos1-1-l]})
           pos++
           pos1++
       }
       for pos1 < tmp {
           ans = append(ans, Edge{subId[pos1-l], subId[tmp-l]})
           pos1++
       }
       for i := r - cnt2 + 1; i <= r; i++ {
           ans = append(ans, Edge{subId[tmp-l], subId[i-l]})
       }
   } else {
       // partition nodes with edge to detection node
       tp := 0
       U := id[detectionPos]
       for i := l; i <= r; i++ {
           u := id[i]
           if s[u][U] == '1' {
               tp++
               sta[tp] = i
           }
       }
       // move these nodes to front of [l,r]
       for i := 1; i <= tp; i++ {
           if sta[i] != l+i-1 {
               id[sta[i]], id[l+i-1] = id[l+i-1], id[sta[i]]
           }
       }
       // update hashed values across partition
       for i := l; i <= l+tp-1; i++ {
           for j := l + tp; j <= r; j++ {
               u := id[i]
               v := id[j]
               if s[u][v] == '1' {
                   hashed[v] ^= has[u]
               }
           }
       }
       for i := l + tp; i <= r; i++ {
           for j := l; j <= l+tp-1; j++ {
               u := id[i]
               v := id[j]
               if s[u][v] == '1' {
                   hashed[v] ^= has[u]
               }
           }
       }
       // find candidate V
       pos2 := 0
       for i := l + tp; i <= r; i++ {
           u := id[i]
           if s[u][u] == '0' || p[u].deci != 0 || (tothas[u]^tothas[U]) != all {
               continue
           }
           ok := true
           for j := l; j <= l+tp-1; j++ {
               if s[id[j]][u] != '0' {
                   ok = false
                   break
               }
           }
           if ok {
               pos2 = i
           }
       }
       if pos2 == 0 {
           flag = true
           return
       }
       V := id[pos2]
       ans = append(ans, Edge{U, V})
       p[V] = Info{1, 2}
       p[U] = Info{1, 2}
       // consistency checks
       for i := l; i <= l+tp-1; i++ {
           for j := l + tp; j <= r; j++ {
               u := id[i]
               v := id[j]
               if v == V {
                   continue
               }
               if s[u][v] != s[V][v] {
                   flag = true
                   return
               }
           }
       }
       for i := l + tp; i <= r; i++ {
           for j := l; j <= l+tp-1; j++ {
               u := id[i]
               v := id[j]
               if v == U {
                   continue
               }
               if s[u][v] != s[U][v] {
                   flag = true
                   return
               }
           }
       }
       solve(l, l+tp-1)
       if flag {
           return
       }
       solve(l+tp, r)
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // read n
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // read adjacency matrix
   s = make([][]byte, n+1)
   var line string
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &line)
       s[i] = make([]byte, n+1)
       for j := 1; j <= n; j++ {
           s[i][j] = line[j-1]
       }
   }
   // special case n==1
   if n == 1 {
       if s[1][1] == '1' {
           writer.WriteString("NO\n")
       } else {
           writer.WriteString("YES\n")
       }
       return
   }
   // init
   id = make([]int, n+1)
   has = make([]uint64, n+1)
   hashed = make([]uint64, n+1)
   tothas = make([]uint64, n+1)
   p = make([]Info, n+1)
   tem = make([]Info, n+1)
   sta = make([]int, n+1)
   ans = make([]Edge, 0, n*2)
   rnd = rand.New(rand.NewSource(20080623))
   // random hashes
   for i := 1; i <= n; i++ {
       id[i] = i
       // generate uint64 from two int63
       has[i] = (uint64(rnd.Int63()) << 32) | uint64(rnd.Int63()&0xffffffff)
       all ^= has[i]
   }
   // compute initial hashed values
   for i := 1; i <= n; i++ {
       for j := 1; j <= n; j++ {
           if s[i][j] == '1' {
               hashed[j] ^= has[i]
           }
       }
   }
   // copy for checks
   copy(tothas, hashed)
   // solve
   solve(1, n)
   if flag {
       writer.WriteString("NO")
   } else {
       writer.WriteString("YES\n")
       for _, e := range ans {
           writer.WriteString(strconv.Itoa(e.u) + " " + strconv.Itoa(e.v) + "\n")
       }
   }
}
