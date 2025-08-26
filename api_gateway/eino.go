package main

import (
	"context"
	"fmt"
	"log"
)

func EinoClient() {

	ctx := context.Background()

	// 使用模版创建messages
	fmt.Printf("===create messages===\n")
	messages := CreateMessagesFromTemplate()
	fmt.Printf("messages: %+v\n\n", messages)

	// 创建llm
	fmt.Printf("===create llm===\n")
	//cm := createOpenAIChatModel(ctx)
	cm := CreateOllamaChatModel(ctx)
	fmt.Printf("create llm success\n\n")

	fmt.Printf("===llm generate===\n")
	//result := generate(ctx, cm, messages)
	//log.Printf("result: %+v\n\n", result)

	fmt.Printf("===llm stream generate===\n")
	//streamResult, err := cm.Stream(ctx, messages)
	//if err != nil {
	//
	//}
	//reportStream(streamResult)
	//result, err := cm.Generate(ctx, messages)
	//if err != nil {
	//
	//}
	log.Printf("===llm generate===\n")
	result := generate(ctx, cm, messages)
	log.Printf("result: %+v\n\n", result)
}
