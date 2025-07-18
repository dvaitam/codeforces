#pragma GCC optimize("Ofast")

#ifdef Debugg

#define _GLIBCXX_DEBUG

#endif



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







const int mod = 1e9 + 7;

const int N = 1e6 * 5 + 10;



mt19937 rnd(0);



void slove(){

    int n;

    cin >> n;

    vector<int> a(n);

    read(a);

    vector<int> cnt(N, 0), num(N, 0);

    vector<pair<int, int>> fresh_popka(N, {0, 0});

    forn (i, n){

        if (cnt[a[i]] == 1){

            fresh_popka[2 * a[i]] = {num[a[i]], i};

        }

        cnt[a[i]]++, num[a[i]] = i;

    }

    int c = 0;

    vector<int> na;

    multiset<int> mr;

    forn (i, N){

        c += cnt[i] / 2;

        forn (fejk, cnt[i] / 2 * 2) {

            if (mr.size() == 4) break;

            mr.insert(i);

        }

        if (cnt[i] > 0) na.pb(i);

    }



    if (c >= 2){

        cout << "Yes\n";

        vector<int> res;

        forn (i, n){

            if (mr.find(a[i]) != mr.end()) {

                res.pb(i);

                mr.erase(mr.find(a[i]));

            }

        }

        if (a[res[0]] == a[res[1]])

            swap(res[0], res[3]);

        for (auto it : res)

            cout << it + 1 << ' ';

        return;

    }



    n = na.size();

    forn (l, n){

        for (int r = l + 1; r < n; r++){

            if (fresh_popka[na[l] + na[r]] != make_pair(0, 0)){

                cout << "Yes\n";

                cout << num[na[l]] + 1 << ' ' << num[na[r]] + 1 << ' ' << fresh_popka[na[l] + na[r]].fi + 1 << ' ' << fresh_popka[na[l] + na[r]].se + 1;

                return;

            }

            fresh_popka[na[l] + na[r]] = make_pair(num[na[l]], num[na[r]]);

        }

    }



    cout << "No\n";

}



int main()

{

    ios_base::sync_with_stdio(false);

    cin.tie(0);

//

//    int t;

//    cin >> t;

//    while (t--)

        slove();

}







// mark.mark.mar.rkmark