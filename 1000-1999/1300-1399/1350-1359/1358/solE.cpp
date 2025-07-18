#include <iostream>
#include <numeric>

using namespace std;

constexpr int N = 2.5e5;

int a[N];

int main() {
  cin.tie(0)->sync_with_stdio(0);
  int n; cin >> n;
  for (int i = 0; i * 2 < n; ++i) {
    cin >> a[i];
  }
  int64_t x; cin >> x;
  int64_t s = accumulate(a, a + (n + 1) / 2, 0ll);
  if (s + n / 2 * x <= 0 && x >= 0) {
    cout << -1;
    return 0;
  }
  int le = n;
  for (int i = 0; le * 2 > n && i + le <= n; ++i) {
    while (le * 2 > n && s + x * (le + ~n / 2 + i) <= 0) --le;
    s -= a[i];
  }
  cout << (le * 2 > n? le: -1);
}