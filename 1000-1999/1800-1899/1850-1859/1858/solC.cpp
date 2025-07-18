#include <bits/stdc++.h>
#define cpu() ios::sync_with_stdio(false); cin.tie(nullptr)
#define pb push_back
#define ff first
#define ss second

typedef long long ll;

using namespace std;

const int MOD = 1e9 + 7, MOD1 = (119 << 23) | 1;
const int INF = 1e9 + 5, MAX = 2e5 + 5;

void solve(){
    int n; cin >> n;
    for(int i = 1; i <= n; i += 2){
        for(int j = i; j <= n; j <<= 1){
            cout << j << " ";
        }
    }
    cout << "\n";
}

int main(){
    cpu();
    int t; cin >> t;
    while(t--){
        solve();
    }
    return 0;
}