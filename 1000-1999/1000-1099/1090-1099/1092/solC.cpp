#include<bits/stdc++.h>
using namespace std;
#define ll long long
const int N = 1e2 + 5;
bool matchpre(string s, string t) {
    for (int i = 0; i < t.size(); i++) {
        if (t[i] != s[i])
            return 0;
    }
    return 1;
}
bool matchsuf(string s, string t) {
    int j = t.size() - 1;
    for (int i = s.size() - 1; j > -1; j--, i--) {
        if (s[i] != t[j])
            return 0;
    }
    return 1;
}
int main()
{
    int n;
    cin >> n;
    vector<pair<int, string> > v, meh;
    for (int i = 0; i < 2 * n - 2; i++) {
        string s;
        cin >> s;
        v.push_back({s.size(), s});
        meh.push_back({i, s});
    }
    sort(v.rbegin(), v.rend());
    string s1 = v[0].second, s2 = v[1].second;
    s1.push_back(v[1].second.back());
    s2.push_back(v[0].second.back());
    //cout << s1 << ' ' << s2 << '\n';
    string res;
    set<string> ss;
    bool ok = 1;
    res.resize(2 * n - 2);
    for (int i = 0; i < meh.size(); i++) {
        if (matchsuf(s1, meh[i].second)) {
            if (ss.find(meh[i].second) == ss.end()) {
                res[meh[i].first] = 'S';
                ss.insert(meh[i].second);
            }
            else
                res[meh[i].first] = 'P';
        }
        else if (matchpre(s1, meh[i].second)) {
            res[meh[i].first] = 'P';
        }
        else {
            ok = 0;
            break;
        }
    }
    if (ok) {
        cout << res;
    }
    else {
        //cout << "HERE?\n";
        res.clear();
        res.resize(2 * n - 2);
        ss.clear();
        for (int i = 0; i < meh.size(); i++) {
            if (matchsuf(s2, meh[i].second)) {
                if (ss.find(meh[i].second) == ss.end()) {
                    res[meh[i].first] = 'S';
                    ss.insert(meh[i].second);
                }
                else
                    res[meh[i].first] = 'P';
            }
            else if (matchpre(s2, meh[i].second)) {
                res[meh[i].first] = 'P';
            }
            else {
                ok = 0;
                break;
            }
        }
        cout << res;
    }
    return 0;
}

/*#include <bits/stdc++.h>
using namespace std;

char m[6][10] = {'s','s','s','s','s','s','s','s','s','s',
                  's','s','s','s','s','s','s','s','s','s',
                  's','c','s','c','s','s','s','s','s','s',
                  's','s','b','s','s','b','b','s','s','c',
                  'M','b','b','s','b','b','b','b','s','s',
                  'b','b','b','b','b','b','b','b','b','b'
                 };
int n = 3;
vector<pair<int, int> > path;
vector<vector<pair<int, int> > > coinpaths;
bool vis[6][10];
void solve(int x, int y, int coins, int felhawa) {
    if (coins == n) {
        coinpaths.push_back(path);
    }
    //cout << x << ' ' << y << '\n';
    if (x < 0 || y < 0 || x >= 6 || y >= 10)
        return;
    if (vis[x][y] || m[x][y] == 'b')
        return;
    vis[x][y] = 1;
    path.push_back({x, y});
    bool c = 0;
    if (m[x][y] == 'c')
        c = 1, m[x][y] = 's';
    if (m[x + 1][y] == 'b') {
        solve(x, y + 1, coins + c, 0);
        solve(x, y - 1, coins + c, 0);
        solve(x - 1, y, coins + c, 1);
    }
    else {
        if (felhawa == 1) {
            solve(x, y + 1, coins + c, 2);
            solve(x, y - 1, coins + c, 2);
        }
        else {
            vis[x][y] = 0;
            solve(x + 1, y, coins + c, 0);
        }
    }
    if (m[x][y] == 'c')
        m[x][y] = 'c';
    path.pop_back();
    vis[x][y] = 0;
}
void show() {
    for (int i = 0; i < 6; i++) {
        for (int j = 0; j < 10; j++) {
           cout << m[i][j] << ' ';
        }
        cout << '\n';
    }
}
int main ()
{
    memset(vis, 0, sizeof vis);
    for (int i = 0; i < 6; i++) {
        for (int j = 0; j < 10; j++) {
            if (m[i][j] == 'M') {
                solve(i, j, 0, 0);
            }
        }
    }
    cout << coinpaths.size() << '\n';
    int prvx = -1, prvy = -1;

    for (int i = 0; i < coinpaths[1].size(); i++) {
        int x = coinpaths[1][i].first;
        int y = coinpaths[1][i].second;
        if (prvx != -1)
            m[prvx][prvy] = 's';
        m[x][y] = 'M';
        show();
        cout << "------------\n";
        prvx = x, prvy = y;
    }
    return 0;
}
//Test 4
/*char m [6][10] ={'s','s','s','s','s','c','s','s','s','s',
                  's','c','b','b','b','b','s','s','s','s',
                  's','s','s','s','s','s','b','s','s','s',
                  's','b','b','b','b','s','s','b','M','s',
                  'b','b','b','b','b','b','c','b','b','s',
                  'b','b','b','b','b','b','b','b','b','b'} ;*/

// CustomTest
/*char m [6][10] = {'s','s','s','s','s','s','s','s','s','s',
                   's','c','s','s','c','s','b','s','s','b',
                   's','s','b','s','b','b','s','s','b','c',
                   'M','b','s','b','c','s','s','s','s','s',
                   'b','s','s','s','b','b','b','s','s','b',
                   'b','b','b','b','b','b','b','b','b','b'} ;*/

//Test 3
/*char m [6][10] ={'s','s','s','s','s','s','s','s','s','s',
                  's','s','s','s','s','s','s','s','s','s',
                  's','s','c','s','s','s','s','s','s','s',
                  's','s','s','b','b','b','b','s','s','c',
                  'M','s','b','b','b','b','b','b','s','s',
                  'b','b','b','b','b','b','b','b','b','b'} ;*/