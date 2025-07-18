#include<bits/stdc++.h>
using namespace std;
using ll = long long;
#define all(x) (x).begin(), (x).end()

void solve() {
    int n;
    cin >> n;

    int sum = 0;
    for (int i = 1; i <= n; i++)
        sum += (i * 2 - 1) * i;

    cout << sum << ' ' << n * 2 - 1 << '\n';

    cout << 1 << ' ' << n << ' ';
    for (int i = 1; i <= n; i++)
        cout << i << " \n"[i == n];

    for (int i = 1; i < n; i++) {
        cout << 2 << ' ' << n - i << ' ';
        for (int j = 1; j <= n; j++)
            cout << j << " \n"[j == n];

        cout << 1 << ' ' << n - i << ' ';
        for (int j = 1; j <= n; j++)
            cout << j << " \n"[j == n];
    }
}

int main() {
    cin.tie(0)->sync_with_stdio(false);

    int t;
    cin >> t;
    
    while (t--)
        solve();

    return 0;
}