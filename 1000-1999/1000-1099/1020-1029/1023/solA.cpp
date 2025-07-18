#include <bits/stdc++.h>

using namespace std;
#define pb push_back
#define mp make_pair
#define all(x) (x).begin(),(x).end()
#define rep(i,n) for(int i=0;i<(n);i++)
constexpr int MOD = 1000000007;
typedef long long ll;
typedef unsigned long long ull;
typedef pair<int,int> pii;
constexpr int dx[] = {1, 0, -1, 0, 1, 1, -1, -1};
constexpr int dy[] = {0, -1, 0, 1, 1, -1, -1, 1};

template <typename T> ostream &operator<<(ostream &os, const vector<T> &vec){os << "["; for (const auto &v : vec) {os << v << ","; } os << "]"; return os; }

void solve() {
    int N, M;
    cin >> N >> M;
    string s, t;
    cin >> s >> t;
    bool flag = false;
    for(int i = 0; i < N; i++) {
        if (s[i] == '*') flag = true;
    }
    if (flag) {
        if (N - 1 > M) {
            cout << "NO" << endl;
            return;
        }
        bool ok = true;
        for(int i = 0; i < N; i++) {
            if (s[i] == '*') break;
            if (s[i] != t[i]) ok = false;
        }
        for(int i = 0; i < N; i++) {
            if (s[N - 1 - i] == '*') break;
            if (s[N - 1 - i] != t[M - 1 - i]) ok = false;
        }
        if (ok) {
            cout << "YES" << endl;
        } else {
            cout << "NO" << endl;
        }
    } else {
        if (s != t) {
            cout << "NO" << endl;
        } else {
            cout << "YES" << endl;
        }
    }

}

int main() {
    std::cin.tie(0);
    std::ios::sync_with_stdio(false);
    cout.setf(ios::fixed);
    cout.precision(16);
    solve();
    return 0;
}