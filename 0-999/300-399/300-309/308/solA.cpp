#include<stack>
#include<algorithm>
#include<cassert>
#include<complex>
#include<map>
#include<iomanip>
#include<sstream>
#include<queue>
#include<set>
#include<string>
#include<cstdio>
#include<vector>
#include<list>
#include<iostream>
#include<cstring>
#include<unordered_set>
#include<unordered_map>

#define FOR(i, a, b) for(int i =(a); i <=(b); ++i)
#define FORD(i, a, b) for(int i = (a); i >= (b); --i)
#define REP(i, n) for(int i = 0;i <(n); ++i)
#define VAR(v, i) __typeof(i) v=(i)
#define FORE(i, c) for(VAR(i, (c).begin()); i != (c).end(); ++i)
#define ALL(x) (x).begin(), (x).end()
#define SZ(x) ((int)(x).size())
#define CLR(x) memset((x), 0, sizeof(x))
#define PB push_back
#define MP make_pair
#define X first
#define Y second 
#define SQR(a) ((a)*(a))
#define DEBUG 1
#define debug(x) {if (DEBUG) cerr << #x <<" = " << x <<endl; }
#define debugv(x) {if (DEBUG) {cerr << #x <<" = "; FORE(__it, (x)) cerr << *__it << ", "; cout << endl; }}
using namespace std;
typedef long long LL;
typedef long double LD;
typedef pair<int, int> P;
typedef vector<int> VI;
const int INF = 1E9+7;
template<class C> void mini(C& a4, C b4){ a4=min(a4, b4); }
template<class C> void maxi(C& a4, C b4){ a4=max(a4, b4); }

const int MX = 1000100;
int a[MX];
int main(){
	ios_base::sync_with_stdio(false);
    LL n,l,t;
    cin >> n >> l >> t;	
    REP(i, n) {
        cin >> a[i];
    }
    t *= 2;
    LL res = 0;
    LL full = t / l;
    LL left = t % l;
    int j = 0;
    bool phase = false;
    REP(i, n) {
        if (left + a[i] >= l) {
            phase = true;
            j = 0;
            left -= l;
        }
        while(j < n && a[j] <= left+a[i]) {
            ++j;
        }
        if (!phase) res += (full+1)*(j-i-1) + full*(n-j+i);
        else res += (full+1)*(n+j-i-1) + full*(i-j);
    }
    cout << setprecision(6) << fixed << ((LD) res/4.0) << endl;
    return 0;
}