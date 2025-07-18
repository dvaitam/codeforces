#include <algorithm>
#include <cstdio>
#include <cstdlib>
#include <list>
#include <map>
#include <queue>
#include <set>
#include <stack>
#include <vector>
#include <cmath>
#include <cstring>
#include <string>
#include <iostream>
#include <complex>
#include <sstream>
using namespace std;
 
typedef long long LL;
typedef unsigned long long ULL;
typedef long double LD;
typedef vector<int> VI;
typedef pair<int,int> PII;
 
#define REP(i,n) for(int i=0;i<(n);++i)
#define SIZE(c) ((int)((c).size()))
#define FOR(i,a,b) for (int i=(a); i<(b); ++i)
#define FOREACH(i,x) for (__typeof((x).begin()) i=(x).begin(); i!=(x).end(); ++i)
#define FORD(i,a,b) for (int i=(a)-1; i>=(b); --i)
#define ALL(v) (v).begin(), (v).end()
 
#define pb push_back
#define mp make_pair
#define st first
#define nd second

int hp[2],dt[2],l[2],r[2],p[2],z[2];

double state[2][2][205];
int w = 0;

void go(int f) {
    int wp = 1 - w;
    
    REP(j,hp[f]+1)
        state[wp][f][j] = p[1-f] * 0.01 * state[w][f][j];
            
    double s = 0;
    FORD(j,hp[f]+1,0) {
        if (j + l[1-f] <= hp[f])
            s += state[w][f][j+l[1-f]];
        if (j + r[1-f] + 1 <= hp[f])
            s -= state[w][f][j+r[1-f]+1];
        state[wp][f][j] += 0.01 * (100 - p[1-f]) / z[1-f] * s;
    }
    REP(j,hp[f]+1) {
        if (j >= r[1-f]) continue;
        int a = j - r[1-f];
        int b = min(j - l[1-f], -1);
        double pq = (b - a + 1) * 1.0 / z[1-f] * (100 - p[1-f]) * 0.01;
        state[wp][f][0] += pq * state[w][f][j];
    }

    REP(j,hp[1-f]+1)
        state[wp][1-f][j] = state[w][1-f][j];

    w = wp; 
}

int main() {
    REP(i,2) {
        scanf("%d%d%d%d%d",&hp[i],&dt[i],&l[i],&r[i],&p[i]);
        z[i] = r[i] - l[i] + 1;
    }
    
    if (p[0] == 100 || p[1] == 100) {
        printf("%0.6lf\n",(p[0] == 100)?0.0:1.0);
        return 0;
    }
    
    state[w][0][hp[0]] = 1.0;
    state[w][1][hp[1]] = 1.0;
    double result = 0.0;
    for(int t = 0;;++t) {
        if (t % dt[0] == 0) {
            result -= state[w][1][0] * (1.0 - state[w][0][0]);       
            go(1);
            result += state[w][1][0] * (1.0 - state[w][0][0]);       

            double q = state[w][0][0] + state[w][1][0] - state[w][0][0] * state[w][1][0]; 
            if (q + 1e-7 > 1) goto result;
        }
        if (t % dt[1] == 0) go(0);        
    }
    result:printf("%0.6lf\n",result);
}