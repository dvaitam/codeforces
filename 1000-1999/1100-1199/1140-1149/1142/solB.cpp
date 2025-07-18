#include <algorithm>
#include <cctype>
#include <cstdio>
#include <cstring>

typedef long long LL;

int readInt() {
  int ans = 0, c, f = 1;
  while (!isdigit(c = getchar()))
    if (c == '-') f *= -1;
  do ans = ans * 10 + c - '0';
  while (isdigit(c = getchar()));
  return ans * f;
}

const int N = 200050;

int n, m, q, p[N], T[N], ls[N], rb[N], Q[N], left[N];

void dfs(int x, int dep) {
  Q[dep] = x;
  left[x] = dep >= n ? Q[dep - n + 1] : -1;
  for (int y = ls[x]; y; y = rb[y])
    dfs(y, dep + 1);
}

int main() {
  n = readInt();
  m = readInt();
  q = readInt();
  for (int i = 1; i <= n; ++i) p[readInt()] = i;
  for (int i = 1; i <= m; ++i) {
    int x = p[readInt()];
    int y = T[x == 1 ? n : x - 1];
    rb[i] = ls[y]; ls[y] = i;
    T[x] = i;
  }
  dfs(0, 0);
  for (int i = 1; i <= m; ++i) left[i] = std::max(left[i], left[i - 1]);
  while (q--) {
    int l = readInt(), r = readInt();
    printf("%d", left[r] >= l);
  }
}