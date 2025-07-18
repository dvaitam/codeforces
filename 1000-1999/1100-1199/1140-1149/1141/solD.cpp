#include <bits/stdc++.h>

#define ll long long

#define ull unsigned long long

#define ff first

#define sc second

#define pb push_back

#define pii pair<int, int>

#define pll pair<ll, ll>

using namespace std;

const ll mod = 1e9 + 7;

 

void solve(){

    int n;

    cin >> n;

    string a, b;

    cin >> a >> b;

    vector <pii> ans;

    unordered_map <char, vector <int>> A;

    unordered_map <char, vector <int>> B;

    for(int i = 0; i < n; i++){

        A[a[i]].pb(i + 1);

        B[b[i]].pb(i + 1);

    }

    for(char i = 'a'; i <= 'z'; i++){

        while(A[i].size() and B[i].size()){

            ans.pb({A[i].back(), B[i].back()});

            A[i].pop_back();B[i].pop_back();

        }

        while(A[i].size() and B['?'].size()){

            ans.pb({A[i].back(), B['?'].back()});

            A[i].pop_back();B['?'].pop_back();

        }

        while(A['?'].size() and B[i].size()){

            ans.pb({A['?'].back(), B[i].back()});

            A['?'].pop_back();B[i].pop_back();

        }

    }

    while(A['?'].size() and B['?'].size()){

        ans.pb({A['?'].back(), B['?'].back()});

        A['?'].pop_back();B['?'].pop_back();

    }

    cout << ans.size() << '\n';

    for(int i = 0; i < ans.size(); i++){

        cout << ans[i].ff << ' ' << ans[i].sc << '\n';

    }

}

 

int main() {

    ios_base::sync_with_stdio(false);

    cin.tie(nullptr); cout.tie(nullptr);

    

    int T = 1;

    //cin >> T;

    while(T--){

        solve();

        cout << '\n';

    }

    

    return 0;

}