#pragma GCC optimize("Ofast")

#ifdef Debugg

#define _GLIBCXX_DEBUG

#endif



#define _USE_MATH_DEFINES



#include <iostream>

#include <vector>

#include <string>

#include <cmath>

#include <set>

#include <map>

#include <queue>

#include <algorithm>

#include <unordered_set>

#include <unordered_map>

#include <ostream>

#include <random>

#include <chrono>

#include <stdlib.h>

#include <cstdio>

#include <iterator>

#include <list>

#include <cmath>



#include <fstream>



#define ll long long

#define ld long double

#define pb push_back

#define se second

#define fi first

#define forn(i, n) for(int i = 0; i < (n); i++)

#define read(a) for(int i = 0; i < a.size(); i++){cin >> a[i];}

#define all(a) a.begin(), a.end()

#define make_unique(a) sort(all(a)); a.resize(unique(a.begin(), a.end()) - a.begin());

#define left v*2 + 1

#define right v*2 + 2

#define mid  (l + r) / 2

using namespace std;







ll mod = /*998244353*/ /*999999937*/ 1e9 + 7;

const int N = 5 * 1e5 + 100;

const int base = 31;



mt19937 rnd(0);



void solve(){

    int n, m;

    cin >> n >> m;

    vector<int> a(n), b(n), c(m);

    vector<int> res(m);

    vector<vector<int>> cnt(n + 1);



    read(a);

    read(b);

    read(c);



    int kal_v = -1;

    for (int i = 0; i < n; i++) {

        if (a[i] != b[i]) cnt[b[i]].pb(i);

        if (a[i] == b[i] && a[i] == c[m - 1]) kal_v = i;

    }

    reverse(all(c));

    int kal = -1;

    forn (i, m){

        if (cnt[c[i]].empty()){

            if (kal == -1 && kal_v == -1){

                cout << "No\n";

                return;

            }

            else if (kal == -1)

                kal = kal_v;



            res[i] = kal;

        }

        else{

            int t = cnt[c[i]].back();

            cnt[c[i]].pop_back();

            if (kal == -1)

                kal = t;

            res[i] = t;

        }

    }



    for (auto it : cnt){

        if (!it.empty()){

            cout << "No\n";

            return;

        }

    }



    cout << "Yes\n";

    reverse(all(res));

    for (auto i : res)

        cout << i + 1 << ' ';

    cout << '\n';

}



int main()

{

    ios_base::sync_with_stdio(false);

    cin.tie(0);





    int t;

    cin >> t;

    while (t--)

        solve();

}







// mark.mark.mar.rkmark