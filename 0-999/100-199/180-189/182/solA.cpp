#include <iostream>
#include <fstream>
#include <algorithm>
#include <bitset>
#include <cassert>
#include <cctype>
#include <cmath>
#include <complex>
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <ctime>
#include <deque>
#include <iomanip>
#include <map>
#include <numeric>
#include <queue>
#include <set>
#include <stack>
#include <sstream>
#include <string>
#include <vector>
using namespace std;

#define EPS 1e-9
#define INF MOD
#define MOD 1000000007LL
#define fir first
#define foreach(it,X) for(it=X.begin();it!=X.end();it++)
#define iss istringstream
#define ite iterator
#define ll long long
#define mp make_pair
#define rep(i,n) rep2(i,0,n)
#define rep2(i,m,n) for(int i=m;i<n;i++)
#define pi pair<int,int>
#define pb push_back
#define sec second
#define sh(i) (1LL<<i)
#define sst stringstream
#define sz size()
#define vi vector<int>
#define vc vector
#define vl vector<ll>
#define vs vector<string>

#define y1 Y1

int a,b,ax,ay,bx,by,n;
int x1[1010],y1[1010],x2[1010],y2[1010],d[1010];

int dist2(int A,int B){
    int dx=0,dy=0;
    if(x2[A] < x1[B])dx=x1[B]-x2[A];
    else if(x2[B] < x1[A])dx=x1[A]-x2[B];
    if(y2[A] < y1[B])dy=y1[B]-y2[A];
    else if(y2[B] < y1[A])dy=y1[A]-y2[B];
    return dx*dx+dy*dy;
}

double dist(int A,int B){
    double dx=0,dy=0;
    if(x2[A] < x1[B])dx=x1[B]-x2[A];
    else if(x2[B] < x1[A])dx=x1[A]-x2[B];
    if(y2[A] < y1[B])dy=y1[B]-y2[A];
    else if(y2[B] < y1[A])dy=y1[A]-y2[B];
    return sqrt(dx*dx+dy*dy);
}

int main(){
    cin>>a>>b>>ax>>ay>>bx>>by>>n;
    rep(i,n){
        cin>>x1[i]>>y1[i]>>x2[i]>>y2[i];
        if(x1[i]==x2[i]){
            if(y1[i]>y2[i])swap(y1[i],y2[i]);
        }else{
            if(x1[i]>x2[i])swap(x1[i],x2[i]);
        }
    }
    x1[n]=x2[n]=ax;
    y1[n]=y2[n]=ay;
    x1[n+1]=x2[n+1]=bx;
    y1[n+1]=y2[n+1]=by;
    queue<int> Q;
    Q.push(n);
    fill(d,d+n+2,INF);
    d[n]=0;
    while(Q.sz){
        int cur=Q.front();Q.pop();
        if(cur==n+1)break;
        rep(i,n+2){
            if(d[i]>d[cur]+1 && dist2(cur,i)<=a*a){

                d[i]=d[cur]+1;
                Q.push(i);
            }
        }
    }
    if(d[n+1]==INF){cout<<-1;return 0;}
    double ans=INF;
    rep(i,n+1)if(d[i]==d[n+1]-1){
        ans=min(ans,(a+b)*d[i]+dist(i,n+1));
    }
    cout<<ans;
}