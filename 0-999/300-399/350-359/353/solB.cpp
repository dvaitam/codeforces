#include<iostream>
using namespace std;
int main(){
		int n,a[200];
		int g[90][2];
		cin>>n;
		for(int i=0;i<90;i++){
			g[i][1]=0;
			g[i][0]=0;
			}
		for(int i=0;i<2*n;i++){
			cin>>a[i];
			g[a[i]-10][0]++;
			}
		int t=0;
		int g1=0,g2=0;
		for(int i=0;i<90;i++){
			if(g[i][0]==1){
				if(t==0){
					g[i][1]=1;
					g1++;
				}
				if(t==1)
					g2++;
				g[i][0]=0;
				t=(t+1)%2;
				}
			}
		for(int i=0;i<90;i++){
			if(g[i][0]>1){
				g[i][1]=g[i][0]/2;
				g1++;
				g2++;
				if(g[i][0]%2==1){
					if(t==0)
						g[i][1]++;
					t=(t+1)%2;
					}
				}
			}
			cout<<g1*g2<<endl;
			for(int i=0;i<2*n;i++){
				if(g[a[i]-10][1]>0){
					g[a[i]-10][1]--;
					cout<<1<<' ';
					}
				else
					cout<<2<<' ';
				}
			return 0;
			}