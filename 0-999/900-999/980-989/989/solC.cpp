#include <bits/stdc++.h>
using namespace std;
typedef long long ll;
const int MOD=1e9+7;
char m[55][55];
void init(){
    for(int i=1;i<=50;i++){
        for(int j=1;j<=50;j++){
            m[i][j]='#';
        }
    }
}
void draw(){
    for(int i=1;i<=50;i++){
        for(int j=1;j<=50;j++){
            cout<<m[i][j];
        }
        cout<<endl;
    }
}
void f(int x,int y,char c){
    int w=0;
    for(int i=1;i<=x;i++){
        m[y+w%3][i+1]=c;
        w++;
    }
}

void ff(int k,int x,int y,char c){
    int w=0;
    for(int i=1;i<=x;i++){
        m[y+w%3][i+k]=c;
        w++;
    }
}
int main(){
    int a,b,c,d;
    scanf("%d%d%d%d",&a,&b,&c,&d);
    init();
    for(int i=1;i<=5;i++){
        m[i][1]='A';
    }
    for(int i=1;i<=50;i++){
        m[1][i]='A';
    }
    for(int i=1;i<=40;i++){
        m[i][50]='A';
    }
    for(int i=2;i<=50;i++){
        m[31][i]='A';
    }
    for(int i=31;i<=45;i++){
        m[i][1]='D';
    }
    //************
    if(a>90){
        ff(3,a-90,27,'A');
        a=90;
    }
    if(b>90){
        ff(14,b-90,27,'B');
        b=90;
    }
    if(c>90){
        ff(25,c-90,27,'C');
        c=90;
    }

    //**********
    if(a>45){
        f(45,7,'A');
        a-=45;
    }
    f(a,3,'A');

    if(b>45){
        f(45,11,'B');
        b-=45;
    }
    f(b,15,'B');

    if(c>45){
        f(45,19,'C');
        c-=45;
    }
    f(c,23,'C');
    //27
    if(d>90){
        ff(3,d-90,32,'D');
        d=90;
    }
    if(d>45){
        ff(3,45,36,'D');
        d-=45;
    }
    f(d,45,'D');
    for(int i=1;i<=31;i++){
        for(int j=1;j<=50;j++){
            if(m[i][j]=='#'){
                m[i][j]='D';
            }
        }
    }
    for(int i=31;i<=50;i++){
        for(int j=1;j<=50;j++){
            if(m[i][j]=='#'){
                m[i][j]='A';
            }
        }
    }
    printf("50 50\n");
    draw();
    return 0;
}