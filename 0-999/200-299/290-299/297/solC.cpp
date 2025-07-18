#include<algorithm>
#include<cassert>
#include<complex>
#include<map>
#include<iomanip>
#include<sstream>
#include<queue>
#include<set>
#include<string>
#include<vector>
#include<iostream>
#include<cstring>
#define FOR(i, a, b) for(int i =(a); i <=(b); ++i)
#define FORD(i, a, b) for(int i = (a); i >= (b); --i)
#define fup FOR
#define fdo FORD
#define REP(i, n) for(int i = 0;i <(n); ++i)
#define VAR(v, i) __typeof(i) v=(i)
#define FORE(i, c) for(VAR(i, (c).begin()); i != (c).end(); ++i)
#define ALL(x) (x).begin(), (x).end()
#define SZ(x) ((int)(x).size())
#define siz SZ
#define CLR(x) memset((x), 0, sizeof(x))
#define PB push_back
#define MP make_pair
#define X first
#define Y second 
#define FI X
#define SE Y
#define SQR(a) ((a)*(a))
#define DEBUG 1
#define debug(x) {if (DEBUG)cerr <<#x <<" = " <<x <<endl; }
#define debugv(x) {if (DEBUG) {cerr <<#x <<" = "; FORE(it, (x)) cerr <<*it <<", "; cout <<endl; }}
using namespace std;

typedef long long LL;
typedef long double LD;
typedef pair<int, int>P;
typedef vector<int>VI;
const int INF=1E9+7;
template<class C> void mini(C&a4, C b4){a4=min(a4, b4); }
template<class C> void maxi(C&a4, C b4){a4=max(a4, b4); }
template<typename T1, typename T2>
ostream& operator<< (ostream &out, pair<T1, T2> pair) { out << "(" << pair.FI << ", " << pair.SE << ")"; return out; }


#define maxn 100005
int co[maxn];
int t[maxn][2];
int res[maxn][2];

int main(){
	ios_base::sync_with_stdio(false);
    int n;
    vector<pair<int, int> > tt;
    cin >> n;
    fup(i, 1, n) {
        int a; 
        a = rand() % 1000;
            //cout << a << endl;
        cin >> a;
        tt.PB(MP(a, i));
    }
    sort(ALL(tt));
    fup(i, 0, siz(tt) - 1) {
        co[i] = tt[i].FI;
    }
    int N = 0;
    while (N < n) N += 3;


    int n3 = N / 3;
    fup(i, 0, n3 - 1) {
        t[i][0] = co[i];
        t[i][1] = 0;
    }
    fup(i, n3, n3 * 2 - 1) {
        t[i][0] = 0;
        t[i][1] = co[i];
    }
    int z = n3 - 1;
    fup(i, n3 * 2, N - 1) {
        t[i][0] = co[i] - z; 
        t[i][1] = z;
        z--;
    }
    /*
    fup(i, 0, N - 1) {
        cout << t[i][0] << " "; 
    }
    cout << endl;

fup(i, 0, N - 1) {
        cout << t[i][1] << " "; 
    }
    cout << endl;
*/


    fup(i, 0, siz(tt) - 1) {
        int p = tt[i].SE;
        res[p][0] = t[i][0];
        res[p][1] = t[i][1];
    }
    cout << "YES" << endl;
    fup(i, 1, n) cout << res[i][0] << " "; cout << endl;
    fup(i, 1, n) cout << res[i][1] << " "; cout << endl;

    


	return 0;
}