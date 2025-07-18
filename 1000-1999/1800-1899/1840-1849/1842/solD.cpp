#include<bits/stdc++.h>
using namespace std;


//template
#pragma region
#pragma GCC optimize("O3", "unroll-loops") // for fast N^2
#pragma GCC target("avx2", "popcnt") // for fast popcount
//some macros
#define all(x) x.begin(), x.end()
#define lef(x) (x<<1)
#define rig(x) (lef(x)|1)
#define mt make_tuple
#define ft first
#define sd second
#define mp make_pair
#define eb emplace_back
#define pb push_back
//loop that goes (begin,..,end] 
#define rep(i, begin, end) for (__typeof(end) i = (begin) - ((begin) > (end)); i != (end) - ((begin) > (end)); i += 1 - 2 * ((begin) > (end)))
 
typedef long long ll;
typedef pair<ll,ll> pll;
typedef pair<int,int> pii;
typedef pair<double,double> pdb;
 
//read and print pair
template<typename T>
ostream& operator<<(ostream& os, pair<T, T> p) {
    os << "(" << p.first << ", " << p.second << ")";
    return os;
}
template<typename T>
istream& operator>>(istream& is, pair<T, T> &p) {
    is >> p.first >> p.second;
    return is;
}
 
//error debug
#define error(args...) { string _s = #args; replace(_s.begin(), _s.end(), ',', ' '); stringstream _ss(_s); istream_iterator<string> _it(_ss); err(_it, args); }
 
void err(istream_iterator<string> it) {}
template<typename T, typename... Args>
void err(istream_iterator<string> it, T a, Args... args) {
	cerr << *it << " = " << a << endl;
	err(++it, args...);
}

mt19937 rng(chrono::steady_clock::now().time_since_epoch().count());
#pragma endregion

const int N = 110;

ll mat[N][N], temp[N];
int msk[N];
void test()
{
    int n, m;
    memset(mat, -1, sizeof(mat));
    memset(temp, -1, sizeof(temp));
    set<pair<ll,int>> bad;
    cin >> n >> m;
    while(m--)
    {
        int u, v;
        ll y;
        cin >> u >> v >> y;
        mat[u][v] = mat[v][u] = y;
    }
    ll rs = 0;
    queue<int> rmv;
    rep(i, 1, n){
        temp[i] = -1;
        msk[i] = 1;
        if(mat[n][i] != -1)
        {
            temp[i] = mat[n][i];
            if(mat[n][i] == 0){
                msk[i] = 0;
                rmv.push(i);
            }else
                bad.emplace(mat[n][i],i);
        }
    }
    rmv.push(n);
    while(!rmv.empty())
    {
        int at = rmv.front();
        rmv.pop();
        rep(i, 1, n){
            if(mat[at][i] != -1 && (temp[i] == -1 || temp[i] > rs+mat[at][i]))
            {
                if(mat[at][i] == 0){
                    temp[i] = rs+mat[at][i];
                    msk[i] = 0;
                    rmv.push(i);
                }else
                {
                    bad.erase(pair<ll,int>{temp[i],i});
                    temp[i] = rs+mat[at][i];
                    bad.emplace(temp[i],i);
                }
            }
        }
    }
    vector<pair<string,int>>resp;
    while(!bad.empty() && msk[1] != 0)
    {
        ll t;
        int id;
        tie(t,id) = *bad.begin();
        string s;
        // cout << t << " " << id << "<<<<\n";
        rep(i, 1, n+1)
        {
            s.pb(msk[i]+'0');
        }
        resp.eb(s,t-rs);
        bad.erase(bad.begin());
        msk[id] = 0;
        queue<int> rmv;
        rmv.push(id);
        rs = t;
        while(!bad.empty() && bad.begin()->first == t)
        {
            tie(t,id) = *bad.begin();
            bad.erase(bad.begin());
            msk[id] = 0;
            rmv.push(id);
        }
        while(!rmv.empty())
        {
            int at = rmv.front();
            rmv.pop();
            rep(i, 1, n){
                if(mat[at][i] != -1 && (temp[i] == -1 || temp[i] > rs+mat[at][i]))
                {
                    if(mat[at][i] == 0){
                        bad.erase(pair<ll,int>{temp[i],i});
                        temp[i] = rs+mat[at][i];
                        msk[i] = 0;
                        rmv.push(i);
                    }else
                    {
                        bad.erase(pair<ll,int>{temp[i],i});
                        temp[i] = rs+mat[at][i];
                        bad.emplace(temp[i],i);
                    }
                }
            }
        }
    }
    if(msk[1] != 0)
        cout << "inf\n";
    else {
        cout << rs << " " << resp.size() << '\n';
        for(auto [s,qt]:resp)
            cout << s << ' ' << qt << '\n';
    }
}


int main()
{
    ios::sync_with_stdio(false);
    cin.tie(NULL);
    int t = 1;
    // cin >> t;
    while(t--)
        test();    
    return 0;
}
// 1 3 1 3 1 3