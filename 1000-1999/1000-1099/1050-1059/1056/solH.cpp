#include <bits/stdc++.h>
using namespace std;
using ull = uint64_t;
using ll = int64_t;
using ld = long double;

const int BL = 550;

const int N = 300003;

vector<int> a[N];
vector<tuple<int, int, int>> sm[N];

int id[N];
int n, q;


void solve() {
    cin >> n >> q;

    for (int i = 0; i < q; ++i) {
        int k;
        cin >> k;
        a[i].resize(k);
        for (int& x: a[i]) {
            cin >> x;
            --x;
        }
    }

    for (int i = 0; i < q; ++i) {
        if (a[i].size() > BL) {
            fill(id, id + n, -1);
            for (int j = 0; j < a[i].size(); ++j) {
                id[a[i][j]] = j;
            }
                
            for (int j = 0; j < q; ++j) {
                if (j == i) {
                    continue;
                }

                int last = n;

                for (int x: a[j]) {
                    int y = id[x];
                    if (y != -1) {
                        if (last < y - 1) {
                            cout << "Human\n";
                            return;
                        }

                        last = y;
                    }
                }
            } 
        } else {
            for (int j = 0; j + 1 < a[i].size(); ++j) {
                sm[a[i][j]].emplace_back(a[i][j + 1], i, j);
            }
        }
    }

    fill(id, id + n, 0);

    int cur = 1;

    for (int i = 0; i < n; ++i) {
        ++cur;
        int lg = cur, rg = cur;
        sort(sm[i].begin(), sm[i].end());
        int prev = -1;
        for (auto& [nx, qid, idx]: sm[i]) {
            ++cur;

            if (nx != prev) {
                rg = cur;
                prev = nx;
            }

            auto& ca = a[qid];
            for (int j = idx + 1; j < ca.size(); ++j) {
                int x = ca[j];
                if (lg <= id[x] && id[x] < rg) {
                    cout << "Human\n";
                    return;
                }

                id[x] = cur;
            }
        }
    }

    cout << "Robot\n";
}

int main() {
#ifdef BZ
    freopen("input.txt", "r", stdin);
#endif
    ios_base::sync_with_stdio(false); cin.tie(nullptr); cout.tie(nullptr); cout.setf(ios::fixed); cout.precision(20);
    int t;
    cin >> t;
    for (int i = 0; i < t; ++i) {
        solve();

        for (int j = 0; j < q; ++j) {
            a[i].clear();
        }   

        for (int j = 0; j < n; ++j) {
            sm[j].clear();
        }
    }
}