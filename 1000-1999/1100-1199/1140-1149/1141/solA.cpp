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

    int n, m;

    cin >> n >> m;

    int ans = INT_MAX;

    for (int i = 0; i <= 30; i++) {

        for (int j = 0; j <= 30; j++) {

            ll temp = n;

            for (int k = 0; k < i && temp < m; k++) {

                temp *= 2;

            }

            for (int k = 0; k < j && temp < m; k++) {

                temp *= 3;

            }

            if (temp == m) {

                ans = min(ans, i + j);

            }



        }

    }

    if (ans == INT_MAX)ans = -1;

    cout << ans;



}