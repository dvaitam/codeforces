///prod by maxana

#include<bits/stdc++.h>

#pragma GCC optimize("O3")

#pragma GCC optimize("Ofast")

#define F first

#define S second

#define pb push_back

using namespace std;

int32_t main() {

    /*#ifdef LOCAL

        freopen("input.txt", "r", stdin);

        freopen("output.txt", "w", stdout);

    #endif // LOCAL*/

    ios_base::sync_with_stdio(0);cin.tie(0);cout.tie(0);

    int t = 1;

    cin >> t;

    while(t--) {

        int n, k1 = 0, k0 = 0;

        cin >> n;

        for(int i = 0; i < n; ++i){

            int a;

            cin >> a;

            if(a)k1++;

            else k0++;

        }

        if(k1 <= n / 2){

            cout << n - k1 - (n == k0) << '\n';

            k0 -= (n == k0);

            while(k0--)cout << "0 ";

            cout << '\n';

            continue;

        }

        if(k0 < n / 2){

            cout << n - k0 - ((n - k1) % 2) << '\n';

            k1 -= ((n - k1) % 2);

            while(k1--)cout << "1 ";

            cout << '\n';

            continue;

        }

    }

    return 0;

}