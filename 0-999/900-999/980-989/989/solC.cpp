#include<iostream>
using namespace std;
int main(){
	int a[5]; puts("44 49"); 
	cin>>a[4]>>a[1]>>a[2]>>a[3];
	for(int i=1;i<=4;++i){
		char su='A'+i-1, pi=i==4?'A':su+1;
		for(int j=1;j<=5;++j,cout<<su<<endl){
			for(int o=1;o<=49;++o)cout<<su; puts("");
			for(int o=1;o<=24;++o)cout<<su<<(--a[i]>0?pi:su);
		}   for(int o=1;o<=49;++o)cout<<su; puts("");
	}return 0;
}
