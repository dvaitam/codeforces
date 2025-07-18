/***************************************

    codeforces = topcoder = sahedsohel

    IIT,Jahangirnagar University(42)

****************************************/

#include <bits/stdc++.h>

using namespace std;



#define ll long long int

#define ull unsigned long long int

#define inf (INT_MAX/10)

#define linf (LLONG_MAX/10LL)

#define sc(a) scanf("%d",&a)

#define sc2(a,b) scanf("%d%d",&a,&b)

#define sc3(a,b,c) scanf("%d%d%d",&a,&b,&c)

#define sc4(a,b,c,d) scanf("%d%d%d%d",&a,&b,&c,&d)

#define f(i,n) for(i=0;i<n;i++)

#define fl(c,i,n) for(i=c;i<n;i++)

#define mem(a) memset(a,0,sizeof(a))

#define memn(a) memset(a,-1,sizeof(a))

#define pb push_back

#define aov(a) a.begin(),a.end()

#define mpr make_pair

#define PI (2.0*acos(0.0)) //#define PI acos(-1.0)

#define xx first

#define yy second

#define mxv(a) *max_element(aov(a))

#define mnv(a) *min_element(aov(a))

#define LB(a,x) (lower_bound(aov(a),x)-a.begin())

#define UB(a,x) (upper_bound(aov(a),x)-a.begin())

#define to_c_string(a) a.c_str()

#define strtoint(c) atoi(&c[0])

#define pll pair< ll , ll >

#define pii pair< int , int >

#define pid pair< int , double >

#define pcs(a) printf("Case %d: ", a)

#define nl puts("")

#define endl '\n'

#define dbg(x) cout<<#x<<" : "<<x<<endl



template <class T> inline T bigmod(T p,T e,T M){ll ret = 1;for(; e > 0; e >>= 1){if(e & 1) ret = (ret * p) % M;p = (p * p) % M;}return (T)ret;}

template <class T> inline T gcd(T a,T b){if(b==0)return a;return gcd(b,a%b);}

template <class T> inline T modinverse(T a,T M){return bigmod(a,M-2,M);}   // M is prime}

template <class T> inline T bpow(T p,T e){ll ret = 1;for(; e > 0; e >>= 1){if(e & 1) ret = (ret * p);p = (p * p);}return (T)ret;}



int toInt(string s){int sm;stringstream ss(s);ss>>sm;return sm;}

int toLlint(string s){long long int sm;stringstream ss(s);ss>>sm;return sm;}





///int mnth[]={-1,31,28,31,30,31,30,31,31,30,31,30,31};  //Not Leap Year

///int dx[]={2,1,-1,-2,-2,-1,1,2};int dy[]={1,2,2,1,-1,-2,-2,-1};//Knight Direction

///int dx[]={-1,+1,0,1,0,-1}; // Hexagonal Direction   **

///int dy[]={-1,+1,1,0,-1,0}; //                       *#*

///                                                     **

///const double eps=1e-9;

///int dx[]={0,1,0,-1};int dy[]={1,0,-1,0}; //4 Direction



/*****************************************************************/

/// ////////////////////   GET SET GO    ////////////////////// ///

/*****************************************************************/





#define intx(i,j,k,l) ((a[i]*b[j]-b[i]*a[j])*(a[k]-a[l])-(a[i]-a[j])*(a[k]*b[l]-b[k]*a[l]))/((a[i]-a[j])*(b[k]-b[l])-(b[i]-b[j])*(a[k]-a[l]))

#define inty(i,j,k,l) ((a[i]*b[j]-b[i]*a[j])*(b[k]-b[l])-(b[i]-b[j])*(a[k]*b[l]-b[k]*a[l]))/((a[i]-a[j])*(b[k]-b[l])-(b[i]-b[j])*(a[k]-a[l]))

#define dst(u,v,x,y) sqrt((x*1.0-u*1.0)*(x*1.0-u*1.0)+(y*1.0-v*1.0)*(y*1.0-v*1.0))

#define area(p1,p2,p3) (p1.xx*p2.yy+p2.xx*p3.yy+p3.xx*p1.yy-p1.yy*p2.xx-p2.yy*p3.xx-p3.yy*p1.xx)



int ts,kk=1;



#define M 200005

#define MD 100000007

#define MX 10001LL



double n,r;



double clc(double q)

{

    if(q<1)return 0.0;

    double q1=min(1.0,q);

    double rs=sqrt(r*r+r*r)*q1+(q1*q1*r);

    q-=q1;

    return rs+sqrt(r*r+r*r)*q*2.0+(q*q*r);

}



int main()

{

    int t,i,j,k;



//    sc(t);

//    for(i=0;i<t*t;i++)

//    {

//        cerr<<i<<" "<<i/t+1<<" "<<t+1+i%t<<endl;

//    }



    cin>>j>>k;

    n=j*1.0;

    r=k*1.0;

    double rs=0;

    for(i=1;i<=j;i++)

    {

        rs+=r+r;

        rs+=r*(i-1)+clc(i-1);

        rs+=r*(n-i)+clc(n-i);

    }

    rs/=(n*n);

//    cerr<<rs<<endl;

    cout<<setprecision(12)<<fixed<<rs<<endl;



    return 0;

}