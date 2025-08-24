package vdb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/zeromicro/go-zero/core/logx"
)

// ElasticsearchClient ES客户端结构体，管理连接和索引信息
type ElasticsearchClient struct {
	client *elasticsearch.Client // ES客户端实例
	index  string                // 默认操作索引
	ctx    context.Context       // 上下文（用于超时控制）
}

// NewElasticsearchClient 创建ES客户端实例
// addrs: ES集群地址列表（如["http://127.0.0.1:9200"]）
// index: 默认操作索引名
func NewElasticsearchClient(addrs []string, index string) (*ElasticsearchClient, error) {
	// 配置ES客户端
	cfg := elasticsearch.Config{
		Addresses: addrs,
		Username:  "user",
		Password:  "pass",
	}

	// 创建客户端
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("初始化ES客户端失败: %w", err)
	}

	return &ElasticsearchClient{
		client: client,
		index:  index,
		ctx:    context.Background(),
	}, nil
}

// CreateIndex 创建索引并设置映射
// mappings: 索引映射配置（JSON结构，如{"properties": {"field": {"type": "text"}}}）
func (c *ElasticsearchClient) CreateIndex(mappings map[string]interface{}) error {
	// 检查索引是否已存在
	existsReq := esapi.IndicesExistsRequest{Index: []string{c.index}}
	existsRes, err := existsReq.Do(c.ctx, c.client)
	if err != nil {
		return fmt.Errorf("检查索引存在性失败: %w", err)
	}
	defer existsRes.Body.Close()

	if existsRes.StatusCode == 200 {
		return nil // 索引已存在，无需创建
	}

	// 序列化映射配置
	mappingsJSON, err := json.Marshal(mappings)
	if err != nil {
		return fmt.Errorf("映射配置序列化失败: %w", err)
	}

	// 创建索引请求
	createReq := esapi.IndicesCreateRequest{
		Index: c.index,
		Body:  bytes.NewReader(mappingsJSON),
	}

	res, err := createReq.Do(c.ctx, c.client)
	if err != nil {
		return fmt.Errorf("创建索引请求失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("创建索引失败: %s", res.Status())
	}
	return nil
}

// IndexDocument 索引文档（新增/更新）
// id: 文档ID（为空时ES自动生成）
// doc: 文档内容（结构体或map）
func (c *ElasticsearchClient) IndexDocument(id string, doc interface{}) error {
	// 序列化文档内容
	docJSON, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("文档序列化失败: %w", err)
	}

	// 构建索引请求
	req := esapi.IndexRequest{
		Index:      c.index,
		DocumentID: id,
		Body:       bytes.NewReader(docJSON),
		Refresh:    "false", // 可选：设置为"true"立即刷新索引（影响性能）
	}

	// 执行请求
	res, err := req.Do(c.ctx, c.client)
	if err != nil {
		return fmt.Errorf("索引请求执行失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("文档索引失败: %s", res.Status())
	}
	return nil
}

// Search 执行搜索查询
// query: 查询条件（JSON结构，如{"query": {"match": {"field": "value"}}}）
// out: 用于接收结果的结构体指针（如&struct{ Hits struct{Hits []Hit} }{}}）
func (c *ElasticsearchClient) Search(query map[string]interface{}, out interface{}) error {
	// 序列化查询条件
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("查询条件序列化失败: %w", err)
	}

	// 构建搜索请求
	req := esapi.SearchRequest{
		Index: []string{c.index},
		Body:  bytes.NewReader(queryJSON),
	}

	// 执行请求
	res, err := req.Do(c.ctx, c.client)
	if err != nil {
		return fmt.Errorf("搜索请求执行失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("搜索执行失败: %s", res.Status())
	}

	// 解析响应结果到目标结构体
	if err := json.NewDecoder(res.Body).Decode(out); err != nil {
		return fmt.Errorf("响应结果解析失败: %w", err)
	}
	return nil
}

// GetDocument 根据ID获取文档
// id: 文档ID
// out: 用于接收文档内容的结构体指针
func (c *ElasticsearchClient) GetDocument(id string, out interface{}) error {
	// 构建获取请求
	req := esapi.GetRequest{
		Index:      c.index,
		DocumentID: id,
	}

	res, err := req.Do(c.ctx, c.client)
	if err != nil {
		return fmt.Errorf("获取文档请求失败: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return fmt.Errorf("文档不存在 (ID: %s)", id)
	}

	if res.IsError() {
		return fmt.Errorf("获取文档失败: %s", res.Status())
	}

	// 解析响应结构
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return fmt.Errorf("响应解析失败: %w", err)
	}

	// 提取_source字段并反序列化到目标结构体
	sourceJSON, err := json.Marshal(result["_source"])
	if err != nil {
		return fmt.Errorf("文档内容序列化失败: %w", err)
	}
	if err := json.Unmarshal(sourceJSON, out); err != nil {
		return fmt.Errorf("文档内容反序列化失败: %w", err)
	}
	return nil
}

// UpdateDocument 更新文档（部分更新）
// id: 文档ID
// doc: 要更新的字段（部分文档）
func (c *ElasticsearchClient) UpdateDocument(id string, doc interface{}) error {
	// 构建更新请求体（ES要求格式：{"doc": {...}}）
	updateBody := map[string]interface{}{"doc": doc}
	updateJSON, err := json.Marshal(updateBody)
	if err != nil {
		return fmt.Errorf("更新内容序列化失败: %w", err)
	}

	// 构建更新请求
	req := esapi.UpdateRequest{
		Index:      c.index,
		DocumentID: id,
		Body:       bytes.NewReader(updateJSON),
	}

	res, err := req.Do(c.ctx, c.client)
	if err != nil {
		return fmt.Errorf("更新请求执行失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("文档更新失败: %s", res.Status())
	}
	return nil
}

// DeleteDocument 根据ID删除文档
func (c *ElasticsearchClient) DeleteDocument(id string) error {
	// 构建删除请求
	req := esapi.DeleteRequest{
		Index:      c.index,
		DocumentID: id,
	}

	res, err := req.Do(c.ctx, c.client)
	if err != nil {
		return fmt.Errorf("删除请求执行失败: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return fmt.Errorf("文档不存在 (ID: %s)", id)
	}

	if res.IsError() {
		return fmt.Errorf("文档删除失败: %s", res.Status())
	}
	return nil
}

// Close 关闭客户端（ES客户端无显式关闭方法，此处为兼容接口）
func (c *ElasticsearchClient) Close() error {
	// 官方客户端无需显式关闭，HTTP连接由底层管理
	logx.Info("ES客户端连接已释放")
	return nil
}
