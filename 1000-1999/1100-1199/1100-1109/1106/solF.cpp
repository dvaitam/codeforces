#include <bits/stdc++.h>
using namespace std;

typedef long long LL;
const int N = 110;
const int mod = 998244353;
const LL sqr = 1ll * (mod - 1) * (mod - 1) * 6;

struct Mat {
  int a[N][N];
  Mat() { memset(a, 0, sizeof a); }
};
int s, a[N], n, m;

inline Mat Mul(const Mat &a, const Mat &b) {
  static Mat tmp;
  for (int i = 1; i <= s; i ++)
    for (int j = 1; j <= s; j ++) {
      LL u = 0;
      for (int k = 1; k <= s; k ++) {
        u += 1ll * a.a[i][k] * b.a[k][j];
        if (u >= sqr) u -= sqr;
      }
      tmp.a[i][j] = u % (mod - 1);
    }
  return tmp;
}

inline Mat Pow(Mat x, int exp) {
  Mat I;
  for (int i = 1; i <= s; i ++) I.a[i][i] = 1;
  for (; exp; exp >>= 1, x = Mul(x, x))
    if (exp & 1) I = Mul(I, x);
  return I;
}

inline LL Pow(LL x, LL exp) {
  LL res = 1;
  for (; exp; exp >>= 1, x = x * x % mod)
    if (exp & 1) res = res * x % mod;
  return res;
}

inline int BSGS(int x) {
  unordered_map <int, int> rec;
  rec.rehash(1000007);
  int s = sqrt(mod) + 1;
  LL u = 1;
  for (int i = 0; i < s; i ++) {
    rec[u] = i; u = u * 3 % mod;
  }
  u = Pow(u, mod - 2);
  for (int i = 0; i <= mod; i += s) {
    if (rec.count(x)) return i + rec[x];
    x = x * u % mod;
  }
  throw;
}

inline void exgcd(LL &x, LL &y, LL a, LL b) {
  if (a == 1 && b == 0) {
    x = 1; y = 0; return;
  }
  exgcd(x, y, b, a % b);
  int t = x;
  x = y;
  y = t - (a / b) * y;
}

inline void solve(LL k, LL m) {
  LL p = mod - 1;
  if (m % __gcd(p, k)) {
    puts("-1"); return;
  }
  LL x, y;
  LL d = __gcd(p, k);
  m /= d;
  p /= d;
  k /= d;
  exgcd(x, y, k, p);
  x = x * m % (mod - 1);
  if (x < mod - 1) x += mod - 1;
  printf("%lld\n", Pow(3, x));
}

int main() {
  scanf("%d", &s);
  for (int i = 1; i <= s; i ++) scanf("%d", &a[i]);
  scanf("%d%d", &n, &m);
  Mat base;
  for (int i = 1; i <= s; i ++) base.a[1][i] = a[i];
  for (int i = 2; i <= s; i ++)
    base.a[i][i - 1] = 1;
  LL pw = Pow(base, n - s).a[1][1];
  m = BSGS(m);
  solve(pw, m);
  return 0;
}