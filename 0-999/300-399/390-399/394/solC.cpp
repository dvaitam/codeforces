//IN THE NAME OF GOD

#include <bits/stdc++.h>

#define F first
#define S second
#define MP make_pair
#define PB push_back
#define foro(i, a, b) for (int i = (int)(a); i < (int)(b); ++i)

using namespace std;
typedef long long ll;
typedef pair < int, int > pii;

const string ans[4] = {"11 ", "10 ", "01 ", "00 "};
int a[4], n, m;

int main() {
    ios_base::sync_with_stdio (0);
    cin >> n >> m;
    foro(i, 0, n)   foro(j, 0, m) {
        int tmp;
        cin >> tmp;
        if (tmp == 11)  ++a[0];
        else if (tmp == 0)  ++a[3];
        else            ++a[1];
    }
    vector < int > p;
    foro(i, 0, n) {
        p.clear ();
        foro(j, 0, m) {
            if (a[0]-- > 0)         p.PB (0);
            else if (a[1]-- > 0)    p.PB (1 + (i % 2));
            else                    p.PB (3);
        }
        if (i % 2)
            foro(i, 0, p.size () )  cout << ans[p[i]];
        else
            foro(i, 0, p.size () )  cout << ans[p[p.size () - i - 1]];
        cout << endl;
    }

    return 0;
}