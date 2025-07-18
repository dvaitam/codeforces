#include<bits/stdc++.h>

using namespace std;

pair<int,int> a[100+10];

int main(){

	int n,k;

	cin>>n>>k;

	for(int i=0;i<n;i++){

		cin>>a[i].first;

		a[i].second=i+1;

	}

	sort(a,a+n);

	int ans=0;

	for(int i=0;i<n;i++){

		k-=a[i].first;

		if(k>=0)

			ans++;

		else

			break;

	}

	cout<<ans<<endl;

	for(int i=0;i<ans;i++)

		cout<<a[i].second<<' ';

}