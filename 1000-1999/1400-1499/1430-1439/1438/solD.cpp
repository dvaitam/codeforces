#include <bits/stdc++.h>

using namespace std;

using ll = long long;

using ld = long double;

// https://space.bilibili.com/672328094



int main() {

  cin.tie(nullptr)->sync_with_stdio(false);

  cout << fixed << setprecision(20);

  int n;

  cin >> n;

  int sum = 0;

  vector<int> a(n);

  for (auto& x : a) cin >> x, sum ^= x;

  if (n & 1) {

    cout << "YES\n";

    cout << n - 2 << "\n";

    for (int i = 0; i + 2 < n; i += 2) cout << i + 1 << " " << i + 2 << " " << i + 3 << "\n";

    for (int i = n - 3; i - 2 >= 0; i -= 2)

      cout << i + 1 << " " << i - 1 + 1 << " " << i - 2 + 1 << "\n";

  } else {

    if (sum == 0) {

      cout << "YES\n";

      cout << n - 3 << "\n";

      n--;

      for (int i = 0; i + 2 < n; i += 2) cout << i + 1 << " " << i + 2 << " " << i + 3 << "\n";

      for (int i = n - 3; i - 2 >= 0; i -= 2)

        cout << i + 1 << " " << i - 1 + 1 << " " << i - 2 + 1 << "\n";

    } else {

      cout << "NO\n";

    }

  }

}