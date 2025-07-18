/****************************************************

         ❤ Bsmellah El Rahman El Rahim ❤

****************************************************/



#include <bits/stdc++.h>



using namespace std;

#define int long long

#define fast ios_base::sync_with_stdio(0),cin.tie(0),cout.tie(0)

#define vi vector <int>

#define pii pair <int,int>

#define vpii vector <pii>

#define vc vector <char>

#define vs vector <string>

#define mii map <int,int>

#define si set <int>

#define ull unsigned long long

#define loop(i, from, to) for(int i = from; i < to; i++)

#define U unsigned int

#define endl "\n"

#define MN 300005

#define INF 100000000000ll

#define all(v) v.begin(),v.end()

#define EPS 1e-6

#define sc2(n, m) cin >> n >> m;

#define sc(n) cin >> n;

#define di deque<int>

#define pq priority_queue <int,vi,greater<>>

clock_t startTime;



double getCurrentTime() {

    return (double) (clock() - startTime) / CLOCKS_PER_SEC;

}



bool isG = 0;





void solve() {

    int n;

    cin >> n;

    vi v(n);

    for (int i = 0; i < n; ++i)

        cin >> v[i];

    std::sort(v.begin(), v.end());

    vi ans(n);

    int ptr = n - 1;

    for (int i = 1; i < n; i += 2) {

        ans[i] = v[ptr--];

    }

    for (int i = 0; i < n; i += 2) {

        ans[i] = v[ptr--];

    }



    for (int i = 0; i < n; ++i)

        cout << ans[i] << ' ';

}



signed main() {

//  =============================================================================

    fast;

#ifndef ONLINE_JUDGE

    freopen("input.txt", "r", stdin);

    freopen("output.txt", "w", stdout);

    freopen("error.txt", "w", stderr);

#endif

//  =============================================================================

    startTime = clock();

//    int TC;

//    cin >> TC;

//    while(TC--){

    solve();

//    }

    return 0;

}