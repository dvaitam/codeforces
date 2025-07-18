#pragma warning(disable:4996)
#define loop(i,start,stop) for(int i = start; i < stop; i++)
#include <stdio.h>

long long r, r0, r1, s, s0, s1, t, t0, t1;
long long Q;

void Euclid(long long r0, long long r1, long long& x, long long& y) {
    s0 = t1 = 1;
    s1 = t0 = 0;
    while (r1 != 0) {
        Q = r0 / r1;
        r = r0 - Q*r1;
        s = s0 - Q*s1;
        t = t0 - Q*t1;
        r0 = r1;
        r1 = r;
        s0 = s1;
        s1 = s;
        t0 = t1;
        t1 = t;
    }
    x = s0;
    y = t0;
}

long long a, b, c, d;
//long long x, y, s, t, x0, y0;
long long x, y, x0, y0;

long long abs(long long x) {
    return x < 0 ? -x : x;
}


long long gcd(long long a, long long b) {
    return (a%b == 0) ? abs(b) : gcd(b, a%b);
}

bool isPossible(long long a, long long b, long long c) {
    return (c%gcd(a, b) == 0);
}


bool swapped;

void swap_(long long& a, long long& b) {
    long long temp = a;
    a = b;
    b = temp;
    swapped = true;
}

void solve() {
    if (isPossible(a, b, c)) {


        if (a < b) {
            swap_(a, b);
        }

        Euclid(a, b, s, t);


        d = gcd(a, b);
        x0 = (c / d) * s;
        y0 = (c / d) * t;

        x = x0;
        y = y0;

        
        long long k = 1;
        if (x < 0) {
            while (x < 0) {
                x = x0 + k * (b / d);
                y = y0 - k * (a / d);
                k++;
            }
        }
        else if (y < 0) {
            while (y < 0) {
                x = x0 - k * (b / d);
                y = y0 + k * (a / d);
                k++;
            }
        }

        if (x < 0 || y < 0) {
            printf("NO\n");
            return;
        }
        puts("YES");
        if (!swapped) {
            printf("%I64d %I64d\n", x, y);
        }
        else {
            printf("%I64d %I64d\n", y, x);
        }
    }
    else {
        puts("NO");
    }
}



void init() {
    swapped = false;
    scanf("%I64d %I64d %I64d", &c, &a, &b);
}
int main() {
       //freopen("451a.in", "r", stdin);
        init();
        solve();
}