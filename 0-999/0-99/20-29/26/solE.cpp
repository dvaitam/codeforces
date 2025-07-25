#include<cstdio>
#include<cstdlib>
#include<cstring>
#include<cmath>
#include<iostream>
#include<fstream>
#include<algorithm>
#include<vector>
#include<map>
#include<set>
#include<string>
using namespace std;
int a[103],t[103];
bool cmp (int u,int v){
	return a[u]<a[v];
}
int main(){
	std::ios::sync_with_stdio(false);
	int i,j,k,m,n;
	cin >> n >> m;
	for (i=1;i<=n;t[i]=i,++i)
		cin >> a[i];
	if (n==1){
		if (m==a[1]){
			cout << "Yes" << endl;
			for (i=1;i<=a[1];++i)
				cout << "1 1 ";
			cout << endl;
		}
		else cout << "No" << endl;
	}
	else{
		sort(t+1,t+n+1,cmp);
		if (m<1 || m==1 && a[t[1]]>1) cout << "No" << endl;
		else{
			for (k=0,i=1;i<=n;++i) k+=a[i];
			if (k<m) cout << "No" << endl;
			else{
				cout << "Yes" << endl;
				if (a[t[1]]<=m){
					cout << t[1];
					for (i=2;i<=n && k>m;++i)
						for (;a[t[i]]>0 && k>m;--k,--a[t[i]])
							cout << ' ' << t[i] << ' ' << t[i];
					cout << ' ' << t[1];
					for (--a[t[1]],i=1;i<=n;++i)
						for (j=1;j<=a[t[i]];++j)
							cout << ' ' << t[i] << ' ' << t[i];
					cout << endl;
				}
				else{
					cout << t[2];
					for (i=m--;i<=a[t[1]];++i)
						cout << ' ' << t[1] << ' ' << t[1];
					cout << ' ' << t[2] << ' ' << t[1];
					for (--a[t[2]],i=2;i<=n;++i)
						for (j=1;j<=a[t[i]];++j)
							cout << ' ' << t[i] << ' ' << t[i];
					cout << ' ' << t[1];
					for (i=1;i<m;++i)
						cout << ' ' << t[1] << ' ' << t[1];
					cout << endl;
				}
			}
		}
	}
	return 0;
}