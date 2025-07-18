#include <algorithm>
#include <cctype>
#include <cstdio>
#include <cstring>
#include <functional>

typedef long long LL;

int readInt() {
  int ans = 0, c, f = 1;
  while (!isdigit(c = getchar()))
    if (c == '-') f = -1;
  do ans = ans * 10 + c - '0';
  while (isdigit(c = getchar()));
  return ans * f;
}

const int N = 500050;

int n, T[N];
LL D[N], S1[N], S2[N], P1[N], P2[N];

// d[i] <= k <=> i >= T[N];

inline LL get(int i, int k) {
  // sum_{j=i...n}min(d[j],k)
  int t = std::max(i, T[k]);
  return (LL)(t - i) * k + S2[t];
}

int main() {
  n = readInt();
  for (int i = 1; i <= n; ++i) D[i] = readInt();
  std::sort(D + 1, D + n + 1, std::greater<int>());
  S1[0] = S2[n + 1] = 0;
  for (int i = 1; i <= n; ++i) S1[i] = S1[i - 1] + D[i];
  for (int i = n; i >= 1; --i) S2[i] = S2[i + 1] + D[i];
  for (int i = 0, j = n + 1; i <= n; ++i) {
    while (j > 1 && D[j - 1] <= i) --j;
    T[i] = j;
  }
  for (int i = 1; i <= n; ++i) {
    P1[i] = S1[i] - (LL)i * (i - 1) - get(i + 1, i);
    P2[i] = (LL)(i + 1) * i + get(i + 1, i + 1) - S1[i];
    if (P1[i] > i) P1[i] = n + 1;
  }
  P1[0] = 0; P2[n + 1] = n + 1;
  for (int i = 1; i <= n; ++i) P1[i] = std::max(P1[i], P1[i - 1]);
  for (int i = n; i >= 1; --i) P2[i] = std::min(P2[i], P2[i + 1]);
  bool ok = false;
  for (int i = S1[n] & 1, j = n + 1; i <= n; i += 2) {
    // before d[j]
    while (j > 1 && i >= D[j - 1]) --j;
    if (P1[j - 1] <= i && P2[j] >= i && (S1[j - 1] + i <= (LL)j * (j - 1) + get(j, j))) {
      printf("%d ", i); ok = true;
    }
  }
  if (!ok) printf("-1");
  return 0;
}