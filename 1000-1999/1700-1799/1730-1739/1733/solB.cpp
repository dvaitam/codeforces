#include<bits/stdc++.h>

#define ll long long

#define ld long double

#define fi first

#define se second

#define pii pair<int, int>

#define pll pair<long long, long long>

using namespace std;



void solve() {

    int n, x, y;

    cin >> n >> x >> y;

    if (x > 0 && y > 0) {

        cout << -1 << "\n";

        return;

    }



    if (x == 0) swap(x, y);

    if (x == 0) {

        cout << -1 << "\n";

        return;

    }



    if ((n - 1) % x) {

        cout << -1 << "\n";

        return;

    }



    int a = 1, b = 2, cnt = 0;

    for (int i = 0; i < n - 1; i++) {

        if (cnt < x) {

            cout << a << " ";

            b++;

            cnt++;

        }

        else {

            cout << b << " ";

            a = b + 1;

            swap(a, b);

            cnt = 1;

        }

    }



    cout << "\n";

    return;

}



int main() {

    ios_base::sync_with_stdio(false);

    cin.tie(NULL);

    int tt;

    cin >> tt;

    while(tt--) {

        solve();

    }



    return 0;

}