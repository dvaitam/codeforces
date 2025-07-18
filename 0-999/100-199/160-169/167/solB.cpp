#include <iostream>
#include <cstdio>

using namespace std;
typedef long double ld;

ld a[201][201], a_[201][201], **a1, **a2, **a3, b[201], b_[201], *b1, *b2, *b3, p1, p0;
int c[201], p[201], kb = 0, ka = 0, ma = 0;

int main()
{
    #ifdef Debug
    freopen("test.in", "r", stdin);
    freopen("test.out", "w", stdout);
    #endif
    int n, l, k;
    cin >> n >> l >> k;
    for (int i = 0; i < n; i++)
        cin >> p[i];
    for (int i = 0; i < n; i++)
        cin >> c[i];
    b[0] = 1;
    b1 = b;
    b2 = b_;
    for (int i0 = 0; i0 < n; i0++)
        if (c[i0] == -1) {
            p1 = p[i0] / (ld)100;
            p0 = (100-p[i0]) / (ld)100;
            b2[0] = 0;
            for (int i = 0; i <= kb; i++) {
                b2[i] += p0*b1[i];
                b2[i+1] = p1*b1[i];
            }
            b3 = b1;
            b1 = b2;
            b2 = b3;
            kb++;
        }
    ma = max(0, kb - k);
    a[0][0] = 1;
    a1 = (ld**)a;
    a2 = (ld**)a_;
    for (int i0 = 0; i0 < n; i0++)
        if (c[i0] != -1) {
            p1 = p[i0] / (ld)100;
            p0 = (100-p[i0]) / (ld)100;
            for (int i = 0; i <= ka; i++)
                for (int j = 0; j <= ma; j++)
                    ((ld (*)[201])a2)[i][j] = 0;
            for (int i = 0; i <= ka; i++)
                for (int j = 0; j <= ma; j++) {
                    ((ld (*)[201])a2)[i][j] += p0*((ld (*)[201])a1)[i][j];
                    ((ld (*)[201])a2)[i+1][min(ma, j+c[i0])] += p1*((ld (*)[201])a1)[i][j];
                }
            a3 = a1;
            a1 = a2;
            a2 = a3;
            ka++;
        }
    ld ans = 0;
    for (int i = 0; i <= ka; i++)
        for (int j = 0; j <= ma; j++)
            for (int t = 0; t <= kb; t++)
                if ((t+i >= l) && (j+k >= t))
                    ans += ((ld (*)[201])a1)[i][j] * b1[t];
    cout.precision(20);
    cout << ans << endl;
    return 0;
}