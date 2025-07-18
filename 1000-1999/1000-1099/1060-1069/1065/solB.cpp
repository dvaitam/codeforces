#include <iostream>
#include <vector>
#include <algorithm>
#include <set>
#include <bits/stdc++.h>
#include <stdlib.h>
 
#define tests int t;cin>>t;while(t--)
#define array(A,n) long long int A[n];for(int i=0;i<n;i++)cin>>A[i];
#define array2(A,B,n,m) long long int A[n] , B[m];for(int i=0;i<n;i++)cin>>A[i]; for(int i=0;i<m;i++)cin>>B[i];
#define print(A,n) for(int i=0;i<n;i++)cout<<A[i]<<" ";cout<<endl;
#define matrix(A,n,m) long long int A[n][m];for(int i=0;i<n;i++){for(int j=0;j<m;j++)cin>>A[i][j];}
#define vit(a) vector<long long int>::iterator a
#define vrit(a) vector<long long int>::reverse_iterator a
#define sit(a) set<long long int>::iterator a
#define vec(a) vector<long long int > a
#define set(a) set<long long int > a 

#define yes "yes"
#define no "no"
#define YES "YES"
#define NO "NO"
#define Yes "Yes"
#define No "No"
long long int  MIN(long long int a,long long int b)
{if(a>b)return b; else return a;}
 
long long int  MAX(long long int a,long long int b)
{if(a>b)return a; else return b;}
long long int  MOD(long long int x)
{if(x>=0)return x;else return-x;}

 
 
using namespace std;
int main() 
{
    long long int n , m;
    cin>>n>>m;
    
    long long int mm=m*2;
    long long int min = n-mm;
    if(min>=0)
    cout<<min<<" ";
    else cout<<0<<" ";
    {
       long long int dd = m;
        long long int ans= 0 ;
        for(long long int i=0;i<=n;i++)
        {
            long long int tt = (i*(i-1))/2;
            if(tt>=dd)
            {
                ans = i;
                break;
            }
        }
         //if(ans!=-1)
        long long int  max = n-ans;
        cout<<max<<endl;
        
        
    } 
  
      return 0 ;
}