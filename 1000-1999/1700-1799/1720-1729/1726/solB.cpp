// The code snippet of Rain Sure

#pragma GCC optimize(2)

#pragma GCC optimize(3)

#include<iostream>

#include<cstring>

#include<algorithm>

#include<vector>

#include<queue>

#include<set>

#include<map>

#include<unordered_map>

#include<unordered_set>

using namespace std;

#define IOS ios::sync_with_stdio(false); cin.tie(0);cout.tie(0);

#define x first

#define y second

#define int long long

#define endl '\n' 

const int inf = 1e9 + 10;

const int maxn = 100010, M = 2000010;

const int mod = 1e9 + 7;

typedef pair<int,int> PII;

typedef long long LL;

typedef unsigned long long ULL;

typedef long double LD;

int h[maxn], e[M], w[M], ne[M], idx;

int dx[4] = {-1, 0, 1, 0}, dy[4] = {0, -1, 0, 1};

int cnt;

void add(int a, int b, int c)

{

    e[idx] = b, w[idx] = c, ne[idx] = h[a], h[a] = idx ++;

}

void add(int a, int b)

{

    e[idx] = b, ne[idx] = h[a], h[a] = idx ++;

}

int qmi(int a,int b){int res=1%mod; a%=mod;while(b) { if(b&1) res=res*a%mod; a=a*a%mod;b>>=1;}return res;}

int gcd(int a,int b) { return b?gcd(b,a%b):a;}

// head

signed main()

{

    IOS;

    int _; cin >> _;

    while(_ -- ) {

        int n, m; cin >> n >> m;

        int v = m / n;

        if(m <= n - 1) {

            cout << "No" << "\n";

            continue;

        }

        if(m % n == 0) {

            cout << "Yes" << "\n";

            for(int i = 1; i <= n; i ++) cout << m / n << " ";

            cout << "\n";

            continue;

        }

        if(n % 2 == 0 && m % 2) {

            cout << "No" << "\n";

            continue;

        }

        if(n % 2) {

            cout << "Yes" << "\n";

            for(int i = 0; i < n - 1; i ++) cout << 1 << ' ';

            cout << m - n + 1 << "\n";

        }else {

            cout << "Yes" << "\n";

            for(int i = 0; i < n - 2; i ++) cout << 1 << ' ';

            cout << (m - (n - 2)) / 2 << ' ' << (m - (n - 2)) / 2 << "\n";

        }

    }

    return 0;

}