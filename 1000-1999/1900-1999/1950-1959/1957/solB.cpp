#include <bits/stdc++.h>

using namespace std;

#define int long long

void _()
{
    int n, k;
    cin >> n >> k;
    if (n == 1)
    {
        cout << k << '\n';
        return;
    }
    int x = 0, i = 0;
    while ((x | (1LL << i)) <= k)
    {
        x |= (1LL << i);
        i++;
    }
    cout << x << ' ' << k - x << ' ';
    for (int i = 3; i <= n; i++) cout << 0 << ' ';
    cout << '\n';
}

signed main()
{
    ios_base::sync_with_stdio(0);
    cin.tie(0);
    int t;
    cin >> t;
    while (t--) _();
}