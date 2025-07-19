package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD uint64 = 92233720368547753
const E uint64 = 29

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read integers and tokens
   N := readInt(reader)
   K := readInt(reader)
   Sbig := readToken(reader)
   M := readInt(reader)
   // build rolling hashes of all length-K substrings of circular Sbig
   NK := N * K
   EK := uint64(1)
   for i := 1; i < K; i++ {
       EK = EK * E % MOD
   }
   H := make([]uint64, NK)
   // initial hash
   var h0 uint64
   for i := 0; i < K; i++ {
       h0 = (h0*E + uint64(Sbig[i]-'a'+1)) % MOD
   }
   if NK > 0 {
       H[0] = h0
   }
   for i := 1; i < NK; i++ {
       // remove Sbig[i-1], add Sbig[(i+K-1)%NK]
       rem := uint64(Sbig[i-1]-'a'+1) * EK % MOD
       h := H[i-1]
       if h < rem {
           h += MOD
       }
       h = (h - rem) % MOD
       h = (h*E + uint64(Sbig[(i+K-1)%NK]-'a'+1)) % MOD
       H[i] = h
   }
   // read pattern hashes
   patMap := make(map[uint64]int, M)
   for i := 1; i <= M; i++ {
       pat := readToken(reader)
       var ph uint64
       for j := 0; j < K; j++ {
           ph = (ph*E + uint64(pat[j]-'a'+1)) % MOD
       }
       // store 1-based index
       patMap[ph] = i
   }
   // visited marks
   Z := make([]int, M)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // try each offset
   for o := 0; o < K; o++ {
       A := make([]int, 0, N)
       s := o
       for i := 0; i < N; i++ {
           ph := H[s]
           j, ok := patMap[ph]
           if !ok {
               break
           }
           // j in [1..M], index j-1
           if Z[j-1] == o+1 {
               break
           }
           Z[j-1] = o + 1
           A = append(A, j)
           s += K
           if s >= NK {
               s -= NK
           }
       }
       if len(A) == N {
           fmt.Fprint(writer, "YES\n")
           for i, v := range A {
               if i == 0 {
                   fmt.Fprint(writer, v)
               } else {
                   fmt.Fprint(writer, " ", v)
               }
           }
           fmt.Fprint(writer, "\n")
           return
       }
   }
   fmt.Fprint(writer, "NO\n")
}

// readInt reads next integer from reader
func readInt(r *bufio.Reader) int {
   var x int
   var neg bool
   // skip non-digit
   for {
       b, err := r.ReadByte()
       if err != nil {
           return x
       }
       if b == '-' {
           neg = true
           break
       }
       if b >= '0' && b <= '9' {
           x = int(b - '0')
           break
       }
   }
   // read digits
   for {
       b, err := r.Peek(1)
       if err != nil {
           break
       }
       if b[0] < '0' || b[0] > '9' {
           break
       }
       r.ReadByte()
       x = x*10 + int(b[0]-'0')
   }
   if neg {
       x = -x
   }
   return x
}

// readToken reads next non-space token
func readToken(r *bufio.Reader) string {
   // skip spaces
   for {
       b, err := r.Peek(1)
       if err != nil {
           return ""
       }
       if b[0] > ' ' {
           break
       }
       r.ReadByte()
   }
   var buf []byte
   for {
       b, err := r.Peek(1)
       if err != nil || b[0] <= ' ' {
           break
       }
       r.ReadByte()
       buf = append(buf, b[0])
   }
   return string(buf)
}
