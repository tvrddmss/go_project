package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// 声明服务端
type Mongo struct {
	name string
}

// func creatMongo() *Mongo {
// 	var m Mongo
// 	m.name = "ceshi"
// 	return &m
// }

func (m *Mongo) Init() {

	fmt.Println("mongo连接测试1")
	TestConn1()
	fmt.Println("mongo连接测试2")
	TestConn2()
}

// 声明结构体
type Test struct {
	ID     bson.ObjectId `bson:"_id,omitempty"` //类型是bson.ObjectId
	Name   string        `bson:"name"`          //这里变量名和数据库里的名字不一致
	Text   string        `bson:"text"`
	Newcol string        `bson:"test"`
}

func TestConn1() {
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI("mongodb://admin:123456@192.168.50.8:27017")

	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	find, err := client.Database("test").Collection("test").Find(context.TODO(), bson.M{"name": "Ale"})
	defer find.Close(context.TODO())
	res := []Test{}
	// 遍历游标获取查询到的文档
	for find.Next(context.TODO()) {
		var cur Test
		// 解码find当前行 存到cur
		err = find.Decode(&cur)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, cur)
	}
	fmt.Printf("Test:%s\n", res[len(res)-1].Name)
}

func TestConn2() {

	dialInfo := &mgo.DialInfo{
		Addrs:     []string{"192.168.50.8:27017"}, //远程(或本地)服务器地址及端口号
		Direct:    true,
		Timeout:   time.Second * 1,
		Database:  "admin", //数据库
		Source:    "",
		Username:  "admin",
		Password:  "123456",
		PoolLimit: 4096, // Session.SetPoolLimit
		//Mechanism: "SCRAM-SHA-1",
		//Dial: mgo.Dial(url)
	}
	//url := "mongodb://admin:123456@192.168.50.8:27017"
	//url := "mongodb%3A%2F%2Fadmin%3A123456%40192.168.50.8%3A27017%2Ftest"
	//session, err := mgo.Dial(url)
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("test")
	//c.Find(bson.M{"name": "Ale"}).One(&result)
	// err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
	// 	&Person{"Cla", "+55 53 8402 8510"})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	result := Person{}
	err = c.Find(bson.M{}).One(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Phone:", result.Newcol)
}

type Person struct {
	ID     bson.ObjectId `bson:"_id,omitempty"` //类型是bson.ObjectId
	Name   string        `bson:"name"`          //这里变量名和数据库里的名字不一致
	Text   string        `bson:"text"`
	Newcol string        `bson:"newcol"`
}
