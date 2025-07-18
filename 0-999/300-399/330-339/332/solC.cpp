#include <map>
#include <set>
#include <queue>
#include <ctime>
#include <cmath>
#include <string>
#include <cstdio>
#include <vector>
#include <cstdlib>
#include <cstring>
#include <iostream>
#include <algorithm>
using namespace std;

#define all(a) a.begin(),a.end()
#define clr(a) memset(a,0,sizeof(a))
#define fill(a,b) memset(a,b,sizeof(a))
#define pb push_back
#define mp make_pair

typedef long long LL;
typedef vector<int> VI;
typedef pair<int,int> PII;
typedef vector<pair<int,int> > VII;
typedef VI::iterator IT;
const double eps = 1e-8;
const double pi = acos(-1.0);

const int N = 100000 + 10;

struct E{
    int a, b, nid;
    bool operator < (const E&o)const{
        if(b==o.b) return a>o.a;
        return b<o.b;
    }
} e[N+N];

struct Node{
    int a, nid;
    Node(int _a=0, int _nid=0):a(_a),nid(_nid){}
    bool operator < (const Node&o)const{
        return a > o.a;
    }
};

priority_queue<Node> q;
int main(){
    int n ,p, k, i, l, r;
    scanf("%d%d%d",&n,&p,&k);
    for(i=1;i<=n;++i) scanf("%d%d",&e[i].a,&e[i].b), e[i].nid = i;
    sort(e+1,e+n+1);
    for(i=n;i>=n-k+1;--i) q.push(Node(e[i].a,i));
    for(i=n-k;i>p-k;--i){
        if(!q.empty() && e[i].a > q.top().a){
            q.pop();
            q.push(Node(e[i].a, i));
        }
    }
    int last = n+1;
    while(!q.empty()){
        printf("%d ", e[q.top().nid].nid);
        last = min(last, q.top().nid);
        q.pop();
    }
    for(i=last-1;i>=last-1-p+k+1;--i) printf("%d ",e[i].nid);
    puts("");
    return 0;
}