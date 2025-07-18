#include <bits/stdc++.h>
#ifdef ALGO
#include "el_psy_congroo.hpp"
#else
#define DUMP(...)
#define CHECK(...)
#endif

namespace {

using LL = long long;
template<typename T, typename U> bool enlarge(T& a, U b) { return a < b ? (a = b, true) : false; }
template<typename T, typename U> bool minify(T& a, U b) { return a > b ? (a = b, true) : false; }

int gcd(int a, int b) { return b == 0 ? a : gcd(b, a % b); }

struct Solver {

  void solve(int ca, std::istream& reader) {
    int n, m;
    reader >> n >> m;
    std::string s;
    reader >> s;
    int B = 0;
    for (char c : s) B += c == '1';
    int L = -1, R = -1;

    // B / n == b / m
    // => b = B / n * m = B / (n / g) * (m / g)
    int g = gcd(n, m);
    if (B % (n / g)) {
      puts("-1");
      return;
    }
    int b = B / (n / g) * (m / g);
    int len = m;

    int cur = 0;
    for (int i = 0; i < len; ++i) {
      cur += s[i] == '1';
    }
    for (int i = 0; i < n; ++i) {
      if (cur == b) {
        if (i + len <= n) {
          printf("1\n");
          printf("%d %d\n", i + 1, i + len);
          return;
        } else {
          L = i;
          R = (i + len - 1) % n;
        }
      }
      cur -= s[i] == '1';
      cur += s[(i + len) % n] == '1';
    }
    if (L == -1) {
      puts("-1");
    } else {
      assert(L > R);
      printf("2\n");
      printf("%d %d\n", 1, R + 1);
      printf("%d %d\n", L + 1, n);
    }
  }
};

}  // namespace

int main() {
  std::ios::sync_with_stdio(false);
  std::cin.tie(nullptr);
  std::istream& reader = std::cin;

  int cas = 1;
  reader >> cas;
  for (int ca = 1; ca <= cas; ++ca) {
    auto solver = std::make_unique<Solver>();
    solver->solve(ca, reader);
  }
}