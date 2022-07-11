# go_project

1.使用github.com/go-sql-driver/mysql v1.6.0，链接本地数据库，并试验了各种路由及分页

2.使用JWT添加token验证,将token放到了header中

3.logging用于写日志

4.并使用了swagger生成了API说明网页
    （1）用// @Tags 分组
    （2）携带token

5.使用nginx转发，并学习了负载均衡

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
7.redis///   命令行：redis-cli（进入命令）quit(推出命令)

