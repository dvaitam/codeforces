#include <bits/stdc++.h>

using namespace std;

#define ln '\n'

#define all(v) v.begin(), v.end()

typedef long long ll;



void solve(){

    int n; cin >> n;

    vector<int>v(n);

    for (int i = 0; i < n; ++i) {

        cin >> v[i];

    }

    int mn = INT_MAX, mx = INT_MIN;

    for (int i = 0; i+1 < n; ++i) {

        if(v[i] < v[i+1]){

            mn = min(mn, (v[i]+v[i+1])/2);

        }

        if(v[i] > v[i+1]){

            mx = max(mx, (v[i]+v[i+1]+1)/2);

        }

    }

    if(mx == INT_MIN){

        cout << 0 << ln;

    }

    else if(mn == INT_MAX){

        cout << v[0] << ln;

    }

    else{

        if(mx <= mn){

            cout << mx << ln;

        }

        else{

            cout << -1 << ln;

        }

    }

}



int main() {

    ios_base::sync_with_stdio(false);

    cin.tie(NULL);

    cout.tie(NULL);

    int t; cin >> t;

    while (t--){

        solve();

    }

    return 0;

}