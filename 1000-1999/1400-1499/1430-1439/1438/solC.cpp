#include<bits/stdc++.h>

using namespace std;

#define int long long

int a[111][111];

void solve(){

    int n, m;

    cin >> n >> m;

    for (int i = 1; i <= n;i++){

        for (int j = 1; j <= m;j++){

            cin >> a[i][j];

        }

    }

    int z = 1;

    for (int i = 1; i <= n;i++){

        z = i % 2;

        for (int j = 1; j <= m;j++){

            if((a[i][j]&1)&&(z&1)){

                z = 1 - z;

                continue;

            }

            if(!(a[i][j]&1)&&!(z&1)){

                z = 1 - z;

                continue;

            }

            a[i][j]++;

            z = 1 - z;

        }

    }

    for (int i = 1; i <= n;i++){

        for (int j = 1; j <= m;j++){

            cout << a[i][j] << " \n"[j == m];

        }

    }

}

signed main(){

    ios::sync_with_stdio(false), cin.tie(0), cout.tie(0);

    int t;

    cin >> t;

    while(t--){

        solve();

    }

}