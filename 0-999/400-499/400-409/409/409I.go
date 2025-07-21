package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   // directions: 0=right,1=down,2=left,3=up
   const INF = 1<<30
   dist := make([][4]int, n)
   for i := 0; i < n; i++ {
       for d := 0; d < 4; d++ {
           dist[i][d] = INF
       }
   }
   // deque for 0-1 BFS: entries of (pos, dir)
   type state struct{ p, d int }
   dq := make([]state, 0, n*4)
   pushFront := func(st state) { dq = append([]state{st}, dq...) }
   pushBack := func(st state) { dq = append(dq, st) }
   popFront := func() state {
       st := dq[0]
       dq = dq[1:]
       return st
   }
   // start at pos 0, dir right
   dist[0][0] = 0
   pushBack(state{0, 0})
   best := INF
   for len(dq) > 0 {
       u := popFront()
       p, d := u.p, u.d
       cdist := dist[p][d]
       // prune if already worse than best
       if cdist >= best {
           continue
       }
       c := s[p]
       // if termination
       if c == '@' {
           best = cdist
           continue
       }
       // compute next directions
       var ndirs []int
       switch c {
       case '>':
           ndirs = []int{0}
       case 'v':
           ndirs = []int{1}
       case '<':
           ndirs = []int{2}
       case '^':
           ndirs = []int{3}
       case '_':
           // pop empty treated as zero => go right
           ndirs = []int{0}
       case '|':
           // pop empty => zero => go down
           ndirs = []int{1}
       case '?':
           ndirs = []int{0, 1, 2, 3}
       default:
           ndirs = []int{d}
       }
       for _, nd := range ndirs {
           // cost for reading input at '&'
           add := 0
           if c == '&' {
               add = 1
           }
           np := p
           // move from p according to nd
           switch nd {
           case 0:
               np = (p + 1) % n
           case 2:
               np = (p - 1 + n) % n
           case 1, 3:
               // up/down stays on same row (only one row)
               np = p
           }
           ndist := cdist + add
           if ndist < dist[np][nd] {
               dist[np][nd] = ndist
               st := state{np, nd}
               if add == 1 {
                   pushBack(st)
               } else {
                   pushFront(st)
               }
           }
       }
   }
   if best == INF {
       fmt.Println("false")
   } else {
       // lex smallest sequence of length best: all zeros
       if best > 0 {
           // build string of zeros
           out := make([]byte, best)
           for i := range out {
               out[i] = '0'
           }
           fmt.Println(string(out))
       } else {
           // no inputs needed: empty sequence => print empty line
           fmt.Println()
       }
   }
}
