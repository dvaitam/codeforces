#include <algorithm>
#include <cctype>
#include <cstdio>
#include <cstring>

typedef long long LL;

inline int readInt() {
  int ans = 0, c, f = 1;
  while (!isdigit(c = getchar()))
    if (c == '-') f *= -1;
  do ans = ans * 10 + c - '0';
  while (isdigit(c = getchar()));
  return ans * f;
}

const int N = 400050;

int F[N];

int Find(int x) { return F[x] ? F[x] = Find(F[x]) : x; }
bool Union(int x, int y) {
  return (x = Find(x)) != (y = Find(y)) && (F[x] = y, true);
}

int main() {
  int n = readInt(), m = readInt(), ans = n + m - 1;
  for (int q = readInt(); q; --q) {
    int x = readInt(), y = readInt();
    ans -= Union(x, y + n);
  }
  printf("%d\n", ans);
  return 0;
}