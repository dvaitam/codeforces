#include <bits/stdc++.h>
using namespace std;

//------------------------- DEBUGGING -----------------------------------------------------------------------
#ifdef DEBUG
#define debug( args... ) { vector<string> _v = split( #args, ','); err( _v.begin(), args); cerr << '\n'; }
vector<string> split( const string& s, char c) {
    vector<string> v;
    stringstream ss( s);
    string x;
    while ( getline( ss, x, c))
        v.emplace_back( x);
    return move( v);
}
void err( vector<string>::iterator it) { }
template<typename T, typename... Args>
void err( vector<string>::iterator it, T a, Args... args) {
    cerr << it -> substr( ( *it)[0] == ' ', it -> length()) << " = " << a << '\t';
    err( ++it, args...);
}
#else
#define debug(...)
#endif
//------------------------ MACROS ----------------------------------------------------------------------------

#define pb push_back
#define eb emplace_back
#define fi first
#define se second
#define fr(i,j,k) for(i = j; i < (k); i++)
#define bck(i,j,k) for(i = j; i > (k); i--)
#define all(x) x.begin(), x.end()
#define el '\n'

typedef long long int ll;
typedef long double DO;
typedef pair<int,int> pii;
typedef pair<ll, ll> pll;
typedef vector<int> vi;
typedef vector<bool> vb;
typedef vector<vi> vvi;
//---------------------------  TEMPLATE ENDS, MAKE CHANGES FROM HERE ----------------------------------------

const int mod = 1e9+7;
const ll inf = 2e18;
const DO eps = 1e-9;
const int NN = 1e5 + 2;

int M, N;
string seg[NN];

pii dxy[] = { {1,1} , {1, -1}, {1, 0}, {0, 1} , {0,-1}, {-1, -1} , {-1, 0}, {-1, 1}};

bool ch(int x, int y) {
    if(0 < x && x < N && 0 < y  && y < M) {
        for(auto it : dxy) {
            if(seg[x + it.fi][y + it.se] != '#')
                return 0;
        }
        return 1;
    }
    return 0;
}

bool check(int i, int j) {
    for( auto &p : dxy ) {
        int xx = p.fi + i;
        int yy = p.se + j;
        if(ch(xx, yy)) {
            return 1;
        }
    }
    return 0;
}

void solve(){
    int i = 0, j = 0, k = 0, n = 0, m = 0;
    cin >> N >> M;
    fr(i, 0, N) {
        cin >> seg[i];
    }
    
    for(i=0; i<N; i++) {
        for(j=0; j < M; j++) {

            if(seg[i][j] == '#') {
                if(!check(i, j))  {
                    cout << "NO\n";
                    return;
                }
            }
        }
    }
    cout << "YES\n";

}

int main(){
#ifndef DEBUG
    ios::sync_with_stdio(false); cin.tie(0); 
#endif
    int T = 1, tc;
    for(tc = 1; tc <= T; tc++){
        solve();
    }
    return 0;
}