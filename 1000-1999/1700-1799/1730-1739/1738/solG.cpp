// LUOGU_RID: 96054987
#include <bits/stdc++.h>

using namespace std;

constexpr int N = 1e3 + 9;

int t, n, k, f[N][N], mx[N][N];

bool vst[N][N];

void solve() {

  cin >> n >> k;

  auto cmax = [](int& x, int y) { x < y && (x = y); };

  for (int i = 1; i <= n; ++i) memset(vst[i] + 1, false, n * sizeof(bool));

  for (int i = 1; i < k; ++i) memset(mx[i] + 1, 0, n * sizeof(int));

  memset(f[n + 1] + 1, 0, (n + 1) * sizeof(int));

  for (int i = 1; i <= n; f[i++][n + 1] = 0)

    for (int j = 1; j <= n; ++j) {

      char c;

      cin >> c;

      f[i][j] = ~c & 1;

    }

  for (int i = n; i; --i)

    for (int j = n; j; --j) {

      cmax(f[i][j] += f[i + 1][j + 1], max(f[i + 1][j], f[i][j + 1]));

      if (f[i][j] == k) {

        cout << "NO" << '\n';

        return;

      }

      if (!mx[f[i][j]][j]) mx[f[i][j]][j] = i;

    }

  for (int i = k - 1; i; --i) {

    for (int j = n - 1; j; --j) cmax(mx[i][j], mx[i][j + 1]);

    int x = n, y = 1;

    while (y <= n && vst[x][y]) ++y;

    while (y <= n) {

      vst[x][y] = true;

      if ((y == n || x != mx[i][y + 1]) && x > 1 && !vst[x - 1][y])

        --x;

      else

        ++y;

    }

  }

  cout << "YES" << '\n';

  for (int i = 1; i <= n; ++i, cout << '\n')

    for (int j = 1; j <= n; ++j) cout << vst[i][j];

}

signed main() {

  cin.tie(nullptr)->sync_with_stdio(false);

  for (cin >> t; t; --t) solve();

  return cout << flush, 0;

}