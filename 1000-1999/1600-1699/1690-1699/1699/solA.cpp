//

//  main.cpp

//  CP

//

//  Created by richard_bachman

//

 

#include <bits/stdc++.h>

#include <iostream>

 

const int MOD = 1000000007;

const int INF = 1061109567;

const long long LL_INF = 1000000000000000;

 

#define all(x) (x).begin(), (x).end()

#define sz(s) (int)((s).size())

#define FOR(i,a,b) for(int i=(a); i<=(b); i++)

#define ROF(i,b,a) for(int i=(b); i>=(a); i--)

 

#define ll long long

#define vbool vector<int>

#define vint vector<int>

#define vdouble vector<double>

#define vll vector<long long>

#define vpint vector<pint>

#define vvint vector<vint>

#define vvll vector<vll>

#define pb push_back

 

#define pint pair<int, int>

#define pll pair<long long, long long>

#define mp make_pair

#define f first

#define s second



using namespace std;



template<typename T> istream& operator>>(istream& in, vector<T>& a) {for(auto &x : a) in >> x; return in;};

template<typename T> ostream& operator<<(ostream& out, vector<T>& a) {for(auto &x : a) out << x << ' '; return out;};



template<typename T1, typename T2> ostream& operator<<(ostream& out, const pair<T1, T2>& x) {return out << x.f << ' ' << x.s;}

template<typename T1, typename T2> istream& operator>>(istream& in, pair<T1, T2>& x) {return in >> x.f >> x.s;}



//class ClassName {

//public:

//};



void mod_add(ll &a, ll &b) {

    a = (a + b) % MOD;

}



void mod_sub(ll &a, ll &b) {

    a = (a + MOD - b) % MOD;

}



void solve() {

    ll n; cin >> n;

    

    if (n % 2 == 1) {

        cout << "-1" << "\n";

    } else {

        cout << "0 0 " << n / 2 << "\n";

    }

}



int main() {

    ios_base::sync_with_stdio(false);

    cin.tie(NULL);

    cout.tie(NULL);



    // cout << "good morning, world!" << "\n";

    

    int t = 1;

    cin >> t;

    FOR(i,0,t-1) {

        solve();

    }

    

//    ClassName c;

//    // input = ...

//    cout << c.method(input) << "\n";

    

    return 0;

}



// Try this:

// binary search (not just the answer)

// invariants (properties of array [sum, center of mass, diff-array, xor-array, etc], graph [cycles, path, mincost, etc])

// dp (deterministic states, no re-computing prev. i)

// greedy algorithm