/*

ID:

LANG: C++

PROB: ride

*/



#include<bits/stdc++.h>

using namespace std;





#define int                             long long

#define pb                              push_back

#define all(x)                          (x).begin(), (x).end()

#define FastIO                          ios::sync_with_stdio(false);cin.tie(0);

#define clk                             cerr<<endl<<(float)clock()/(float)CLOCKS_PER_SEC<<" sec"<<endl;

#define IO                              freopen("input.txt", "r", stdin); freopen("output.txt", "w", stdout);



#define gcd(a, b)                       __gcd(a, b)

#define lcm(a,b)                        (a/gcd(a,b)*b)

//#include <ext/pb_ds/assoc_container.hpp>

//#include <ext/pb_ds/tree_policy.hpp>

//#define ordered_set tree<pair<int,int>, null_type,less<pair<int,int>>, rb_tree_tag,tree_order_statistics_node_update>

//#define ordered_set tree<int, null_type, less_equal, rb_tree_tag, tree_order_statistics_node_update>

//using namespace __gnu_pbds;





int  comp_double(double a, double b)    {if(fabs(a-b) <= 1e-10)return 0;return a < b ? -1 : 1;}

void setIO(string s)                    {freopen((s + ".in").c_str(), "r", stdin);freopen((s + ".out").c_str(), "w", stdout);}

int dx[8]                               = {2, 2, -2, -2, 1, 1, -1, -1};

int dy[8]                               = {1, -1, 1, -1, -2, 2, 2, -2};



//int dx[]={+1,-1,+0,-0}

//int dy[]={+0,-0,+1,-1}



//memset(memo, -1, sizeof memo);

// index = (index + 1) % n; // index++; if (index >= n) index = 0;

// index = (index + n - 1) % n; // index--; if (index < 0) index = n - 1;

// int ans = (int)((double)d + 0.5); // for rounding to nearest integer

//h = sqrt(x2+y2) double hypot(double x, double y)

//getline(cin >> ws, s2); // eat whitespace



void solve(){

    int n;

    cin >> n;

    for(int i=0; i<n; i++) cout << "(";

    for(int i=0; i<n; i++) cout << ")";

    cout << endl;



    for(int i=0; i<n-1; i++){

        for(int j=0; j<=i; j++){

            cout << "()";

        }

        int nn = n*2 - (i+1)*2;



        for(int j=0; j<nn/2; j++){

            cout << "(";

        }

        for(int j=0; j<nn/2; j++) {

            cout << ")";

        }

        cout << endl;

    }

}

signed main(){

    //IO

    //FastIO

    //setIO("ride");



    int t; cin >> t; while(t--)

    solve();





    return 0;

}