package main

import (
   "bufio"
   "os"
   "runtime"
)

const base = 91138233

var (
   reader *bufio.Reader
   writer *bufio.Writer
   n      int
   adj    [][]int
   parent []int
   depth  []int
   size   []int
   heavy  []int
   head   []int
   pos    []int
   curPos int
   ba     []uint64
   powB   []uint64
   hf     []uint64
   hr     []uint64
)

// fast IO
func initIO() {
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
}
// fast integer reading without allocations
func readInt() int {
   var x int
   var c byte
   var err error
   // skip non-digits
   for {
       c, err = reader.ReadByte()
       if err != nil {
           return 0
       }
       if c >= '0' && c <= '9' {
           break
       }
   }
   x = int(c - '0')
   for {
       c, err = reader.ReadByte()
       if err != nil || c < '0' || c > '9' {
           break
       }
       x = x*10 + int(c-'0')
   }
   return x
}
func writeInt(x int) {
   if x == 0 {
       writer.WriteByte('0')
       writer.WriteByte('\n')
       return
   }
   var buf [20]byte
   i := len(buf)
   for x > 0 {
       i--
       buf[i] = byte('0' + x%10)
       x /= 10
   }
   writer.Write(buf[i:])
   writer.WriteByte('\n')
}

func main() {
   runtime.GOMAXPROCS(runtime.NumCPU())
   initIO()
   n = readInt()
   s := make([]byte, n+1)
   for i := 1; i <= n; i++ {
       ch, _ := reader.ReadByte()
       for ch < 'a' || ch > 'z' {
           ch, _ = reader.ReadByte()
       }
       s[i] = ch
   }
   adj = make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       u := readInt()
       v := readInt()
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   parent = make([]int, n+1)
   depth = make([]int, n+1)
   size = make([]int, n+1)
   heavy = make([]int, n+1)
   // dfs1: parent, depth, size, heavy
   type stackEntry struct{ u, idx int }
   stack := make([]stackEntry, 0, n*2)
   stack = append(stack, stackEntry{1, 0})
   parent[1] = 0
   depth[1] = 0
   order := make([]int, 0, n)
   // pre and post
   for len(stack) > 0 {
       top := &stack[len(stack)-1]
       u := top.u
       if top.idx == 0 {
           order = append(order, u)
       }
       if top.idx < len(adj[u]) {
           v := adj[u][top.idx]
           top.idx++
           if v != parent[u] {
               parent[v] = u
               depth[v] = depth[u] + 1
               stack = append(stack, stackEntry{v, 0})
           }
       } else {
           // post
           size[u] = 1
           maxSz := 0
           for _, v := range adj[u] {
               if v != parent[u] {
                   size[u] += size[v]
                   if size[v] > maxSz {
                       maxSz = size[v]
                       heavy[u] = v
                   }
               }
           }
           stack = stack[:len(stack)-1]
       }
   }
   // hld
   head = make([]int, n+1)
   pos = make([]int, n+1)
   ba = make([]uint64, n+1)
   curPos = 1
   // iterative hld
   type hldEntry struct{ u, h int }
   hstack := []hldEntry{{1, 1}}
   for len(hstack) > 0 {
       ent := hstack[len(hstack)-1]
       hstack = hstack[:len(hstack)-1]
       u, h := ent.u, ent.h
       head[u] = h
       pos[u] = curPos
       ba[curPos] = uint64(s[u] - 'a' + 1)
       curPos++
       // push light children first
       for _, v := range adj[u] {
           if v != parent[u] && v != heavy[u] {
               hstack = append(hstack, hldEntry{v, v})
           }
       }
       if heavy[u] != 0 {
           hstack = append(hstack, hldEntry{heavy[u], h})
       }
   }
   // hashes
   powB = make([]uint64, n+1)
   hf = make([]uint64, n+1)
   hr = make([]uint64, n+1)
   powB[0] = 1
   for i := 1; i <= n; i++ {
       powB[i] = powB[i-1] * base
       hf[i] = hf[i-1]*base + ba[i]
       hr[i] = hr[i-1]*base + ba[n-i+1]
   }
   // queries
   m := readInt()
   for qi := 0; qi < m; qi++ {
       a := readInt()
       b := readInt()
       c := readInt()
       d := readInt()
       l1 := lca(a, b)
       l2 := lca(c, d)
       len1 := depth[a] + depth[b] - 2*depth[l1] + 1
       len2 := depth[c] + depth[d] - 2*depth[l2] + 1
       maxL := len1
       if len2 < maxL {
           maxL = len2
       }
       seg1 := makeSegments(a, b, l1)
       seg2 := makeSegments(c, d, l2)
       // two pointers over segs
       i1, i2 := 0, 0
       o1, o2 := 0, 0
       matched := 0
       res := 0
       for i1 < len(seg1) && i2 < len(seg2) && matched < maxL {
           s1 := seg1[i1]
           s2 := seg2[i2]
           rem1 := s1.length() - o1
           rem2 := s2.length() - o2
           rem := rem1
           if rem2 < rem {
               rem = rem2
           }
           if maxL-matched < rem {
               rem = maxL - matched
           }
           // compare full rem
           if hashSegment(s1, o1, rem) == hashSegment(s2, o2, rem) {
               matched += rem
               o1 += rem
               if o1 == s1.length() {
                   i1++
                   o1 = 0
               }
               o2 += rem
               if o2 == s2.length() {
                   i2++
                   o2 = 0
               }
               if matched == maxL {
                   break
               }
               continue
           }
           // mismatch inside, binary search
           lo, hi := 0, rem
           for lo < hi {
               mid := (lo + hi + 1) >> 1
               if hashSegment(s1, o1, mid) == hashSegment(s2, o2, mid) {
                   lo = mid
               } else {
                   hi = mid - 1
               }
           }
           res = matched + lo
           matched = maxL // to break
           break
       }
       if matched < maxL {
           res = matched
       } else if res < matched {
           res = matched
       }
       writeInt(res)
   }
   writer.Flush()
}

// segment on base array
type segment struct{ l, r int; dir int }
func (s segment) length() int {
   if s.dir == 1 {
       return s.r - s.l + 1
   }
   return s.r - s.l + 1
}

// makeSegments for path a->b through lca
func makeSegments(a, b, l int) []segment {
   segsUp := make([]segment, 0)
   u := a
   for head[u] != head[l] {
       segsUp = append(segsUp, segment{pos[head[u]], pos[u], -1})
       u = parent[head[u]]
   }
   if u != l {
       segsUp = append(segsUp, segment{pos[l] + 1, pos[u], -1})
   }
   // down segments
   tmp := make([]segment, 0)
   v := b
   for head[v] != head[l] {
       tmp = append(tmp, segment{pos[head[v]], pos[v], 1})
       v = parent[head[v]]
   }
   // final
   tmp = append(tmp, segment{pos[l], pos[v], 1})
   // reverse tmp to get from l to b
   segsDown := make([]segment, len(tmp))
   for i := range tmp {
       segsDown[i] = tmp[len(tmp)-1-i]
   }
   // combine
   segs := append(segsUp, segsDown...)
   return segs
}

// hash functions
func hashF(l, r int) uint64 {
   return hf[r] - hf[l-1]*powB[r-l+1]
}
func hashR(l, r int) uint64 {
   // reverse of ba[l..r]
   rl := n - r + 1
   rr := n - l + 1
   return hr[rr] - hr[rl-1]*powB[rr-rl+1]
}
func hashSegment(s segment, offset, length int) uint64 {
   if length == 0 {
       return 0
   }
   if s.dir == 1 {
       lpos := s.l + offset
       return hashF(lpos, lpos+length-1)
   }
   // dir -1
   rpos := s.r - offset
   return hashR(rpos-length+1, rpos)
}

// LCA via HLD
func lca(u, v int) int {
   for head[u] != head[v] {
       if depth[head[u]] > depth[head[v]] {
           u = parent[head[u]]
       } else {
           v = parent[head[v]]
       }
   }
   if depth[u] < depth[v] {
       return u
   }
   return v
}
