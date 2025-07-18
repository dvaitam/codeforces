#include<iostream>
#include<cstdio>
using namespace std;
int main()
{
    int n,p1,p2,p3,t1,t2;
    while(cin>>n>>p1>>p2>>p3>>t1>>t2)
    {
                                     int l[n];
                                     int r[n];
                                     long long pc=0;
                                     cin>>l[0]>>r[0];
                                     pc=(r[0]-l[0])*p1;
                                     for(int i=1;i<n;i++)
                                     {
                                             cin>>l[i]>>r[i];
                                             pc+=(r[i]-l[i])*p1;
                                             int x=l[i]-r[i-1];
                                             if(x<=t1)
                                             {
                                                      pc+=x*p1;
                                                      continue;
                                             }
                                             pc+=t1*p1;
                                             x-=t1;
                                             if(x<=t2)
                                             {
                                                      pc+=(x)*p2;
                                                      continue;
                                             }
                                             pc+=t2*p2;
                                             x-=t2;
                                             pc+=x*p3;
                                     }
                                     cout<<pc<<endl;
    }
}