#include <bits/stdc++.h>

using namespace std;
typedef unsigned long long ull;
const ull l = 0xFFFFFFFFFFFFFFFF;
const int p[] = {3, 5, 17, 257, 641, 65537, 6700417}, m[] = {1, 2, 6, 16, 25, 256, 2588}, I[] = {2, 2, 1, 32, 590, 16384, 3883315};
ull F[65][65], T[65][257][257], a, b, ans; int t, r[7];
ull mul(ull, ull) ;
ull mul2(ull a, ull b) {
if (a == 1 || b == 1) return a * b;
if (a < b) swap(a, b);
int x = __builtin_ctzll(a), y = __builtin_ctzll(b), z = 1 << 63 - __builtin_clz(x);
if (y & z) return mul((1ull << z >> 1) * 3, mul2(a >> z, b >> z));
return mul2(a >> z, b) << z;
}
ull mul(ull a, ull b) {
ull x, y, c = 0;
for (ull i=a; i; i^=x) {
x = 1ull << __builtin_ctzll(i);
for (ull j=b; j; j^=y)
y = 1ull << __builtin_ctzll(j), c ^= mul2(x, y);
} return c;
}
ull prod(ull a, ull b) {
ull c = 0;
for (int i=0; i<8; i++)
for (int j=0; j<8; j++)
c ^= T[i<<3|j][a>>(i<<3)&255][b>>(j<<3)&255];
return c;
}
ull qpow(ull x, ull y) {
ull t = 1;
for (; y; y>>=1) {
if (y & 1) t = prod(t, x);
x = prod(x, x);
} return t;
}
int main() {
for (int x=0; x<64; x++)
for (int y=0; y<64; y++)
F[x][y] = mul2(1ull << x, 1ull << y);
for (int i=0; i<8; i++)
for (int j=0; j<8; j++) {
for (int x=0; x<8; x++)
for (int y=0; y<8; y++)
T[i<<3|j][1<<x][1<<y] = F[i<<3|x][j<<3|y];
for (int a=1; a<256; a++)
for (int b=1; b<256; b++)
if (a^a&-a) T[i<<3|j][a][b] = T[i<<3|j][a&-a][b] ^ T[i<<3|j][a^a&-a][b];
else if (b^b&-b) T[i<<3|j][a][b] = T[i<<3|j][a][b&-b] ^ T[i<<3|j][a][b^b&-b];
}
cin >> t;
while (t --) {
scanf ("%llu%llu", &a, &b), ans = 0;
for (int i=0; i<7; i++) {
ull A = qpow(a, l / p[i]), B = qpow(b, l / p[i]), Am = qpow(A, m[i]), P = 1; int res = -1;
unordered_map <ull, int> id;
for (int j=0; j<m[i]; j++)
id[B] = j, B = prod(B, A);
for (int j=0, u=(p[i]-1)/m[i]; j<=u; j++) {
if (id. find(P) != id. end()) {
res = (j * m[i] - id[P] + p[i]) % p[i];
break;
} P = prod(P, Am);
}
if (~ res) r[i] = res;
else { puts ("-1"); goto E; }
}
for (int i=0; i<7; i++)
ans = (ans + (__int128) I[i] * (l / p[i]) * r[i]) % l;
printf ("%llu\n", ans); E : ;
}
}