# go_project

1.使用github.com/go-sql-driver/mysql v1.6.0，链接本地数据库，并试验了各种路由及分页

2.使用JWT添加token验证,将token放到了header中

3.logging用于写日志

4.并使用了swagger生成了API说明网页:swag init

    （1）用// @Tags 分组

    （2）携带token

5.使用nginx转发，并学习了负载均衡 /etc/nginx/nginx.conf

    sudo nginx -t//检查配置
    nginx -s reload
    service nginx restart

    #查看nginx的启动用户及使用用户
    ps aux | grep nginx



6.数据库链接执行SQL语句

    (1)执行sql
        db.Exec("insert into tsys_user (comp_cd,user_cd,user_name,user_pwd,user_type,user_phone,user_email) values (?,?,?,?,?,?,?)",
		    compCd, userCd, userName, usePwd, userType, userPhone, userEmail)

    (2)查询
        userId int
        userList []User
        db.Raw("select * from tsys_user where user_id > ?", userId).Scan(&userList)

    (3)修改
        //测试修改功能---一次修改多条
        func TestEdit(updMsg string, idSlice []int64) bool {
            db.Exec("update tsys_user set add_user=? where user_id in (?) ", updMsg, idSlice)
    	    return true
        }

7.redis,重启不会丢失///   命令行：redis-cli（进入命令）quit(推出命令)

    (1)redis设置单独的中间件类，默认存储10分钟

8.Channel了解,必须在开始进行定义并作为参数传递，不能像类一样使用包

9.最新版本vscode与ubuntu系统不兼容，会闪烁

10.闭包和类有些像，获取闭包就像类的实例化

11.用微服务RPC实现登陆日志(写数据库)

12.用容器Docker部署了MongoDB简单连接实验
    //                   conterid servername 管理员权限
    sudo docker exec -it mongo mongo admin

    (2)无法连接--据说要降低版本到mongo4.0

    (3)降低到4.0版本后可以连接，但是查询不到数据

        查询到了：结构体，字段名，首字母大写，注释部分与数据一致 例如：

        type Person struct {
	        ID     bson.ObjectId `bson:"_id,omitempty"` //类型是bson.ObjectId
	        Name   string        `bson:"name"`          
	        Text   string        `bson:"text"`
	        Newcol string        `bson:"newcol"`
        }

13.做了首页及登陆页面，与后段Golang连接
    熟悉了解了一下React,据说是2022年最或前段架构，本来想用这个做前段的发现不太合适，React更注重交互及复用。

未完成：

4.自动生成网页

8.cookies

