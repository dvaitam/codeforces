#include <bits/stdc++.h>

using namespace std;

typedef long long ll;

typedef vector<int> vi;

typedef vector<long long> vll;

typedef pair<int, int> pii;

typedef pair<ll, ll> pll;

typedef pair<double, double> pdd;

const double PI = acos(-1);

const ll mod7 = 1e9 + 7;

const ll mod9 = 998244353;

const ll INF = 2 * 1024 * 1024 * 1023;

const char nl = '\n';

#define forn(i, n) for (int i = 0; i < int(n); i++)

#define all(v) v.begin(),v.end()



#pragma GCC optimize("O2")

#pragma GCC target("avx2,bmi,bmi2,lzcnt,popcnt")

 

ll l, r, k, n, m, p, q, u, v, w, x, y, z;

string s, t;

 

////////////////////  LIBRARY CODE ////////////////////



///////////////////////////////////////////////////////





bool multiTest = 0;

void solve(int tt){



    cin >> n;

    vector<pii> ans;



    int numPairs = n/2;

    int gap = 1;

    int num = 1;



    vi leftOver;

    if(n % 2) leftOver.push_back(1);



    while(numPairs) {

        vector<int> changed(n+1, 0);

        int ptr = 1;

        forn(i, numPairs * num) {

            while(changed[ptr]) ptr++;

            ans.push_back({ptr, ptr + gap});

            changed[ptr] = 1;

            changed[ptr + gap] = 1;

        }

        if(numPairs % 2) leftOver.push_back(num * 2);

        numPairs /= 2;

        gap *= 2;

        num *= 2;

    }



    leftOver.pop_back();



    if(leftOver.size() < 2) {

        cout << ans.size() << nl;

        for(auto [y, z] : ans) cout << y << ' ' << z << nl;

        return;

    }



    vi cur;

    forn(i, leftOver[0]) cur.push_back(n - i);



    int ptr = 1;

    int backPtr = n - leftOver[0];

    int leftOverPointer = 1;



    while(leftOverPointer != leftOver.size()) {

        vi newCur;

        if(leftOver[leftOverPointer] == cur.size()) {

            for(int z : cur) {

                ans.push_back({z, backPtr});

                newCur.push_back(backPtr);

                backPtr--;

            }

            leftOverPointer++;

        }

        else {

            for(int z : cur) {

                ans.push_back({z, ptr});

                newCur.push_back(ptr);

                ptr++;

            }

        }

        for(int z : newCur) cur.push_back(z);

    }



    cout << ans.size() << nl;



    for(auto [y, z] : ans) cout << y << ' ' << z << nl;



}



int main() {



    ios::sync_with_stdio(0);

    cin.tie(0);

    cout.tie(0);



    #ifdef yangster

    freopen("input.txt", "r", stdin);

    freopen("output.txt", "w", stdout);

    // #else 

    // freopen("cowpatibility.in", "r", stdin);

    // freopen("cowpatibility.out", "w", stdout);

    #endif



    int t = 1;

    if (multiTest) cin >> t;

    for (int ii = 0; ii < t; ii++) {

        solve(ii + 1);

    }

}