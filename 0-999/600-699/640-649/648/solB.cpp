#include<iostream>

int main(){
   int n,l,c=0,b=0;
   bool exist = false;
   std::cin>>n;
   int array[2*n];
   for(int i=0;i<(n*2);i++){
     std::cin>>array[i];
   }
   int array2[n*2][(n*2)-1];
     for(int i=0;i<n*2;i++){
        for(int j=0;j<n*2;j++){
           if(i!=j){
            array2[c][b]=array[i]+array[j];
            b++;
           }
        } c++;
        b=0;
     }






      for(int k=0;k<(n*2)-1;k++){
         l = array2[0][k];
     for(int i=1;i<n*2;i++){
        for(int j=0;j<(n*2)-1;j++){
           if(l==array2[i][j]){
               exist = true;
               break;
           }
           
        } if(exist == false){
            break;
        }else{
            if(i<(n*2)-1){
            exist = false;
            }else{
                exist = true;
            }
        }
     } 
     if(exist==true){
        break;
     }
}  
   

     for(int i=0;i<n*2;i++){
        if(array[i]!='\0'){ 
        for(int j=i;j<(n*2)-1;j++){
            if(array[i]+array[j+1]==l){
                std::cout<<array[i]<<' '<<array[j+1]<<'\n';
                array[j+1]='\0';
                 break;
            }
        }
        }
     }
    return 0;
}