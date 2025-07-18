#include <bits/stdc++.h>

#include <ext/pb_ds/assoc_container.hpp>

#include <ext/pb_ds/tree_policy.hpp>

using namespace __gnu_pbds;

#define ll long long

#define ld long double

#define el '\n'

#define pi acos(-1)

#define F first

#define S second

#define Baba_Sevawy  ios_base::sync_with_stdio(0);cin.tie(0);cout.tie(0);

#define ordered_set tree<ll , null_type,less<ll>, rb_tree_tag,tree_order_statistics_node_update>

using namespace std;

const ll N = 2e5 + 6, M = 1e6 + 5, mod = 1e9 + 7, K = 21, inf = 2e18;

const ld EPS = 1e-12;

void go() {

    int n, d12, d13, d23;

    cin >> n >> d12 >> d23 >> d13;

    if(d23 == d12 + d13){ // root = 1

        if(d12 + d13 + 1 > n){

            cout << "NO\n";

            return;

        }

        cout << "YES\n";

        int node = 4;

        int second = 1, first;

        for(int i = 1; i <= d12 - 1; i++){

            first = second;

            second = node;

            cout << first << " " << second << el;

            node++;

        }

        cout << second << " " << 2 << el;

        second = 1;

        for(int i = 1; i <= d13 - 1; i++){

            first = second;

            second = node;

            cout << first << " " << second << el;

            node++;

        }

        cout << second << " " << 3 << el;

        while (node <= n) cout << 1 << " " << node++ << el;

    }else if(d12 == d23 + d13){ // root = 3

        if(d23 + d13 + 1 > n){

            cout << "NO\n";

            return;

        }

        cout << "YES\n";

        int node = 4;

        int second = 3, first;

        for(int i = 1; i <= d13 - 1; i++){

            first = second;

            second = node;

            cout << first << " " << second << el;

            node++;

        }

        cout << second << " " << 1 << el;

        second = 3;

        for(int i = 1; i <= d23 - 1; i++){

            first = second;

            second = node;

            cout << first << " " << second << el;

            node++;

        }

        cout << second << " " << 2 << el;

        while (node <= n) cout << 1 << " " << node++ << el;

    }else if(d13 == d12 + d23){ // root = 2

        if(d12 + d23 + 1 > n){

            cout << "NO\n";

            return;

        }

        cout << "YES\n";

        int node = 4;

        int second = 2, first;

        for(int i = 1; i <= d12 - 1; i++){

            first = second;

            second = node;

            cout << first << " " << second << el;

            node++;

        }

        cout << second << " " << 1 << el;

        second = 2;

        for(int i = 1; i <= d23 - 1; i++){

            first = second;

            second = node;

            cout << first << " " << second << el;

            node++;

        }

        cout << second << " " << 3 << el;

        while (node <= n) cout << 1 << " " << node++ << el;

    }else{ // root = 4

        for(int j = 1; j < d12; j++){

            int d41 = j, d42 = d12 - j, d43 = d13 - j;

            if(d43 <= 0 || d42 <= 0 || d42 + d43 != d23 || d41 + d42 + d43 + 1 > n) continue;

            cout << "YES\n";

            int node = 5;

            int second = 4, first;

            for(int i = 1; i <= d41 - 1; i++){

                first = second;

                second = node;

                cout << first << " " << second << el;

                node++;

            }

            cout << second << " " << 1 << el;

            second = 4;

            for(int i = 1; i <= d43 - 1; i++){

                first = second;

                second = node;

                cout << first << " " << second << el;

                node++;

            }

            cout << second << " " << 3 << el;

            second = 4;

            for(int i = 1; i <= d42 - 1; i++){

                first = second;

                second = node;

                cout << first << " " << second << el;

                node++;

            }

            cout << second << " " << 2 << el;

            while (node <= n) cout << 4 << " " << node++ << el;

            return;

        }

        cout << "NO\n";

    }

}



signed main() {

    Baba_Sevawy

    int tt = 1;

    cin >> tt;

    while (tt--)

        go();

    return 0;

}