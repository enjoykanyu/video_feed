package main

import (
	"api_gateway/user"
	"context"
	"github.com/cloudwego/eino-ext/components/embedding/ollama"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/schema"

	//"github.com/cloudwego/eino-ext/components/model/ollama"
	"log"
	"os"
	"time"
	"video_douyin/middleware"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/milvus-io/milvus-sdk-go/v2/client"

	"github.com/cloudwego/eino-ext/components/indexer/milvus"
)

func main() {
	// 初始化 Hertz 服务器实例，监听本地 8889 端口
	// WithHostPorts 配置项指定服务监听地址:ml-citation{ref="6,7" data="citationList"}
	hz := server.New(server.WithHostPorts("localhost:8889"))
	// 注册全局中间件
	hz.Use(middleware.AuthMiddleware())
	// 创建路由分组（符合 RESTful 风格）
	// 第一级分组 /douyin 作为 API 根路径:ml-citation{ref="2" data="citationList"}
	douyin := hz.Group("/douyin")
	// 第二级分组 /user 用于用户相关操作:ml-citation{ref="3" data="citationList"}
	userGroup := douyin.Group("/user")
	// 注册用户注册接口，POST 方法对应创建操作
	// 处理函数指向 user 包的 Handler 方法:ml-citation{ref="3" data="citationList"}
	userGroup.POST("/register", user.Register)
	userGroup.POST("/login/code", user.Login)
	userGroup.POST("/code", user.Code)
	// 启动 HTTP 服务，若失败则记录错误日志
	// Run() 会阻塞直到服务终止:ml-citation{ref="6" data="citationList"}
	//引入eino
	//ctx := context.Background()
	//
	//// 使用模版创建messages
	//fmt.Printf("===create messages===\n")
	//messages := CreateMessagesFromTemplate()
	//fmt.Printf("messages: %+v\n\n", messages)
	//
	//// 创建llm
	//fmt.Printf("===create llm===\n")
	////cm := createOpenAIChatModel(ctx)
	//cm := CreateOllamaChatModel(ctx)
	//fmt.Printf("create llm success\n\n")
	//
	//fmt.Printf("===llm generate===\n")
	////result := generate(ctx, cm, messages)
	////log.Printf("result: %+v\n\n", result)
	//
	//fmt.Printf("===llm stream generate===\n")
	////streamResult, err := cm.Stream(ctx, messages)
	////if err != nil {
	////
	////}
	////reportStream(streamResult)
	////result, err := cm.Generate(ctx, messages)
	////if err != nil {
	////
	////}
	//log.Printf("===llm generate===\n")
	//result := generate(ctx, cm, messages)
	//log.Printf("result: %+v\n\n", result)
	embeddingTest()
	//启动hertz服务
	if err := hz.Run(); err != nil {
		log.Fatal(err)
	}
}

func embeddingTest() {
	ctx := context.Background()

	baseURL := os.Getenv("OLLAMA_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:11434" // 默认本地
	}
	model := os.Getenv("OLLAMA_EMBED_MODEL")
	if model == "" {
		model = "nomic-embed-text"
	}

	embedder, err := ollama.NewEmbedder(ctx, &ollama.EmbeddingConfig{
		BaseURL: baseURL,
		Model:   model,
		Timeout: 10 * time.Second,
	})
	if err != nil {
		log.Fatalf("NewEmbedder of ollama error: %v", err)
		return
	}

	log.Printf("===== call Embedder directly =====")

	vectors, err := embedder.EmbedStrings(ctx, []string{"hello", "how are you"})
	if err != nil {
		log.Fatalf("EmbedStrings of Ollama failed, err=%v", err)
	}

	log.Printf("vectors : %v", vectors)

	// you can use WithModel to specify the model
	vectors, err = embedder.EmbedStrings(ctx, []string{"hello", "how are you"}, embedding.WithModel(model))
	if err != nil {
		log.Fatalf("EmbedStrings of Ollama failed, err=%v", err)
	}

	log.Printf("vectors : %v", vectors)
}

func Indexer() {
	// Get the environment variables
	addr := os.Getenv("MILVUS_ADDR")
	username := os.Getenv("MILVUS_USERNAME")
	password := os.Getenv("MILVUS_PASSWORD")
	arkApiKey := os.Getenv("ARK_API_KEY")
	arkModel := os.Getenv("ARK_MODEL")

	// Create a client
	ctx := context.Background()
	cli, err := client.NewClient(ctx, client.Config{
		Address:  addr,
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return
	}
	defer cli.Close()

	// Create an embedding model
	emb, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey: arkApiKey,
		Model:  arkModel,
	})
	if err != nil {
		log.Fatalf("Failed to create embedding: %v", err)
		return
	}

	// Create an indexer
	indexer, err := milvus.NewIndexer(ctx, &milvus.IndexerConfig{
		Client:    cli,
		Embedding: emb,
	})
	if err != nil {
		log.Fatalf("Failed to create indexer: %v", err)
		return
	}
	log.Printf("Indexer created success")

	// Store documents
	docs := []*schema.Document{
		{
			ID:      "milvus-1",
			Content: "milvus is an open-source vector database",
			MetaData: map[string]any{
				"h1": "milvus",
				"h2": "open-source",
				"h3": "vector database",
			},
		},
		{
			ID:      "milvus-2",
			Content: "milvus is a distributed vector database",
		},
	}
	ids, err := indexer.Store(ctx, docs)
	if err != nil {
		log.Fatalf("Failed to store: %v", err)
		return
	}
	log.Printf("Store success, ids: %v", ids)
}
