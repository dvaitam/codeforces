#include <bits/stdc++.h>

using namespace std;

const int max_n = 2000;
const int MOD = 1e9 + 7;
const long long INF = 1e18;

#define int long long
#define fori(i, n) for (int i = 0; i < n; ++i)
#define all(v) v.begin(), v.end()
#define AIDAR ASSADULLIN
#define endl "\n"

signed main() {
//    freopen("input.txt", "r", stdin);
//    freopen("output.txt", "w", stdout);
    cin.tie(nullptr), cout.tie(nullptr), ios_base::sync_with_stdio(false), cout.precision(10);
    int qoq;
    cin >> qoq;
//    int qoq = 1;
    while (qoq--) {
        int a, b;
        cin >> a >> b;
        if (b == 1) {
            cout << "NO\n";
            continue;
        }
        int x, y;
        x = a * (b / 2);
        y = a * ((b + 1) / 2);
        if (b % 2 == 0) y += a * b;
        cout << "YES\n" << x << " " << y << " " << x + y << "\n";
    }
}