#include <bits/stdc++.h>

using namespace std;



#define lld long double

#define ll long long

#define inf 0x3f3f3f3f

#define linf 0x3f3f3f3f3f3f3f3fll

#define ull unsigned long long

#define PII pair<int, int>

#define fi first

#define se second

#define mod 1000000007

#define getunique(v) {sort(v.begin(), v.end()); v.erase(unique(v.begin(), v.end()), v.end());}

#define fire ios::sync_with_stdio(false);cin.tie(0);cout.tie(0);

int n, m;

void solve(){

    cin >> n >> m;

    vector<int> u(m + 1), v(m + 1), deg(n + 1);

    for(int i = 1; i <= m; i++){

        cin >> u[i] >> v[i];

        deg[u[i]]++, deg[v[i]]++;

    }

    for(int x = 1; x <= n; x++){

        if(deg[x] < n - 1){

            cout << 2 << '\n';

            for(int i = 1; i <= m; i++){

                cout << ((u[i] == x || v[i] == x) ? 1 : 2) << " \n"[i == m];

            }

            return;

        }

    }

    cout << 3 << '\n';

    for(int i = 1; i <= m; i++){

        int x;

        if(u[i] + v[i] == 3){

            x = 1;

        }

        else if(u[i] == 1 || v[i] == 1){

            x = 2;

        }

        else x = 3;

        cout << x << " \n"[i == m];

    }



}



int main()

{

    fire;

    int t; cin >> t;

    while(t--) solve();

    return 0;

}



/*

    If get stucked, write down something, stay calm and organized

*/