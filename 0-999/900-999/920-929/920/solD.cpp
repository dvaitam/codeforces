#include <bits/stdc++.h>
using namespace std;

const double eps = 1e-9;

typedef long long ll;
typedef unsigned long long ull;
typedef vector<int> vi;
typedef vector<string> vs;
typedef pair<int, int> pii;

#define sz(c) int((c).size())
#define all(c) (c).begin(), (c).end()
#define FOR(i,a,b) for (int i = (a); i < (b); i++)
#define FORD(i,a,b) for (int i = int(b)-1; i >= (a); i--)
#define FORIT(i,c) for (__typeof__((c).begin()) i = (c).begin(); i != (c).end(); i++)
#define mp make_pair
#define pb push_back
#define mscanf(...) if(0 == scanf(__VA_ARGS__)){fprintf(stderr, "Could not parse arguments\n");}

const int MAXN = 5050;
int N,K,V;
int A[MAXN];
int L[MAXN];
int main(){
    mscanf("%d%d%d",&N,&K,&V);
    FOR(i,0,N)mscanf("%d",A+i);
    queue<tuple<int,int,int> > operations;
    if(V%K==0){
        operations.push(make_tuple(A[0],0,1));
        A[1] += A[0];
        A[0] = 0;
    }
    FOR(k,0,K)L[k] = -1;
    L[A[0]%K]=0;
    int aim = V%K;
    FOR(n,0,N)if(L[aim] == -1){
        int v = A[n]%K;
        if(L[v] == -1)L[v] = n;
        FOR(k,0,K)if(L[k] != -1 && L[k] != n){
            int kk = k + v;
            if(kk >= K)kk -= K;
            if(L[kk] == -1)L[kk] = n;
        }
    }
    if(L[aim] == -1){
        printf("NO\n");
        return 0;
    }
    int target = L[aim];
    V -= A[target];
    A[target] = 0;
    while(V%K != 0){
        int v =  V%K;
        if(v < 0)v += K;
        assert(A[L[v]] != 0);
        operations.push(make_tuple(A[L[v]],L[v],target));
        V -= A[L[v]];
        A[L[v]] = 0;
    }
    int other = 0;
    if(target == 0)other = 1;
    if(V < 0){
        operations.push(make_tuple((-V)/K,target,other));
    } else if(V > 0) {
        FOR(i,0,N)if(A[i] > 0 && i != other && i != target){
            operations.push(make_tuple(A[i],i,other));
            A[other] += A[i];
            A[i] = 0;
        }
        if(A[other] < V){
            printf("NO\n");
            return 0;
        }
        operations.push(make_tuple(V/K,other,target));
    }
    printf("YES\n");
    while(!operations.empty()){
        auto op = operations.front();
        operations.pop();
        printf("%d %d %d\n", max(1,get<0>(op)), get<1>(op)+1, get<2>(op)+1);
    }
 return 0;
}