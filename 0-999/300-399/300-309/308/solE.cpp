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

const int MX = 2010;
struct sheep {
    int l,r,i;
    sheep(int _l=0, int _r=0, int _i=0): l(_l), r(_r), i(_i) {};

    bool operator < (sheep const & s) const {
        return r < s.r;
    }
};
sheep v[MX];
vector<vector<bool> > T;
int cnt[MX];
int ogr[MX];
int pos[MX];
int on[MX];

bool ok(int n, int med, bool print) {
    REP(i, n) {
        ogr[i] = n-1;
        cnt[i] = 0;
        pos[i] = -1;
    };
    cnt[n-1] = n;
    REP(i, n) {
        int sum  =0;
      //  REP(j, n) cout << ogr[j] << " ";
      //  cout << endl;
        REP(j, n) {
            sum += cnt[j];
            if (sum > j+1) {
      //          cout << med << "NO"<< endl;
                return false;
            }
            if (sum == j+1 && j >= i) break;
        }
        int cur = -1;
        REP(j, n) {
            if (ogr[j] <= sum-1 && pos[j] == -1) {
                cur = j;
                break;
            }
        }
     //   cout << "choose" << cur << endl;
        cnt[ogr[cur]]--;
        cnt[i]++;
        ogr[cur] = i;
        pos[cur] = i;
        on[i] = cur;
        REP(j, n) if (T[cur][j]) {
            if (i+med < ogr[j]) {
                cnt[ogr[j]]--;
                ogr[j] = i+med;
                cnt[ogr[j]]++;
            }
      //      cout << "impact " << j << i+med << endl;
        }
    }
    if (print) {
        REP(i,n)cout << v[on[i]].i+1 << " ";
        cout << endl;
    }
   
   //cout << med << " YES"<< endl;
    return true; 
}


int main(){
    ios_base::sync_with_stdio(false);
    int n;
    cin >> n;    
    REP(i, n) {
        int l,r;
        cin >> l >> r;
        v[i] = sheep(l,r,i);
            
    }
    sort(v,v+n);
    T = vector<vector<bool> >(n, vector<bool>(n,false));
    REP(i, n) REP(j, n) {
        T[i][j] = !(v[i].l > v[j].r || v[j].l > v[i].r);
    }
    int beg = 0, end = n;
    while(beg < end) {
        int med = (beg+ end)/2;
        if (ok(n, med, false)) end = med;
        else beg = med+1;
    } 
    ok(n, beg, true);
    return 0;
}