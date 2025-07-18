#include <bits/stdc++.h>



using namespace std;

using ll = long long;



const int N = 100;



int n, k;

ll dp[N][N];

ll a[N], b[N];

bool lst[N][N];



bool cmp(int i, int j)

{

    return b[i] < b[j];

}



void solve()

{

    cin >> n >> k;



    vector < int > vec;



    for (int i = 1; i <= n; i++) {

        cin >> a[i] >> b[i];

        vec.push_back(i);

    }



    sort(vec.begin(), vec.end(), cmp);



    for (int i = 0; i <= n; i++) {

        for (int j = 0; j <= k; j++) dp[i][j] = -1e15;

    }



    dp[0][0] = 0;



    for (int i = 1; i <= n; i++) {

        for (int j = 0; j <= k; j++) {

            int p = vec[i - 1];

            lst[i][j] = 0;

            dp[i][j] = dp[i - 1][j] + (k - 1) * b[p];

            if (j && dp[i - 1][j - 1] + b[p] * (j - 1) + a[p] > dp[i][j]) {

                lst[i][j] = 1;

                dp[i][j] = dp[i - 1][j - 1] + b[p] * (j - 1) + a[p];

            }

        }

    }



    vector < int > e, f;



    int p = k;



    for (int i = n; i >= 1; i--) {

        if (!lst[i][p]) {

            f.push_back(-vec[i - 1]);

            f.push_back(vec[i - 1]);

        } else {

            e.push_back(vec[i - 1]);

            p--;

        }

    }



    reverse(e.begin(), e.end());

    reverse(f.begin(), f.end());



    cout << e.size() + f.size() << "\n";



    for (int i = 0; i < k - 1; i++) cout << e[i] << " ";

    for (auto x : f) cout << x << " ";

    cout << e.back() << "\n";

}



int main()

{

    ios_base::sync_with_stdio(0); cin.tie(0); cout.tie(0);



    int t;

    cin >> t;

    while (t--) solve();

}