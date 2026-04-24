package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const (
	URI      = "mongodb://root:emoTrashMongo123!@127.0.0.1:27017/emo_trash?authSource=admin"
	Database = "emo_trash"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// ─── 连接 ─────────────────────────────────────────────────────────────────
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer client.Disconnect(ctx)

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("Ping 失败: %v", err)
	}
	fmt.Println("✓ 连接成功")

	db := client.Database(Database)

	// ─── 创建 Collection ──────────────────────────────────────────────────────
	if err = createCollection(ctx, db); err != nil {
		log.Fatalf("创建 Collection 失败: %v", err)
	}

	// ─── 创建索引 ─────────────────────────────────────────────────────────────
	if err = createIndexes(ctx, db); err != nil {
		log.Fatalf("创建索引失败: %v", err)
	}

	// ─── 插入默认数据 ─────────────────────────────────────────────────────────
	if err = seedDefaultUser(ctx, db); err != nil {
		log.Fatalf("插入默认数据失败: %v", err)
	}

	fmt.Println("\n✓ 初始化完成")
}

// ─── 创建 Collection（带 JSON Schema 校验）────────────────────────────────────
func createCollection(ctx context.Context, db *mongo.Database) error {
	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": bson.A{"email", "password", "status", "created_at", "updated_at"},
			"properties": bson.M{
				"email":      bson.M{"bsonType": "string", "maxLength": 128},
				"password":   bson.M{"bsonType": "string", "maxLength": 128},
				"nickname":   bson.M{"bsonType": "string", "maxLength": 32},
				"avatar":     bson.M{"bsonType": "string", "maxLength": 256},
				"status":     bson.M{"bsonType": "int", "enum": bson.A{0, 1}},
				"created_at": bson.M{"bsonType": "date"},
				"updated_at": bson.M{"bsonType": "date"},
				"deleted_at": bson.M{"bsonType": bson.A{"date", "null"}},
			},
		},
	}

	opts := options.CreateCollection().
		SetValidator(validator).
		SetValidationLevel("strict").
		SetValidationAction("error")

	err := db.CreateCollection(ctx, "users", opts)
	if err != nil {
		// 已存在则跳过
		if isCollectionExistsError(err) {
			fmt.Println("⚠ Collection 'users' 已存在，跳过创建")
			return nil
		}
		return err
	}

	fmt.Println("✓ Collection 'users' 创建完成")
	return nil
}

// ─── 创建索引 ─────────────────────────────────────────────────────────────────
func createIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("users")

	indexes := []mongo.IndexModel{
		// email 唯一索引（仅未删除文档）
		{
			Keys: bson.D{{Key: "email", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetPartialFilterExpression(bson.M{"deleted_at": bson.M{"$exists": false}}).
				SetName("uidx_email_active"),
		},
		// status + created_at 复合索引
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().SetName("idx_status_created"),
		},
		// 活跃用户 created_at 部分索引
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().
				SetPartialFilterExpression(bson.M{"deleted_at": bson.M{"$exists": false}}).
				SetName("idx_active_users_created"),
		},
		// updated_at 索引
		{
			Keys:    bson.D{{Key: "updated_at", Value: -1}},
			Options: options.Index().SetName("idx_updated_at"),
		},
	}

	names, err := col.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return err
	}

	for _, name := range names {
		fmt.Printf("✓ 索引 '%s' 创建完成\n", name)
	}
	return nil
}

// ─── 插入默认管理员 ───────────────────────────────────────────────────────────
func seedDefaultUser(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("users")

	// 已存在则跳过
	count, err := col.CountDocuments(ctx, bson.M{"email": "admin@emo-trash.com"})
	if err != nil {
		return err
	}
	if count > 0 {
		fmt.Println("⚠ 默认用户已存在，跳过插入")
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("Admin@123456"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("生成密码 hash 失败: %w", err)
	}

	now := time.Now()
	user := bson.M{
		"_id":        primitive.NewObjectID(),
		"email":      "admin@emo-trash.com",
		"password":   string(hash),
		"nickname":   "管理员",
		"avatar":     "",
		"status":     1,
		"created_at": now,
		"updated_at": now,
	}

	result, err := col.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	fmt.Printf("✓ 默认用户插入完成\n")
	fmt.Printf("  id:       %v\n", result.InsertedID)
	fmt.Printf("  email:    admin@emo-trash.com\n")
	fmt.Printf("  password: Admin@123456  (请登录后立即修改)\n")
	fmt.Printf("  nickname: 管理员\n")
	return nil
}

// ─── 工具函数 ─────────────────────────────────────────────────────────────────
func isCollectionExistsError(err error) bool {
	if cmdErr, ok := err.(mongo.CommandError); ok {
		return cmdErr.Code == 48 // NamespaceExists
	}
	return false
}
