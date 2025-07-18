#include <bits/stdc++.h>
         
using namespace std;
         
typedef long long ll;
const int M=1e6+7;
const int M2=1e5;
const ll q=7057594037927903;
const ll prime=2137;
#define pb push_back
#define mp make_pair
#define f first 
#define s second
#define ios ios_base::sync_with_stdio(0) 

int n,k,per1[M],per2[M],ind1[M],ind2[M],pom,pom2,nic,nic2,wazne=INT_MAX,zlicz[M]; char napis[M],akt='z',napis2[M];

int main()
{
    ios;
    cin>>n>>k;
    
    for(int i=1; i<=n; i++)
       napis[i]='Z',napis2[i]='Z';
    
    for(int i=1; i<=n; i++)
    {
		cin>>per1[i];
		ind1[per1[i]]=i;
	}
	
	for(int i=1; i<=n; i++)
	{
		cin>>per2[i];
		ind2[per2[i]]=i;
	}
	
	napis[per1[n]]=akt;
	for(int j=ind2[per1[n]]; j<=n; j++)
	{
		napis2[per2[j]]=akt;
		napis[per2[j]]=akt;
		wazne=min(wazne,ind1[per2[j]]);
	}
	
	for(int i=n-1; i>=1; i--)
	{
		if(napis[per1[i]]=='Z')
		{
			if(wazne<i)
			{
				napis[per1[i]]=akt;
				for(int j=ind2[per1[i]]; j<=n; j++)
            	{
					if(napis2[per2[j]]!='Z')
					  break;
		            napis2[per2[j]]=akt;
		            napis[per2[j]]=akt;
		            wazne=min(wazne,ind1[per2[j]]);
	            }
			}
			else
			{
				if(akt != 'a')
				  akt=(char)(akt-1);
				for(int j=ind2[per1[i]]; j<=n; j++)
            	{
					if(napis2[per2[j]]!='Z')
					  break;
					  
		            napis2[per2[j]]=akt;
		            napis[per2[j]]=akt;
		            wazne=min(wazne,ind1[per2[j]]);
	            }
			}
		}
	}
	
	for(int j=1; j<=n; j++)
	  if(!zlicz[(int)napis[j]])
	  {
		  nic++;
		  zlicz[(int)napis[j]]=1;
	  }
	  
	if(nic>=k)
	{
		cout<<"YES\n";
		for(int i=1; i<=n; i++)
		  cout<<napis[i];
	}
    else
    {
		cout<<"NO\n";
	}
    
    //1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27
    
    return 0;	
}