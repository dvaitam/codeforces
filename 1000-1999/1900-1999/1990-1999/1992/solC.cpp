#include <bits/stdc++.h>
using namespace std;

typedef long long ll;
#define all(x) x.begin(), x.end()
#define vecin(name, len) vector<int> name(len); for (auto &_ : name) cin >> _;


void solve() {
    int n, m, k; cin >> n >> m >> k;
    for (int i = n; i >= k; i --)
        cout << i << " ";
    for (int i = m + 1; i < k; i ++)
        cout << i << " ";
    for (int i = 1; i <= m; i ++)
        cout << i << " ";
    cout << endl;
}

int main() {
	ios_base::sync_with_stdio(false);
    cin.tie(NULL);
    int tt = 1;

    cin >> tt;

    while (tt--) solve();
    return 0;
}