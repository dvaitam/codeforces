#include <bits/stdc++.h>



#define ll long long



#define SZ(x) (int)(x).size()

#define pb push_back



template<class T>inline void chkmax(T &x, const T &y) {if(x < y) x = y;}

template<class T>inline void chkmin(T &x, const T &y) {if(x > y) x = y;}



template<class T>

inline void read(T &x) {

    char c;int f = 1;x = 0;

    while(((c=getchar()) < '0' || c > '9') && c != '-');

    if(c == '-') f = -1;else x = c-'0';

    while((c=getchar()) >= '0' && c <= '9') x = x*10+c-'0';

    x *= f;

}

static int outn;

static char out[(int)2e7];

template<class T>

inline void write(T x) {

    if(x < 0) out[outn++] = '-', x = -x;

    if(x) {

        static int tmpn;

        static char tmp[20];

        tmpn = 0;

        while(x) tmp[tmpn++] = x%10+'0', x /= 10;

        while(tmpn) out[outn++] = tmp[--tmpn];

    }

    else out[outn++] = '0';

}



const int N = 80;



int n;

int L[N+9], R[N+9];

int v[N*2+9], vn;



int p[N+9], pn, pre;

double in[N+9], small[N+9];



double f[N+9][N+9][N+9];

double g[N+9][N+9];



int cnt;

double b[7][N+9][N+9];



double ans[N+9][N+9];



inline void backup(int d) {

    for(int i = 0; i < pn; ++i)

        for(int j = 0; i+j < pn; ++j)

            b[d][i][j] = g[i][j];

}



inline void recover(int d) {

    for(int i = 0; i < pn; ++i)

        for(int j = 0; i+j < pn; ++j)

            g[i][j] = b[d][i][j];

}



inline void ins(int l, int r) {

    for(int i = l; i <= r; ++i) {

        double x = in[i], y = small[i];

        for(int a = cnt; a >= 0; --a)

            for(int b = cnt-a; b >= 0; --b) {

                g[a+1][b] += g[a][b]*y;

                g[a][b+1] += g[a][b]*x;

                g[a][b] *= 1-x-y;

            }

        cnt++;

    }

}



void dac(int l, int r, int d) {

    if(l > r) return ; 

    if(l == r) {

        for(int i = 0; i < pn; ++i)

            for(int j = 0; i+j < pn; ++j)

                f[p[l]][i+pre][j+1] += g[i][j]*in[l];

        return ;

    }

    int mid = (l+r)>>1;

    backup(d), ins(l, mid), dac(mid+1, r, d+1), recover(d), cnt -= mid-l+1;

    backup(d), ins(mid+1, r), dac(l, mid, d+1), recover(d), cnt -= r-mid;

}



inline void work(int l, int r) {

    pn = pre = 0;

    for(int i = 1; i <= n; ++i)

        if(L[i] <= l && r <= R[i]) {

            p[++pn] = i;

            double len = R[i]-L[i];

            in[pn] = (r-l)/len;

            small[pn] = (l-L[i])/len;

        }

        else if(R[i] <= l) pre++;

    g[0][0] = 1;

    dac(1, pn, 0);

}



int main() {

    read(n);

    for(int i = 1; i <= n; ++i) {

        read(L[i]), read(R[i]);

        v[++vn] = L[i], v[++vn] = R[i];

    }

    std::sort(v+1, v+vn+1);

    vn = std::unique(v+1, v+vn+1)-v-1;

    for(int i = 1; i < vn; ++i)

        work(v[i], v[i+1]);

    for(int i = 1; i <= n; ++i)

        for(int j = 0; j < n; ++j)

            for(int k = 1; k+j <= n; ++k) {

                double t = f[i][j][k]/k;

                for(int r = 1; r <= k; ++r)

                    ans[i][j+r] += t;

            }

    for(int i = 1; i <= n; ++i)

        for(int j = 1; j <= n; ++j)

            printf("%.10lf%c", ans[i][j], j==n?'\n':' ');

    return 0;

}