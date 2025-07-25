/* Sat 23 Feb 2013 05:56:02 PM ICT */
#include<cstdio>
#include<cstring>
#include<iostream>
using namespace std;

#define REP(i, n) for(int i = 0; i < (n); ++i)

const int L = 7;
struct state {
    int i, j, k;
    bool r, ax, ay, az;

    state() {}
    state(int i, int j, int k, bool r, bool ax, bool ay, bool az):
        i(i), j(j), k(k), r(r), ax(ax), ay(ay), az(az) {}
    inline bool operator == (const state &o) const {
        return i == o.i && j == o.j && k == o.k && r == o.r;
    }
} t[L+1][L+1][L+1][2][2][2][2];
int f[L+1][L+1][L+1][2][2][2][2];
char x[L+1], y[L+1], z[L+1];

bool min(int &a, int b) {
    if(b < a) {
        a = b;
        return 1;
    }
    return 0;
}

void enter() {
    int a, b, c; scanf("%d+%d=%d",&a,&b,&c);
    sprintf(x+1,"%d",a); sprintf(y+1,"%d",b); sprintf(z+1,"%d",c);
}

void solve() {
    REP(ax, 2) REP(ay, 2) REP(az, 2) {
        f[0][0][0][1][ax][ay][az] = (ax | ay | az) == 0 ? 1 << 30 : 1;
        t[0][0][0][1][ax][ay][az] = state(0, 0, 0, 0, 0, 0, 0);
    }
    int lx = strlen(x+1), ly = strlen(y+1), lz = strlen(z+1);
    REP(i, lx) x[i+1] -= '0'; REP(i, ly) y[i+1] -= '0'; REP(i, lz) z[i+1] -= '0';
    REP(i, lx+1) REP(j, ly+1) REP(k, lz+1) if(i | j | k) REP(r, 2) REP(ax, 2) REP(ay, 2) REP(az, 2) {
        int &ff = f[i][j][k][r][ax][ay][az] = 1 << 30;
        state &tt = t[i][j][k][r][ax][ay][az];
        if((x[i]+y[j]+r)%10 == z[k] && k)
            REP(pax, ax+1) REP(pay, ay+1) REP(paz, az+1)
                if(min(ff, f[i-!!i][j-!!j][k-1][x[i]+y[j]+r>9][pax][pay][paz]+(!i&pax)+(!j&pay)))
                    tt = state(i-!!i, j-!!j, k-1, x[i]+y[j]+r>9, pax, pay, paz);
        if(ax) REP(pax, ax+1) REP(pay, ay+1) REP(paz, az+1)
            if(min(ff, f[i][j-!!j][k-!!k][z[k]-y[j]-r<0][pax][pay][paz]+(!j&pay)+(!k&paz)+1))
                tt = state(i, j-!!j, k-!!k, z[k]-y[j]-r<0, pax, pay, paz);
        if(ay) REP(pax, ax+1) REP(pay, ay+1) REP(paz, az+1)
            if(min(ff, f[i-!!i][j][k-!!k][z[k]-x[i]-r<0][pax][pay][paz]+(!i&pax)+(!k&paz)+1))
                tt = state(i-!!i, j, k-!!k, z[k]-x[i]-r<0, pax, pay, paz);
        if(az) REP(pax, ax+1) REP(pay, ay+1) REP(paz, az+1)
            if(min(ff, f[i-!!i][j-!!j][k][x[i]+y[j]+r>9][pax][pay][paz]+(!i&pax)+(!j&pay)+1))
                tt = state(i-!!i, j-!!j, k, x[i]+y[j]+r>9, pax, pay, paz);
    }
    long long pw = 1, a = 0, b = 0, c = 0;
    for(int i=lx, j=ly, k=lz, r=0, ax=1, ay=1, az=1; i | j | k | r; pw *= 10) {
        state &tt = t[i][j][k][r][ax][ay][az];
        if(!(i | j | k)) c += pw;
        else if((x[i]+y[j]+r)%10 == z[k] && k && state(i-!!i,j-!!j,k-1,x[i]+y[j]+r>9,0,0,0) == tt) 
            a += x[i]*pw, b += y[j] * pw, c += z[k] * pw;
        else if(ax && state(i, j-!!j, k-!!k, z[k]-y[j]-r<0, 0, 0, 0) == tt)
            a += ((z[k]-y[j]-r)%10+10)%10 * pw, b += y[j] * pw, c += z[k] * pw;
        else if(ay && state(i-!!i, j, k-!!k, z[k]-x[i]-r<0, 0, 0, 0) == tt)
            a += x[i] * pw, b += ((z[k]-x[i]-r)%10+10)%10 * pw, c += z[k] * pw;
        else
            a += x[i] * pw, b += y[j] * pw, c += (x[i]+y[j]+r)%10 * pw;
        i = tt.i; j = tt.j; k = tt.k; r = tt.r; ax = tt.ax; ay = tt.ay; az = tt.az;
    }
    cout << a << '+' << b << '=' << c << "\n";
}

int main() {
#ifndef ONLINE_JUDGE
    freopen("18.inp", "r", stdin);
    freopen("18.out", "w", stdout);
#endif
    enter();
    solve();
    return 0;
}