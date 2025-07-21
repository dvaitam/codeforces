#include <bits/stdc++.h>
using namespace std;
typedef long long ll;
const int MOD = 998244353;

int main(){
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    int t;
    if(!(cin >> t)) return 0;
    while(t--){
        int n;
        cin >> n;
        vector<ll> f(n+2), s(n+2);
        for(int i = 1; i <= n; i++){
            ll x;
            cin >> x;
            f[i] = (s[i - 1] + x * (n - i + 1)) % MOD;
            s[i] = (s[i - 1] + f[i]) % MOD;
        }
        cout << f[n] << "\n";
    }

    return 0;
}
