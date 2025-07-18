/* 

   KUMAR GAURAV       

*/



#include <bits/stdc++.h>

using namespace std;

#define loop(i,start,end) for (lli i=start;i<=end;i++)

#define pool(i,start,end) for(lli i=start;i>=end;i--)

#define kg() lli t;cin>>t;while(t--)

#define vi(v) vector <long long  int> v;

#define pb(n) push_back(n)

#define mp(a,b) make_pair(a,b)

#define fill(a,value) memset(a,value,sizeof(a))

#define MOD 1000000007

#define PI  3.14159265358979323846

#define MAX 1000002

#define vpi(v) vector <pair <long long int, long long int> > v

#define lli long long int

#define debug() cout<<"######"<<endl



double dp[101][101][101]={0};



double func(int r,int s, int p)

{

   if (dp[r][s][p] != 0)

        return dp[r][s][p];

    if(r==0)

        return 0;

    else if(r!=0 && p==0)

        return 1;

    else if(r!=0 && s==0)

        return 0;

    double d= (double)(r*p)/(double)(r*s+s*p+p*r)*(double)func(r-1,s,p) + (double)r*s/(double)(r*s+s*p+p*r)*(double)func(r,s-1,p) + (double)s*p/(double)(r*s+s*p+p*r)* (double)func(r,s,p-1);



   dp[r][s][p]=d;

    return d;



}





int main()

{

    int  r,s,p;

    cin>>r>>s>>p;

//r//s//p

printf("%.9f %.9f %.9f",func(r,s,p),func(s,p,r),func(p,r,s));

  // cout<< func(r,s,p)<<" "<<func(s,p,r)<<" "<<func(p,r,s);







 return 0;

}