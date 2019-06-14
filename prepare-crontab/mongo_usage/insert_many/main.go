package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// 日志时间
type LogTime struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

// 日志的记录
type LogRecord struct {
	LogName string   `bson:"logName"`
	Command string   `bson:"command"`
	Err     string   `bson:"err"`
	Content string   `bson:"content"`
	LogTime *LogTime `bson:"logTime"`
}

// 查询的条件
type FindLogName struct {
	LogName string `bson:"logName"`
}

// 删除过滤条件{"$lt":时间}
type BeforeTimeCond struct {
	Before int64 `bson:"$lt"`
}

// 删除条件 {"logTime.startTime":{"$lt":当前时间}}
type DelCond struct {
	BeforeCond BeforeTimeCond `bson:"logTime.startTime"`
}

func main() {
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
		logRecord  *LogRecord
		resultOne  *mongo.InsertOneResult
		logs       []interface{}
		resultMany *mongo.InsertManyResult
		id         interface{}
		cursor     *mongo.Cursor
		cond       *FindLogName
		skip       int64
		limit      int64
		delCond    *DelCond
		delResult  *mongo.DeleteResult
	)

	// 连接mongodb
	if client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://127.0.0.1:27017")); err != nil {
		panic(err)
	}

	defer client.Disconnect(context.TODO())

	// 选择数据库
	database = client.Database("cron")

	// 选择表
	collection = database.Collection("log")

	logRecord = &LogRecord{
		LogName: "echo",
		Command: "echo hello",
		Err:     "",
		Content: "hello",
		LogTime: &LogTime{StartTime: time.Now().Unix(), EndTime: time.Now().Unix() + 10},
	}

	fmt.Println("插入一条数据。")
	// 插入一条记录
	if resultOne, err = collection.InsertOne(context.TODO(), logRecord); err != nil {
		panic(err)
	}

	fmt.Println(resultOne.InsertedID)
	// 插入多条记录
	logs = []interface{}{logRecord, logRecord, logRecord, logRecord, logRecord}
	if resultMany, err = collection.InsertMany(context.TODO(), logs); err != nil {
		panic(err)
	}

	fmt.Println("多次插入打印。")
	for _, id = range resultMany.InsertedIDs {
		fmt.Println(id)
	}

	// 构建查询条件
	cond = &FindLogName{LogName: "echo"}
	skip = 0
	limit = 3

	if cursor, err = collection.Find(context.TODO(), cond, &options.FindOptions{Skip: &skip, Limit: &limit}); err != nil {
		fmt.Println(err)
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		logRecord = &LogRecord{}
		if err = cursor.Decode(logRecord); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(*logRecord)
	}

	// 删除操作
	// 操作开始时间要早于当前时间的所有记录($lt是less than)
	// delete({"logTime.startTime":{"$lt":当前时间}})
	delCond = &DelCond{BeforeCond: BeforeTimeCond{Before: time.Now().Unix()}}

	if delResult, err = collection.DeleteMany(context.TODO(), delCond); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("共删除了：", delResult.DeletedCount)
}
