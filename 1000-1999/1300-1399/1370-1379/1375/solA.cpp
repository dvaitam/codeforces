#include <bits/stdc++.h>



using namespace std;

#define ll long long





void solve() {

    int n;

    cin >> n;

    vector<int> a(n);

    for (int i = 0; i < n; ++i) {

        cin >> a[i];

        if(i%2==1){

            a[i] = -abs(a[i]);

        }

        else {

            a[i] = abs(a[i]);

        }

    }

    for (int i = 0; i < n; ++i) {

        cout << a[i] << " ";

    }

    cout << "\n";

}



int main() {



    ios::sync_with_stdio(false);

    cin.tie(nullptr);

    int T;

    cin >> T;

    while (T--) {

        solve();

    }

    return 0;

}