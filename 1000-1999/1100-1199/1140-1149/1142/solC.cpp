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

const int N = 100050;

int n;
struct Point {
  LL x, y;
  Point() {}
  Point(LL x, LL y) : x(x), y(y) {}
} C[N], P[N];

inline bool cmp(const Point &a, const Point &b) { return a.x == b.x ? a.y > b.y : a.x < b.x; }
Point operator-(const Point &a, const Point &b) { return Point(a.x - b.x, a.y - b.y); }
inline LL Cross(const Point &a, const Point &b) { return a.x * b.y - a.y * b.x; }

int top;

int main() {
  n = readInt();
  for (int i = 0; i < n; ++i) {
    P[i].x = readInt();
    P[i].y = readInt() - (LL)P[i].x * P[i].x;
  }
  std::sort(P, P + n, cmp);
  top = 0;
  for (int i = 0; i < n; ++i) if (!i || P[i].x != P[i - 1].x) {
    while (top > 1 && Cross(C[top] - C[top - 1], P[i] - C[top]) >= 0) --top;
    C[++top] = P[i];
  }
  printf("%d\n", top - 1);
}