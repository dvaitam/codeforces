#include <cstdio>
#include <queue>
#include <utility>
#include <algorithm>
#ifdef ONLINE_JUDGE
#define freopen(a, b, c)
#endif

typedef long long int ll;

namespace IPT {
  const int L = 1000000;
  char buf[L], *front=buf, *end=buf;
  char GetChar() {
    if (front == end) {
      end = buf + fread(front = buf, 1, L, stdin);
      if (front == end) return -1;
    }
    return *(front++);
  }
}

template <typename T>
inline void qr(T &x) {
  char ch = IPT::GetChar(), lst = ' ';
  while ((ch > '9') || (ch < '0')) lst = ch, ch=IPT::GetChar();
  while ((ch >= '0') && (ch <= '9')) x = (x << 1) + (x << 3) + (ch ^ 48), ch = IPT::GetChar();
  if (lst == '-') x = -x;
}

namespace OPT {
  char buf[120];
}

template <typename T>
inline void qw(T x, const char aft, const bool pt) {
  if (x < 0) {x = -x, putchar('-');}
  int top=0;
  do {OPT::buf[++top] = static_cast<char>(x % 10 + '0');} while (x /= 10);
  while (top) putchar(OPT::buf[top--]);
  if (pt) putchar(aft);
}

const int maxn = 1005;
const int dx[] = {0, 1, 0, -1}, dy[] = {1, 0, -1, 0};
const char pv[] = {'R', 'D', 'L', 'U'};
const char dp[] = {'L', 'U', 'R', 'D'};

int n;
int x[maxn][maxn], y[maxn][maxn];
char ans[maxn][maxn];
std::queue<std::pair<int, int> > Q;

int main() {
  freopen("1.in", "r", stdin);
  qr(n);
  for (int i = 1; i <= n; ++i) {
    for (int j = 1; j <= n; ++j) {
      qr(x[i][j]); qr(y[i][j]);
    }
  }
  for (int i = 1; i <= n; ++i) {
    for (int j = 1; j <= n; ++j) if (x[i][j] == -1) {
      for (int t = 0; t < 4; ++t) if (x[i + dx[t]][j + dy[t]] == -1) {
        ans[i][j] = pv[t];
        break;
      }
    } else if ((x[i][j] == i) && (y[i][j] == j)) {
      ans[i][j] = 'X';
      for (Q.push({i, j}); !Q.empty(); Q.pop()) {
        int u = Q.front().first, v = Q.front().second;
        for (int t = 0; t < 4; ++t) {
          int fx = dx[t] + u, fy = dy[t] + v;
          if (ans[fx][fy]) continue;
          if ((x[fx][fy] == i) && (y[fx][fy] == j)) {
            ans[fx][fy] = dp[t];
            Q.push({fx, fy});
          }
        }
      }
    }
  }
  for (int i = 1; i <= n; ++i) {
    for (int j = 1; j <= n; ++j) if (ans[i][j] == 0) {
      puts("INVALID");
      return 0;
    }
  }
  puts("VALID");
  for (int i = 1; i <= n; ++i) {
    puts(ans[i] + 1);
  }
  return 0;
}