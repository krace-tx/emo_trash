package vdb

import (
	"context"
	"fmt"
)

// Vectorizer 文本向量化接口，支持不同模型实现
type Vectorizer interface {
	// Vectorize 将文本转换为向量
	Vectorize(text string) ([]float32, error)
}

// RAGVectorDB RAG向量数据库结构体
type RAGVectorDB struct {
	esClient   *ElasticsearchClient
	vectorizer Vectorizer // 文本向量化器
	dimensions int        // 向量维度
	ctx        context.Context
}

// NewRAGVectorDB 创建RAG向量数据库实例
func NewRAGVectorDB(esClient *ElasticsearchClient, vectorizer Vectorizer, dimensions int) *RAGVectorDB {
	return &RAGVectorDB{
		esClient:   esClient,
		vectorizer: vectorizer,
		dimensions: dimensions,
		ctx:        context.Background(),
	}
}

// Document 向量数据库文档结构
type Document struct {
	ID      string                 `json:"id"`
	Content string                 `json:"content"`  // 原始文本内容
	Vector  []float32              `json:"vector"`   // 文本向量
	Meta    map[string]interface{} `json:"meta"`     // 文档元数据
	ChunkID string                 `json:"chunk_id"` // 文档分块ID（用于长文档）
}

// CreateVectorIndex 创建向量索引（包含向量字段映射）
func (r *RAGVectorDB) CreateVectorIndex() error {
	// 定义向量索引映射，包含向量字段和文本字段
	mappings := map[string]interface{}{
		"properties": map[string]interface{}{
			"content": map[string]interface{}{
				"type": "text", // 文本内容，用于展示和全文检索
			},
			"vector": map[string]interface{}{
				"type":       "dense_vector",
				"dims":       r.dimensions, // 向量维度
				"index":      true,         // 支持向量检索
				"similarity": "cosine",     // 使用余弦相似度
			},
			"meta": map[string]interface{}{
				"type": "object", // 元数据字段
			},
			"chunk_id": map[string]interface{}{
				"type": "keyword", // 分块ID，用于精确匹配
			},
		},
	}

	// 调用已有的ES客户端创建索引
	return r.esClient.CreateIndex(mappings)
}

// AddDocument 添加文档到向量数据库（自动向量化）
func (r *RAGVectorDB) AddDocument(doc *Document) error {
	// 如果向量为空，则自动向量化
	if doc.Vector == nil || len(doc.Vector) == 0 {
		vec, err := r.vectorizer.Vectorize(doc.Content)
		if err != nil {
			return fmt.Errorf("文本向量化失败: %w", err)
		}
		doc.Vector = vec
	}

	// 检查向量维度是否匹配
	if len(doc.Vector) != r.dimensions {
		return fmt.Errorf("向量维度不匹配，期望: %d, 实际: %d", r.dimensions, len(doc.Vector))
	}

	// 使用ES客户端索引文档
	return r.esClient.IndexDocument(doc.ID, doc)
}

// SearchSimilar 搜索相似文档
// queryText: 查询文本
// topK: 返回前K个相似结果
// out: 用于接收结果的结构体指针
func (r *RAGVectorDB) SearchSimilar(queryText string, topK int, out interface{}) error {
	// 将查询文本向量化
	queryVec, err := r.vectorizer.Vectorize(queryText)
	if err != nil {
		return fmt.Errorf("查询文本向量化失败: %w", err)
	}

	// 构建向量相似度查询
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"script_score": map[string]interface{}{
				"query": map[string]interface{}{
					"match_all": map[string]interface{}{},
				},
				"script": map[string]interface{}{
					"source": "cosineSimilarity(params.query_vector, 'vector') + 1.0",
					"params": map[string]interface{}{
						"query_vector": queryVec,
					},
				},
			},
		},
		"size": topK,
	}

	// 使用ES客户端执行搜索
	return r.esClient.Search(query, out)
}

// UpdateDocument 更新向量数据库中的文档
func (r *RAGVectorDB) UpdateDocument(doc *Document) error {
	// 如果提供了新内容，重新向量化
	if doc.Content != "" {
		vec, err := r.vectorizer.Vectorize(doc.Content)
		if err != nil {
			return fmt.Errorf("文本向量化失败: %w", err)
		}
		doc.Vector = vec
	}

	// 构建更新内容
	updateDoc := map[string]interface{}{}
	if doc.Content != "" {
		updateDoc["content"] = doc.Content
	}
	if doc.Vector != nil {
		updateDoc["vector"] = doc.Vector
	}
	if doc.Meta != nil {
		updateDoc["meta"] = doc.Meta
	}
	if doc.ChunkID != "" {
		updateDoc["chunk_id"] = doc.ChunkID
	}

	// 使用ES客户端更新文档
	return r.esClient.UpdateDocument(doc.ID, updateDoc)
}

// DeleteDocument 从向量数据库删除文档
func (r *RAGVectorDB) DeleteDocument(id string) error {
	return r.esClient.DeleteDocument(id)
}

// SearchByMetadata 结合元数据过滤的相似性搜索
func (r *RAGVectorDB) SearchByMetadata(queryText string, metadata map[string]interface{}, topK int, out interface{}) error {
	// 将查询文本向量化
	queryVec, err := r.vectorizer.Vectorize(queryText)
	if err != nil {
		return fmt.Errorf("查询文本向量化失败: %w", err)
	}

	// 构建带元数据过滤的向量查询
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"script_score": map[string]interface{}{
						"query": map[string]interface{}{
							"match_all": map[string]interface{}{},
						},
						"script": map[string]interface{}{
							"source": "cosineSimilarity(params.query_vector, 'vector') + 1.0",
							"params": map[string]interface{}{
								"query_vector": queryVec,
							},
						},
					},
				},
				"filter": buildMetadataFilters(metadata),
			},
		},
		"size": topK,
	}

	return r.esClient.Search(query, out)
}

// 构建元数据过滤条件
func buildMetadataFilters(metadata map[string]interface{}) []map[string]interface{} {
	filters := make([]map[string]interface{}, 0, len(metadata))
	for k, v := range metadata {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				fmt.Sprintf("meta.%s", k): v,
			},
		})
	}
	return filters
}
