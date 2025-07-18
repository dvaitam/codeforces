/**░░░░░░░░░▄▄▄▄▄▀▀▀██▀▀▀▀▀▀▀▄▄▄▄▄░░░░░░░░░

       ░░░░░▄▄▀▀▀░░▄▄▄▄░██░░░░░░░░░░░▀▀▀▄▄░░░░░

       ░░▄█▀░░▄▄██████▀░███░░░░░░▄████▄▄░░▀█▄░░

       ░█▀░░▄█████████░███▀░░░░░▄████████▄░░▀█░

       █▀░░███████████░███░░░░▄████████████░░░█

       █▄░░███████████░█▀░░░░▄█████████████░░▄█

       ░█▄░░▀█████████░░░░▄██▀░░█████████▀░░▄█░

       ░░▀█▄░░▀▀██████▄██████░████████▀▀░░▄█▀░░

       ░░░░░▀▀▄▄▄░░▀▀▀▀▄▄▄▄▄▄▄▀▀▀▀▀░░▄▄▄▀▀░░░░░

       ░░░░░░░░░▀▀▀▀▀▄▄▄▄▄▄▄▄▄▄▄▄▀▀▀▀▀░░░░░░░░░**/

    #pragma GCC optimize ("Ofast,unroll-loops")

    #include <ext/pb_ds/assoc_container.hpp>

    #include <ext/pb_ds/tree_policy.hpp>

    #define pii pair <int , int>

    #include <bits/stdc++.h>

    #define sz size

    #define ll long long

    #define pb push_back

    #define _ << ' ' <<

    #define ve vector

    #define endl '\n'

    #define S second

    #define F first

    #define ld long double

    #define int ll

    using namespace std;

    using namespace __gnu_pbds;



    mt19937 gen (chrono::system_clock().now().time_since_epoch().count());



    typedef tree <ll , null_type ,less<ll> , rb_tree_tag , tree_order_statistics_node_update> oset;



    const int N = int32_t(3e5) + 300;

    const int mod = 9+7;

    int32_t main()

    {

        #ifdef LOCAL

            freopen("input.txt" , "r" , stdin);

            freopen("output.txt" , "w" , stdout);

        #endif

        ios_base::sync_with_stdio(0);

        cin.tie(0); cout.tie(0);

        int t;

        cin >> t;

        while (t--) {

            int n;

            cin >> n;

            char c[n][n];

            for (int i=0;i<n;i++) {

                for (int j=0;j<n;j++) {

                    cin >> c[i][j];

                }

            }

            vector<pii> ans;

            int k0=bool(c[0][1]=='0') + bool(c[1][0]=='0');

            int k1=2-k0;

            if (c[n-1][n-2]!=c[n-2][n-1]) {

                if (k1<k0) {

                    if (c[n-2][n-1]=='1') {

                ans.pb({n,n-1});

                c[n-1][n-2]=c[n-2][n-1];

                    }

                    else {

                        ans.pb({n-1,n});

                        c[n-2][n-1]=c[n-1][n-2];

                    }

                }

                else {

                if (c[n-2][n-1]=='0') {

                ans.pb({n,n-1});

                c[n-1][n-2]=c[n-2][n-1];

                    }

                    else {

                        ans.pb({n-1,n});

                        c[n-2][n-1]=c[n-1][n-2];

                    }

                }

            }

            if (c[n-2][n-1]=='1') {

                if (c[0][1]=='1') {

                    ans.pb({1,2});

                }

                if (c[1][0]=='1') {

                    ans.pb({2,1});

                }

            }

            else {

                if (c[0][1]=='0') {

                    ans.pb({1,2});

                }

                if (c[1][0]=='0') {

                    ans.pb({2,1});

                }

            }

            cout << ans.sz() << endl;

            for (auto to:ans) cout << to.F << " " << to.S << endl;

        }

        return 0;

    }