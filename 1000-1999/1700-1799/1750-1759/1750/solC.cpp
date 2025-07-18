#include <bits/stdc++.h>



#include <ext/pb_ds/assoc_container.hpp>

#include <ext/pb_ds/tree_policy.hpp>



#define x0 abc

#define y0 ABC

#define x1 abcd

#define y1 ABCD

#define xn abcde

#define yn ABCDE



#define lb lower_bound

#define ub upper_bound



#define in insert

#define er erase



#define fi first

#define se second

#define mp make_pair



#define pb push_back

#define pob pop_back



#define uns unsigned



#define ld long double

#define ll long long



#define cont continue

#define re return



#define endl '\n'



#define fbo find_by_order

#define ook order_of_key



#define MAXLL 9000000000000000000LL

#define MAXINT 2000000000



#define MINLL -9000000000000000000LL

#define MININT -2000000000



#define OUT cout << "-------" << endl;



#define bpc __builtin_popcount

#define bpcll __builtin_popcountll



#define tm qwerty



using namespace std;



using namespace __gnu_pbds;



mt19937 rd(chrono::steady_clock::now().time_since_epoch().count());



typedef tree < pair <ll, ll>, null_type, less < pair <ll, ll> >, rb_tree_tag, tree_order_statistics_node_update > ordered_set;







ld pi = acos(-1.), E = exp(1.);



ll N, M = 998244353, D = 5000, T, Q;



vector < pair < int, int > > vec;



char c1[250001], c2[250001];



int n, i, k;



bool fl1, fl2;



    void solve()

    {

        cin >> n;



        for (i = 1; i <= n; i++) cin >> c1[i];



        for (i = 1; i <= n; i++) cin >> c2[i];



        fl1 = fl2 = false;



        for (i = 1; i <= n; i++)

            if (c1[i] == c2[i]) fl1 = true;

            else fl2 = true;



        if (fl1 && fl2) {

            cout << "NO" << endl;



            re;

        }



        vec.clear();



        k = 0;



        for (i = 1; i <= n; i++)

            if (c1[i] == '1') {

                vec.pb({i, i});



                if (i > 1)

                k++;

            }



        if ((c2[1] - '0' + k) % 2 != 0) {

            vec.pb({1, n - 1});

            vec.pb({n, n});

            vec.pb({1, n});

        }



        cout << "YES" << endl;



        cout << vec.size() << endl;



        for (auto e : vec) cout << e.fi << " " << e.se << endl;

    }



int main()

{

//    freopen("input.txt", "r", stdin);

//    freopen("output.txt", "w", stdout);



    ios::sync_with_stdio(0);

    cin.tie(0);



    cin >> T;



    while (T--) solve();



    re 0;

}