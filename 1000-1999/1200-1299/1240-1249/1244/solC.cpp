#include <string.h>
#include <unordered_map>

#include <sstream>
#include <stdio.h>
#include <assert.h>
#include <math.h>
#include <bitset>
#include <algorithm>
#include <iostream>
#include <stack>
#include <queue>
#include <set>

#include <map>
#include <vector>
#include <string>
#include <stdlib.h>

#define ll long long
#define clr(x) memset(x,0,sizeof(x))
#define _clr(x) memset(x,-1,sizeof(x))
#define fr(i,a,b) for(int i = a; i < b; ++i)
#define frr(i,a,b) for(int i = a; i > b; --i)
#define pb push_back
#define sf scanf

#define pf printf
#define mp make_pair
#define N 1000000

using namespace std;

ll gcd(ll a, ll b) {
    return b?gcd(b,a%b):a;
}

ll exgcd(ll a, ll b, ll &x, ll &y) {
    if(!b) {
        x=1;
        y = 0;
        return a;
    }
    ll d = exgcd(b,a%b,x,y);
    ll t = x;
    x = y;
    y = t-(a/b)*y;
    return d;
}

void sol(){
    ll n, p, w,d;
    cin>>n>>p>>w>>d;
    ll t = gcd(w,d);
    if(p%t) {
        printf("-1\n");
        return;
    }
    p/=t;
    w/=t;
    d/=t;
    ll x,y;
    exgcd(w,d,x,y);

   x*=p;
   y*=p;
   ll min_x = x%d;
   if(min_x <0) min_x += d;

   ll min_y = y%w;
   if(min_y<0) min_y += w;

   for(ll i = 0; i <=y; ++i) {
       if((p-i*d)%w==0) {
           min_y= i;
           break;
       }
   }

   for(ll i = 0; i <=d; ++i) {
       if((p-i*w)%d==0) {
           ll x = i;
           y = (p-i*w)/d;
           if(x>=0&&y>=0&&x+y<=n) {
                cout<<x<<" "<<y<<" "<<n-(x+y)<<endl;
                return;
           }
       }
   }

   for(ll i = 0; i <=w; ++i) {
       if((p-i*d)%w==0) {
           ll y = i;
           ll x = (p-i*d)/w;
           if(x>=0&&y>=0&&x+y<=n) {
                cout<<x<<" "<<y<<" "<<n-(x+y)<<endl;
                return;
           }
       }
   }



   /*
   x = min_x;
   y = (p-x*w)/d;
    if(x>=0&&y>=0&&x+y<=n) {
        cout<<x<<" "<<y<<" "<<n-(x+y)<<endl;
        return;
    }

    y = min_y;
    x = (p-y*d)/w;

    if(x>=0&&y>=0&&x+y<=n) {
        cout<<x<<" "<<y<<" "<<n-(x+y)<<endl;
        return;
    }
    */
    printf("-1\n");
    return;


   //printf("w = %lld d = %lld p = %lld\n",w,d,p);
   //printf("x = %lld y = %lld\n",x,y);

   /*
    if(x>=0&&y>=0&&x+y<=n) {
        cout<<x<<" "<<y<<" "<<n-(x+y)<<endl;
    }
    else {
        if(y<=0) {
            printf("-1\n");
            return;
        }

        ll min_t = (ll)(-y/w);
        ll max_t = (ll)(x/d);

        if(min_t>max_t) {
            printf("-1\n");
            return;
        }
        ll ft = (n-(x+y))/(w-d);
        max_t = min(max_t,ft);
        if(min_t>max_t) {
            printf("-1\n");
            return;
        }
        ll t = min_t;
        //printf("t = %lld\n",t);
        x=x-d*t;
        y=y+w*t;
        if(x>=0&&y>=0&&x+y<=n)
            cout<<x<<" "<<y<<" "<<n-(x+y)<<endl;
        else
            cout<<-1<<endl;
    }
    */
    //printf("x*w + y*d = p %lld\n",x*w+y*d);
}

int main() {
    sol();
}