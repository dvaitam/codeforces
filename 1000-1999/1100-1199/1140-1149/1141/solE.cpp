#include <bits/stdc++.h>



using namespace std;

#define sz(s) (int)(s.size())

#define all(v) v.begin(),v.end()

#define clr(d, v) memset(d,v,sizeof(d))

#define ll long long





void file() {

    std::ios_base::sync_with_stdio(0);

    cin.tie(NULL);

    cout.tie(NULL);

}



int main() {

    file();

    ll h;

    ll n;

    cin >> h >> n;

    ll sum = 0;

    ll ans = -1;

    vector<ll> v(n);

    ll mx = 0;

    ll H = h;

    for (int i = 0; i < n; i++) {

        cin >> v[i];

        sum -= v[i];

        H += v[i];

        if (H <= 0 && ans == -1) {

            ans = i + 1;

        }

        mx = max(mx, sum);

    }

    if (ans != -1) {

        cout << ans;

    } else if (sum <= 0) {

        cout << -1;

    } else {



        ll full = (h - mx) / sum;

        h -= full * sum;

        ll cnt = full * n;

        for (int i = 0;; i++) {

            h += v[i % n];

            cnt++;

            if (h <= 0) {

                cout << cnt;

                return 0;

            }

        }

    }



}