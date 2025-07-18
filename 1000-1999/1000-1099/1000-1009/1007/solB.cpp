#include <bits/stdc++.h>
#define MAXN 100000

int d[1 + MAXN];

inline int gcd(int a, int b){
    int r;
    while(b > 0){
        r = a % b;
        a = b;
        b = r;
    }
    return a;
}

struct Elem{
    int x, y, z;
};
std::map <std::tuple <int, int, int>, int> M;

int main(){
    FILE*fi,*fo;
    //fi = fopen("a.in","r");
    //fo = fopen("a.out","w");
    fi = stdin;
    fo = stdout;

    for(int i = 1; i <= MAXN; i++)
        for(int j = i; j <= MAXN; j += i)
            d[j]++;

    int t;
    fscanf(fi,"%d", &t);
    //t = 10000;
    //srand(time(NULL));
    for(int z = 1; z <= t; z++){
        int a, b, c;
        fscanf(fi,"%d%d%d", &a, &b, &c);
        //a = rand() % 15+ 1;
        //b = rand() % 15 + 1;
        //c = rand() % 15 + 1;
        long long ans = 0, tmp;
        int x = gcd(gcd(a, b), c);
        int m = gcd(a, b);
        int n = gcd(b, c);
        int p = gcd(c, a);


        ///toate combinatiile
        ans += 1LL * d[a] * d[b] * d[c];
        ///1 numar care divide doar a, b sau c
        tmp = d[a] - d[m] - d[p] + d[x];
        ans -= 1LL * tmp * d[n] * (d[n] - 1) / 2;
        tmp = d[b] - d[m] - d[n] + d[x];
        ans -= 1LL * tmp * d[p] * (d[p] - 1) / 2;
        tmp = d[c] - d[n] - d[p] + d[x];
        ans -= 1LL * tmp * d[m] * (d[m] - 1) / 2;
        ///toate numerele divid pe toata lumea
        //sunt diferite
        ans -= 5LL * d[x] * (d[x] - 1) * (d[x] - 2) / 6;
        //doua egale
        ans -= 2LL * d[x] * (d[x] - 1);
        //trei egale
        ans -= 0;
        ///un numar divide 2, 2 divid pe toata lumea
        //trei diferite
        tmp = d[m] - d[x];
        ans -= 3LL * tmp * d[x] * (d[x] - 1) / 2;
        tmp = d[n] - d[x];
        ans -= 3LL * tmp * d[x] * (d[x] - 1) / 2;
        tmp = d[p] - d[x];
        ans -= 3LL * tmp * d[x] * (d[x] - 1) / 2;
        //ultimele doua egale
        tmp = d[m] - d[x];
        ans -= 1LL * tmp * d[x];
        tmp = d[n] - d[x];
        ans -= 1LL * tmp * d[x];
        tmp = d[p] - d[x];
        ans -= 1LL * tmp * d[x];
        ///un numar le divide pe toate, 2 divid 2
        //ultimele divid aceleasi numere
            //sunt diferite
            tmp = d[m] - d[x];
            ans -= 1LL * d[x] * tmp * (tmp - 1) / 2;
            tmp = d[n] - d[x];
            ans -= 1LL * d[x] * tmp * (tmp - 1) / 2;
            tmp = d[p] - d[x];
            ans -= 1LL * d[x] * tmp * (tmp - 1) / 2;
            //sunt egale
            ans -= 0;
        //divid doua perechi diferite
            ans -= 2LL * d[x] * (d[m] - d[x]) * (d[n] - d[x]);
            ans -= 2LL * d[x] * (d[m] - d[x]) * (d[p] - d[x]);
            ans -= 2LL * d[x] * (d[n] - d[x]) * (d[p] - d[x]);
        ///toate 3 divid doar 2
        tmp = d[n] - d[x] + d[p] - d[x];
        ans -= 1LL * tmp * (d[m] - d[x]) * (d[m] - d[x]- 1) / 2;
        tmp = d[m] - d[x] + d[n] - d[x];
        ans -= 1LL * tmp * (d[p] - d[x]) * (d[p] - d[x]- 1) / 2;
        tmp = d[m] - d[x] + d[p] - d[x];
        ans -= 1LL * tmp * (d[n] - d[x]) * (d[n] - d[x]- 1) / 2;

        ans -= 1LL * (d[m] - d[x]) * (d[n] - d[x]) * (d[p] - d[x]);

        fprintf(fo,"%lld\n", ans);
    }

    return 0;
}