#include <bits/stdc++.h>

#define ll long long



using namespace std;



void solve()

{

    ll n;

    cin >> n;



    vector<ll> ev;



    for (int i = 0; i < n; ++i) {

        ll a;

        cin >> a;

        if (a % 2) cout << a << ' ';

        else ev.push_back(a);

    }



    for (long long i : ev) {

        cout << i << ' ';

    }

    

    cout << endl;

}



int main()

{

    ios::sync_with_stdio(false);

    cin.tie(nullptr);



    int t;

    cin >> t;

    for (int i = 0; i < t; ++i) {

        solve();

    }



    return 0;

}