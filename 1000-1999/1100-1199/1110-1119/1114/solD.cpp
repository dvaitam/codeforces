#include <iostream>
#include <vector>
#include <algorithm>
#include <set>
#include <iomanip>
#include <stack>
#include <queue>
#include <string>
#include <cmath>
#include <map>
#include <cstdlib>
#include <unordered_set>
#include <unordered_map>
#include <random>

using namespace std;

#define INF 1000000010
#define MOD 1000000007
#define pb push_back
#define f first
#define s second
#define mp make_pair
#define pii pair <int, int>

const long long MAXN = 1e6 + 10;
//const long double eps = 0.0000001;

typedef long long ll;
typedef unsigned long long ull;
typedef long double ld;
typedef vector <int> vi;


int n, a[MAXN], maxx[MAXN], ans;
vi v;

int main() {
    ios_base::sync_with_stdio(0);
    cin.tie(0);
    cout.tie(0);
    //freopen("network.in", "r", stdin);
    //freopen("network.out", "w", stdout);
    cin >> n;
    for(int i = 0; i < n; i++){
        cin >> a[i];
    }
    v.pb(a[0]);
    for(int i = 1; i < n; i++){
        if(a[i] != a[i - 1])
            v.pb(a[i]);
    }
    for(int i = 0; i < v.size(); i++){
        int curmax = 0;
        for(int j = v.size() - 1; j > i; j--){
            int old = maxx[j];
            if(v[i] == v[j]){
                maxx[j] = max(maxx[j], curmax + 1);
            }
            curmax = max(curmax, old);
        }
    }
    for(int i = 0; i < v.size(); i++){
        ans = max(ans, maxx[i]);
    }
    cout << v.size() - 1 - ans;
    
    
    
    return 0;
}