#include <bits/stdc++.h>

using namespace std;
typedef pair<int, int> pii;
int main()
{
    ios::sync_with_stdio(false), cin.tie(0);
    cout.tie(0);

    int t;
    cin >> t;
    while (t--)
    {
        int n;
        cin >> n;
        n *= 2;
        vector<int> a(n);
        for (int i = 0; i < n; i++)
            cin >> a[i];
        sort(a.begin(), a.end());

        vector<pii> po;
        for (int i = 0; i < n / 2; i++)
        {
            po.push_back({a[i], a[n - 1 - i]});
        }

        int ans = 0;
        for (int i = 0; i < po.size() - 1; i++)
        {
            ans += (abs(po[i].first - po[i + 1].first) + abs(po[i].second - po[i + 1].second));
        }
        cout << ans << endl;
        for (auto x : po)
        {
            cout << x.first << " " << x.second << endl;
        }
    }

    return 0;
}