#include <bits/stdc++.h>



typedef long long ll;

#define fast ios::sync_with_stdio(false);cin.tie(nullptr);

using namespace std;



void sol() {

    int n;

    cin >> n;

    vector<int> ans(n);

    int x = 0;

    while (x*x <= n) {

        x++;

    }

    int i = n-1;

    while (i >= 0) {

        int nx = n;

        while (i >= 0 && x*x - i < n) {

            ans[i] = x*x - i;

            nx--;

            i--;

        }

        n = nx;

        x--;

    }

    for (auto j : ans) cout << j << ' ';

    cout << '\n';

}



signed main() {

    fast;

    int t;

    cin >> t;

    while (t--)

        sol();

}