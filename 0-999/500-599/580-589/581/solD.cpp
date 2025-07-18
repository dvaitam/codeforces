/*  Made by

    Akylbek Aitkali */

#include <cstring>

#include <vector>

#include <list>

#include <map>

#include <set>

#include <deque>

#include <stack>

#include <bitset>

#include <algorithm>

#include <functional>

#include <numeric>

#include <utility>

#include <sstream>

#include <iostream>

#include <iomanip>

#include <cstdio>

#include <cmath>

#include <queue>

#include <cstdlib>

#include <ctime>

#include <cassert>



#define pii pair<int, int>



#define f first

#define s second



#define pb push_back



#define mp make_pair



using namespace std;



typedef unsigned long long ull;

typedef long long ll;

typedef long double ld;



const int N = 1e3 + 30;

const int NN = 20;

const int TMXN = 1e5 * 8 + 1;

const int INF = 1e9 + 7;



int x, y, xx, yy, xxx, yyy;

int mx;

int ans;

char a[4];



void check (int x, int y, int xx, int yy, int xxx, int yyy) {

    if (x == xx && y + yy == mx && xxx + x == mx && yyy == mx) {

        cout << mx << '\n';

        for (int i = 1; i <= xxx; i++) {

            for (int i = 1; i <= yyy; i++) {

                cout << a[3];

            }

            cout << '\n';

        }

        for (int i = 1; i <= x; i++) {

            for (int j = 1; j <= mx; j++) {

                if (j <= y) {

                    cout << a[1];

                } else {

                    cout << a[2];

                }

            }

            cout << '\n';

        }

        exit(0);

    }

    if (x + xx + xxx == mx && y == mx && yy == mx && yyy == mx) {

        cout << mx << '\n';

        for (int i = 1; i <= mx; i++) {

            for (int j = 1; j <= mx; j++) {

                if (i <= x) {

                    cout << a[1];

                } else if (x + xx >= i) {

                    cout << a[2];

                } else {

                    cout << a[3];

                }

            }

            cout << '\n';

        }

        exit(0);

    }

}



int main() {

    scanf("%d%d%d%d%d%d", &x, &y, &xx, &yy, &xxx, &yyy);

    a[1] = 'A';

    a[2] = 'B';

    a[3] = 'C';

    ans = (x * y) + (xx * yy) + (xxx * yyy);

    mx = max(x, max(xx, max(xxx, max(y, max(yy, yyy)))));

    if (mx * mx != ans) {

        cout << -1;

        return 0;

    }

    int i = 0;

    while (i < 3) {

        i++;

        if (i == 2) {

            swap(xx, xxx);

            swap(yy, yyy);

            swap(a[2], a[3]);

        }

        if (i == 3) {

            swap(x, xxx);

            swap(y, yyy);

            swap(a[1], a[3]);

        }

        int t = 0;

        while (t < 8) {

            t++;

            check(x, y, xx, yy, xxx, yyy);

            if (t % 2 == 1) {

                swap(xxx, yyy);

            }

            if (t % 4 == 2) {

                swap(xx, yy);

            }

            if (t % 8 == 4) {

                swap(x, y);

            }

        }

    }

    cout << -1;

    return 0;

}