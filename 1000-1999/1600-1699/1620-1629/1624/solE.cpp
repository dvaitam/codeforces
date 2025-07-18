#include <bits/stdc++.h>

#include <ext/pb_ds/assoc_container.hpp>

#include <ext/pb_ds/tree_policy.hpp>

 

using ll = long long;

using ld = long double;

using namespace std;

using namespace __gnu_pbds;

#define ordered_set tree<int, null_type, less<int>, rb_tree_tag, tree_order_statistics_node_update>

#define all(x) (x).begin(), (x).end()

#define rall(x) (x).rbegin(), (x).rend()

#define ft first

#define sd second

 

constexpr ll N = ll(2e5) + 5;

constexpr int MOD = int(1e9) + 7;

constexpr int inf = 0x1f1f1f1f;



string s; 

int at(int x){

    return s[x] - '0';

}

void solve(){

    int n, m;

    cin >> n >> m;



    array<int, 3> a2[10][10]{}, a3[10][10][10]{}, ans[m]{};

    for(int i = 1; i <= n; i++){

        cin >> s;



        if(m > 1) a2[at(0)][at(1)] = {1, 2, i};

        for(int j = 2; j < m; j++){

            int x = at(j - 2), y = at(j - 1), z = at(j);



            a2[y][z] = {j, j + 1, i};

            a3[x][y][z] = {j - 1, j + 1, i};

        }

    }

    cin >> s;

    

    int p[m]{};

    memset(p, -1, sizeof(p));



    if(m > 1) ans[1] = a2[at(0)][at(1)];

    if(m > 2) ans[2] = a3[at(0)][at(1)][at(2)];

    for(int i = 3; i < m; i++){

        array<int, 3> a = a2[at(i - 1)][at(i)];

        if(ans[i - 2][0] && a[0]){

            p[i] = i - 2;

            ans[i] = a;

        }

        else{

            a = a3[at(i - 2)][at(i - 1)][at(i)];

            if(ans[i - 3][0] && a[0]){

                p[i] = i - 3;

                ans[i] = a;

            }

        }

    }

    if(!ans[m - 1][0]){

        cout << "-1\n";

        return;

    }

    int i = m - 1;

    vector<array<int, 3>> v{ans[i]};

    while(p[i] != -1){

        i = p[i];

        v.push_back(ans[i]);

    }

    cout << v.size() << '\n';

    for(i = v.size() - 1; ~i; i--){

        for(auto& j : v[i]) cout << j << ' ';

        cout << '\n';

    } 

}

int main() {

    //freopen("angry.in", "r", stdin);

    //freopen("angry.out", "w", stdout);



    ios::sync_with_stdio(false);

    cin.tie(nullptr);

 

    int tt = 1;

    cin >> tt;

 

    for(int i = 0; i < tt; i++) {

        solve();

    }

    

}