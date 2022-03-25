# gin-pagination
pagination with gin and gorm2
使用  
go get github.com/wujiyu98/gin-pagination 

p := pagination.Init(1,10, 5,100, "/") 
fmt.Println(p.GetList())  
fmt.Println(p.BsPage())  


BsPage() 是bootstrap5 版本的完整分页  
SimpleBsPage() 是bootstrap5 版本的简单分页,适用于手机端  

使用地址：  
github.com/wujiyu98/gin-pagination
