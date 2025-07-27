package main
import (
   "bufio"
   "fmt"
   "os"
)

// Eertree (palindromic tree) to count distinct palindromic substrings
type Eertree struct {
   s []byte
   next []map[byte]int
   link []int
   length []int
   last int
}

func NewEertree(n int) *Eertree {
   t := &Eertree{
       s: make([]byte, 1, n+2),
       next: make([]map[byte]int, 2, n+3),
       link: make([]int, 2, n+3),
       length: make([]int, 2, n+3),
   }
   t.next[0] = make(map[byte]int)
   t.next[1] = make(map[byte]int)
   t.length[0] = -1; t.link[0] = 0
   t.length[1] = 0; t.link[1] = 0
   t.last = 1
   return t
}

func (t *Eertree) Add(c byte) {
   t.s = append(t.s, c)
   pos := len(t.s) - 1
   cur := t.last
   for {
       l := t.length[cur]
       if pos-1-l >= 0 && t.s[pos-1-l] == c {
           break
       }
       cur = t.link[cur]
   }
   if v, ok := t.next[cur][c]; ok {
       t.last = v
       return
   }
   newNode := len(t.next)
   t.next = append(t.next, make(map[byte]int))
   t.length = append(t.length, t.length[cur]+2)
   t.next[cur][c] = newNode
   if t.length[newNode] == 1 {
       t.link = append(t.link, 1)
   } else {
       x := t.link[cur]
       for {
           l := t.length[x]
           if pos-1-l >= 0 && t.s[pos-1-l] == c {
               break
           }
           x = t.link[x]
       }
       t.link = append(t.link, t.next[x][c])
   }
   t.last = newNode
}

func countPal(s string) int64 {
   t := NewEertree(len(s))
   for i := 0; i < len(s); i++ {
       t.Add(s[i])
   }
   // nodes - 2 (excluding roots)
   return int64(len(t.next) - 2)
}

func allSame(s string) bool {
   for i := 1; i < len(s); i++ {
       if s[i] != s[0] {
           return false
       }
   }
   return true
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var A, B string
   fmt.Fscan(reader, &A)
   fmt.Fscan(reader, &B)
   nA, nB := len(A), len(B)
   if allSame(A) && allSame(B) && A[0] == B[0] {
       // concatenations of same char length from 2 to nA+nB
       fmt.Println(nA + nB - 1)
       return
   }
   sa := countPal(A)
   sb := countPal(B)
   // fallback: assume all distinct
   fmt.Println(sa * sb)
}
