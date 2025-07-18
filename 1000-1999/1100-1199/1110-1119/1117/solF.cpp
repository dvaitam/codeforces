#include <bits/stdc++.h>

using namespace std;

typedef long long ll;
typedef long double ld;

bool f(int p, const vector<char> &s, const vector<vector<int> > &m) {
    int sz = s.size();
    for (int i = 0; i < sz; i++) {
        if (s[i] == p + 'a') {
            int j = i;
            while (j < sz && s[j] == s[i]) {
                j++;
            }
            if (i == 0 || j == sz) {
                i = j;
                continue;
            }
            int q1 = s[i - 1] - 'a', q2 = s[j] - 'a';
            if (m[q1][q2] != 1) {
                return false;
            }
            i = j;
        } else {
            if (i != 0 && s[i - 1] != p + 'a') {
                int q1 = s[i] - 'a', q2 = s[i - 1] - 'a';
                if (m[q1][q2] != 1) {
                    return false;
                }
            }
            if (i + 1 != sz && s[i + 1] != p + 'a') {
                int q1 = s[i] - 'a', q2 = s[i + 1] - 'a';
                if (m[q1][q2] != 1) {
                    return false;
                }
            }
        }
    }
    return true;
}

int main()
{
    ios_base::sync_with_stdio(0);
    cin.tie(0);
    cout.tie(0);
    int n, p;
    cin >> n >> p;
    string s;
    cin >> s;
    vector<vector<int> > m(p);
    for (int i = 0; i < p; i++) {
        m[i].resize(p);
        for (int j = 0; j < p; j++) {
            cin >> m[i][j];
        }
    }
    vector<char> tmp(n);
    for (int i = 0; i < n; i++) {
        tmp[i] = s[i];
    }
    vector<int> used(p, 0);
    for (int i = 0; i < p; i++) {
        for (int j = 0; j < p; j++) {
            if (!used[j]) {
                if (f(j, tmp, m)) {
                    used[j] = 1;
                    int cnt = 0, sz = tmp.size();
                    for (int k = 0; k < sz; k++) {
                        if (tmp[k] != j + 'a') {
                            tmp[cnt] = tmp[k];
                            cnt++;
                        }
                    }
                    tmp.resize(cnt);
                }
            }
        }
    }
    cout << tmp.size();
    return 0;
}