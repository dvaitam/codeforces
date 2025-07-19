package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
)

const INF = 2000000000

// Treap node for ordered map
type Node struct {
   key, value int
   prio       int
   left, right *Node
}

func rotateRight(t *Node) *Node {
   l := t.left
   t.left = l.right
   l.right = t
   return l
}

func rotateLeft(t *Node) *Node {
   r := t.right
   t.right = r.left
   r.left = t
   return r
}

func insertNode(t *Node, key, value int) *Node {
   if t == nil {
       return &Node{key: key, value: value, prio: rand.Int()}
   }
   if key == t.key {
       t.value = value
       return t
   }
   if key < t.key {
       t.left = insertNode(t.left, key, value)
       if t.left.prio > t.prio {
           t = rotateRight(t)
       }
   } else {
       t.right = insertNode(t.right, key, value)
       if t.right.prio > t.prio {
           t = rotateLeft(t)
       }
   }
   return t
}

func merge(a, b *Node) *Node {
   if a == nil {
       return b
   }
   if b == nil {
       return a
   }
   if a.prio > b.prio {
       a.right = merge(a.right, b)
       return a
   }
   b.left = merge(a, b.left)
   return b
}

func eraseNode(t *Node, key int) *Node {
   if t == nil {
       return nil
   }
   if key < t.key {
       t.left = eraseNode(t.left, key)
   } else if key > t.key {
       t.right = eraseNode(t.right, key)
   } else {
       t = merge(t.left, t.right)
   }
   return t
}

// upperBound returns pointer to node with smallest key > k, or nil.
func upperBound(t *Node, k int) *Node {
   var res *Node
   cur := t
   for cur != nil {
       if cur.key > k {
           res = cur
           cur = cur.left
       } else {
           cur = cur.right
       }
   }
   return res
}

// predecessor returns node with greatest key <= k, or nil.
func predecessor(t *Node, k int) *Node {
   var res *Node
   cur := t
   for cur != nil {
       if cur.key <= k {
           res = cur
           cur = cur.right
       } else {
           cur = cur.left
       }
   }
   return res
}

type Event struct{ color int; diff int64 }

func main() {
   rand.Seed(1)
   in := bufio.NewReader(os.Stdin)
   var N int
   var K int64
   fmt.Fscan(in, &N, &K)
   lo := make([]int, N)
   hi := make([]int, N)
   for i := 0; i < N; i++ {
       fmt.Fscan(in, &lo[i], &hi[i])
   }
   events := make([][]Event, N)
   var root *Node
   // initial map entries
   root = insertNode(root, 0, N)
   root = insertNode(root, INF, N)
   // build events
   for i := N - 1; i >= 0; i-- {
       l, h := lo[i], hi[i]
       it := upperBound(root, l)
       nit := predecessor(root, l)
       nv := nit.value
       // event for left segment
       minh := it.key
       if minh > h {
           minh = h
       }
       events[i] = append(events[i], Event{nit.value, int64(l - minh)})
       nl := nv
       // remove covered segments
       for it.key < h {
           nl = it.value
           // next node
           next := upperBound(root, it.key)
           // diff for this segment
           minh2 := next.key
           if minh2 > h {
               minh2 = h
           }
           events[i] = append(events[i], Event{it.value, int64(it.key - minh2)})
           // erase it
           root = eraseNode(root, it.key)
           it = next
       }
       // insert potential split at h
       if it.key != h {
           root = insertNode(root, h, nl)
       }
       // insert new segment
       root = insertNode(root, l, i)
       // event for new segment
       events[i] = append(events[i], Event{i, int64(h - l)})
   }
   // prepare cover array
   ccover := make([]int64, N+1)
   var works = func(M int64) bool {
       for i := range ccover {
           ccover[i] = 0
       }
       ccover[N] = INF
       ctot := int64(INF)
       cstop := N
       var ncnt int64
       for i := N - 1; i >= 0; i-- {
           for _, ev := range events[i] {
               ccover[ev.color] += ev.diff
               if ev.color <= cstop {
                   ctot += ev.diff
               }
           }
           for ctot-int64(ccover[cstop]) >= M {
               ctot -= int64(ccover[cstop])
               cstop--
           }
           ncnt += int64(N - cstop)
       }
       return ncnt >= K
   }
   var vrun = func(M int64) int64 {
       for i := range ccover {
           ccover[i] = 0
       }
       ccover[N] = INF
       ctot := int64(INF)
       cstop := N
       var cc int64
       var ncnt int64
       var ans int64
       for i := N - 1; i >= 0; i-- {
           for _, ev := range events[i] {
               ccover[ev.color] += ev.diff
               if ev.color <= cstop {
                   ctot += ev.diff
               } else {
                   cc += ev.diff * int64(N-ev.color)
               }
           }
           for ctot-int64(ccover[cstop]) >= M {
               ctot -= int64(ccover[cstop])
               cc += int64(N-cstop) * int64(ccover[cstop])
               cstop--
           }
           nv := N - cstop
           ans += cc + int64(nv)*ctot
           ncnt += int64(nv)
       }
       return ans - (ncnt-K)*M
   }
   // binary search best M
   var hlo, hhi int64 = 0, 1000000000
   for hlo < hhi {
       mid := (hlo + hhi + 1) >> 1
       if works(mid) {
           hlo = mid
       } else {
           hhi = mid - 1
       }
   }
   res := vrun(hlo)
   fmt.Println(res)
}
