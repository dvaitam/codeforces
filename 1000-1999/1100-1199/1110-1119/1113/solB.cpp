#include<bits/stdc++.h>
using namespace std;
#include<stdio.h>
#include<math.h>
#define pb push_back
#define mp make_pair
#define endl '\n'
#define lli long long int
#define li long int 
#define FOR(i,n) for (li i = 0; i < n; i++)
#define loop(i,a,b) for (li i = a; i < b; i++)
#define FORD(i,n) for (li i = n-1; i >= 0; i--)
#define loopD(i,a,b) for (li i = a; i >= b; i--)
#define vi vector<int>
#define vl vector<li>
#define vll vector<lli>
#define CLEAR(a)  memset(a,0,sizeof(a))
#define CLEARN(a) memset(a,-1,sizeof(a))
#define F first
#define S second
#define gcd(a,b) __gcd((a),(b))
#define lcm(a,b) ((a)*(b))/gcd((a),(b))
#define mii map<int,int>
#define mll map<li,li>
#define pii pair<int,int>
#define pll pair<li,li>
#define all(v) v.begin(),v.end()
#define pq priority_queue<li>
#define spq priority_queue <li, vector<li>, greater<li>>
#define pql priority_queue<pair<li,li>>
#define pqs priority_queue<pair<li,li>,vector< pair<li,li> >,greater<pair<li,li>>>
#define lb(v,x) lower_bound(v.begin(),v.end(),x)-v.begin()
long int mod=1e9+7;
#define PI 3.1415926535897
li i,j,n,m,k,l,b[105];
vector<li> adj[105];
int main()
{
    #ifndef ONLINE_JUDGE
    freopen("input.txt","r",stdin);
    freopen("output.txt","w",stdout);
    #endif  
    ios::sync_with_stdio(0);
    cin.tie(0); 
     for(j=1 ; j<=100 ; j++)
    {
   for (i=2; i<=sqrt(j); i++) 
    { 
        if (j%i==0) 
        { 
            if (j/i == i) 
           adj[j].pb(i);
            else
            { 
               adj[j].pb(i);
                adj[j].pb(j/i); 
            } 
        } 
    }
    }
    li ans=0;
    cin>>n;
    loop(i,0,n)
    {
        cin>>j;
       b[j]++;
       ans+=j;
    }
    li bha=ans;
    int flag=0;
    for(j=100 ; j>0 ; j--){
        if(b[j]){
        if(adj[j].size()){
          //  cout<<"j is "<<j<<endl;
          for(auto it=adj[j].begin();it!=adj[j].end();it++){
              flag=0;
              loop(i,1,101)
              if(b[i]){
             li prevalue=(i+j); 
             li curvalue=(i*(*it))+(j/(*it));
             if(curvalue<prevalue)
             {
              //   cout<<"i & j "<<i<<" "<<j<<endl;
                bha=min(bha,ans-prevalue+curvalue);
              //flag=1;   
             }
              }
          }
      }
     }
    }
    li sum=0;
  /*  loop(i,1,101)
    sum+=(i*b[i]); */
    cout<<bha;
    return 0;
}