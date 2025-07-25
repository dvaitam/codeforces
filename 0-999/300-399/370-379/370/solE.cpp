#include<iostream>
using namespace std;
int a[200010],cnt[100002],n;
pair<int,int> Max[200010],Min[200010];

int main()
{
	ios::sync_with_stdio(0); 
	cin>>n;
	for(int i=1;i<=n;++i)
		cin>>a[i];
	if(a[1]>1)  cout<<"-1"<<endl;
	else{
		bool flag=1;
		Max[1]=Min[1]=make_pair(1,1);
		for(int i=2;i<=n&&flag;++i)
		{
			if(Max[i-1].second>=2)  Max[i]=make_pair(Max[i-1].first+1,1);
			else Max[i]=make_pair(Max[i-1].first,Max[i-1].second+1);
      		                            	          // 上一本书数目大于2。能过，不断地加 112233  小于二，则添加
			if(Min[i-1].second==5)  Min[i]=make_pair(Min[i-1].first+1,1);
			else Min[i]=make_pair(Min[i-1].first,Min[i-1].second+1); 
			                                          //  i[1]=5  满书，放下一个   111112222233333
			if(a[i]){
				if(a[i]<Min[i].first || a[i]>Max[i].first) flag=0;                //cout<<"\1"<<endl;}
				else{ 
					Max[i]=min(Max[i],make_pair(a[i],5));
					Min[i]=max(Min[i],make_pair(a[i],1));
				}  // 状态转移 
			}
		}
		if(flag)
		{
			pair<int,int> ans = Max[n];
			if(ans.second==1)  ans=make_pair(ans.first-1,5);  // 最后一天要大于等于2 
			if(ans<Min[n])  cout<<-1<<endl;
			else{
				cout<<ans.first<<endl;
				a[n]=ans.first;
				cnt[a[n]]=1;
				for(int i=n-1;i;--i)
				{
					a[i]=min(Max[i].first,a[i+1]);
					if(cnt[a[i]]==5)  a[i]--;
					cnt[a[i]]++;
				}
				for(int i=1;i<=n;++i)
					cout<<a[i]<<' ';
				cout<<endl;
			} 
		}
		else cout<<-1<<endl;
	}	
	return 0;
}