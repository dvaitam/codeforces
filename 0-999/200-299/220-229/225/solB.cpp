#include <iostream>
#include <vector>
using namespace std;
int main()
{long long s,k; cin>>s>>k;
vector <long long> v;
 int f[1000]={0};
 f[0]=0;f[1]=1;
int i=1;
do
 {i++;
  for (int j=1;j<=k;j++)
  {f[i]+=f[i-j];
    if (i-j<=0) break;
  }
 }
 while (f[i]<s);
 while (s)
 {if (s>=f[i]) {s-=f[i];v.push_back(f[i]);}
   i--;
 }
v.push_back(0);
 cout<<v.size()<<endl;
   for(int i=0; i<v.size(); i++)
   cout<<v[i]<<" ";
return 0;
}