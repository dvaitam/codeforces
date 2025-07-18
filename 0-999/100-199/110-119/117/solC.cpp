//BE NAMASH


#include<iostream>
#include<stdio.h>
using namespace std;
const int MN=5100;
bool a[MN][MN];
short int A,B,C,l[MN],t[MN];
int n;
int find(int p,int q){
	if(q-p<=2){
		return 0;
	}
	int ef1=p+1,ef2=q-1;
	int v=l[p];
	t[p]=l[p];
	for(int i=p+1;i<q;i++){
		if(a[v][l[i]]){
			t[ef2--]=l[i];
		}else{
			t[ef1++]=l[i];
		}
	}
	for(int i=p+1;i<ef1;i++){
		int v2=t[i];
		for(int j=ef2+1;j<q;j++){
			if(a[t[j]][v2]){
				A=v;
				B=t[j];
				C=t[i];
				return 1;
			}
		}
	}
	for(int i=p;i<q;i++){
		l[i]=t[i];
	}
	if(find(p+1,ef1)||find(ef2+1,q)){
		return 1;
	}return 0;
}
char inp[50000000];
int main(){
//	scanf("%d",&n);
	cin>>n;
	cin.read(inp,49999999);
	int st=0;
	for(int i=0;i<n;i++){
		for(int j=0;j<n;j++){
			while(inp[st]!='0'&&inp[st]!='1'){
				st++;
			}
			a[i][j]=inp[st]-'0';
			st++;
		}
		l[i]=i;
	}
	if(find(0,n)){
		cout<<A+1<<" "<<B+1<<" "<<C+1<<endl;
	}else{
		cout<<"-1"<<endl;
	}return 0;
}