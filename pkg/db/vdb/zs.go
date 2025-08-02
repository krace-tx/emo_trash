package vdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

type Index struct {
	Name        string   `json:"name"`
	StorageType string   `json:"storage_type"`
	ShardNum    int      `json:"shard_num"`
	Mappings    Mappings `json:"mappings"`
	Settings    Settings `json:"settings"`
}

type Mappings struct {
	Properties map[string]Propertie `json:"properties"`
}

type Propertie struct {
	Type          string `json:"type"`
	Index         bool   `json:"index"`
	Store         bool   `json:"store"`
	Sortable      bool   `json:"sortable"`
	Highlightable bool   `json:"highlightable"`
}

type Settings struct {
	NumberOfShards   int `json:"number_of_shards"`
	NumberOfReplicas int `json:"number_of_replicas"`
}

type SearchResponse struct {
	Took     int    `json:"took"`
	TimedOut bool   `json:"timed_out"`
	Shards   Shards `json:"_shards"`
	Hits     Hits   `json:"hits"`
}

func (s *SearchResponse) Unmarshal(obj interface{}) ([]interface{}, error) {
	// 获取目标类型的反射值
	targetType := reflect.TypeOf(obj)
	if targetType.Kind() != reflect.Ptr || targetType.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("obj must be a pointer to a struct")
	}

	// 创建结果列表
	results := make([]interface{}, 0, len(s.Hits.Hits))

	// 遍历 Hits.Hits
	for _, hit := range s.Hits.Hits {
		// 创建目标类型的实例
		newObj := reflect.New(targetType.Elem()).Interface()

		// 将 data 反序列化到新实例中
		if hit.Source.Type == reflect.TypeOf(obj).String() {
			dataBytes, err := json.Marshal(hit.Source.Data)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal data: %v", err)
			}
			if err := json.Unmarshal(dataBytes, newObj); err != nil {
				return nil, fmt.Errorf("failed to unmarshal data: %v", err)
			}
			// 将新实例添加到结果列表
			results = append(results, newObj)
		}
	}

	return results, nil
}

type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type Hits struct {
	Total struct {
		Value int `json:"value"`
	} `json:"total"`
	MaxScore float64   `json:"max_score"`
	Hits     []HitItem `json:"hits"`
}

type HitItem struct {
	Index     string    `json:"_index"`
	Type      string    `json:"_type"`
	ID        string    `json:"_id"`
	Score     float64   `json:"_score"`
	Timestamp time.Time `json:"@timestamp"`
	Source    HitSource `json:"_source"`
}

type HitSource struct {
	Timestamp time.Time   `json:"@timestamp"`
	Data      interface{} `json:"data"`
	Type      string      `json:"type"`
}

// Params 表示查询对象
type Params struct {
	SearchType string   `json:"search_type,omitempty"`
	Query      Query    `json:"query,omitempty"`
	SortFields []string `json:"sort_fields,omitempty"`
	From       int32    `json:"from,omitempty"`
	MaxResults int32    `json:"max_results,omitempty"`
	Source     []string `json:"_source,omitempty"`
}

type Query struct {
	Term      string    `json:"term,omitempty"`
	Field     string    `json:"field,omitempty"`
	StartTime time.Time `json:"start_time,omitempty"`
	EndTime   time.Time `json:"end_time,omitempty"`
}

//// Merge 用于合并多个map，返回合并后的map
//func (q *Params) Merge(maps ...map[string]interface{}) {
//	// 初始化一个新的map来存储合并的结果
//	if q.Query == nil {
//		q.Query = make(map[string]interface{})
//	}
//
//	var merge func(existing, new map[string]interface{}) map[string]interface{}
//	merge = func(existing, new map[string]interface{}) map[string]interface{} {
//		for key, value := range new {
//			if existingVal, exists := existing[key]; exists {
//				if existingMap, ok := existingVal.(map[string]interface{}); ok {
//					if newMap, ok := value.(map[string]interface{}); ok {
//						existing[key] = merge(existingMap, newMap) // 递归合并
//					} else {
//						existing[key] = value // 新值不是map，直接替换
//					}
//				} else {
//					existing[key] = value // 现有值不是map，直接替换
//				}
//			} else {
//				// 如果现有map没有这个key，直接添加
//				existing[key] = value
//			}
//		}
//		return existing
//	}
//
//	// 遍历所有传入的map，合并到q.Params
//	for _, m := range maps {
//		q.Query = merge(q.Query, m)
//	}
//}

// ZincClient  客户端
type ZincClient struct {
	BaseURL  string
	username string
	password string
	cli      *http.Client
	params   *Params
}

func NewZincSearchClient(baseURL, username, password string) *ZincClient {
	return &ZincClient{
		BaseURL:  baseURL,
		username: username,
		password: password,
		cli:      &http.Client{},
	}
}

// 通用请求方法
func (client *ZincClient) doRequest(method, url string, body interface{}) ([]byte, error) {
	var requestBody []byte
	var err error
	if body != nil {
		if reflect.TypeOf(body).Kind() == reflect.Struct || reflect.TypeOf(body).Kind() == reflect.Ptr || reflect.TypeOf(body).Kind() == reflect.Map {
			requestBody, err = json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request body: %v", err)
			}
		} else {
			requestBody = body.([]byte)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.SetBasicAuth(client.username, client.password)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed, status: %s", res.Status)
	}

	return ioutil.ReadAll(res.Body)
}

// 创建索引
func (client *ZincClient) CreateIndex(index Index) error {
	url := fmt.Sprintf("%s/api/index", client.BaseURL)
	_, err := client.doRequest(http.MethodPut, url, index)
	return err
}

// 插入文档
func (client *ZincClient) Insert(indexName, docID string, document interface{}) error {
	docType := reflect.TypeOf(document).String()

	url := fmt.Sprintf("%s/api/%s/_doc", client.BaseURL, indexName)

	body := map[string]interface{}{
		"_id":  docType + ":" + docID,
		"type": docType,
		"data": document,
	}

	_, err := client.doRequest("PUT", url, body)
	return err
}

// 更新文档
func (client *ZincClient) Update(indexName, docID string, document interface{}) error {
	docType := reflect.TypeOf(document).String()

	url := fmt.Sprintf("%s/api/%s/_doc/%s", client.BaseURL, indexName, docType+":"+docID)

	body := map[string]interface{}{
		"type": docType,
		"data": document,
	}

	_, err := client.doRequest("PUT", url, body)
	return err
}

// 删除文档
func (client *ZincClient) Delete(indexName, docID string, document interface{}) error {
	docType := reflect.TypeOf(document).String()

	url := fmt.Sprintf("%s/api/%s/_doc/%s", client.BaseURL, indexName, docType+":"+docID)
	_, err := client.doRequest("DELETE", url, nil)
	return err
}

// 查询文档
/*
{
    "search_type": "match",
    "params": {
        "term": "20232127687",
        "start_time": "2021-12-25T15:08:48.777Z",
        "end_time": "2025-12-28T16:08:48.777Z"
    },
    "sort_fields": ["-@timestamp"],
    "from": 0,
    "max_results": 20
}
*/
func (client *ZincClient) Search(indexName string) (*SearchResponse, error) {
	url := fmt.Sprintf("%s/api/%s/_search", client.BaseURL, indexName)
	if client.params == nil {
		client.params = &Params{
			SearchType: "",
			Query: Query{
				Term:      "",
				Field:     "",
				StartTime: time.Time{},
				EndTime:   time.Time{},
			},
			SortFields: nil,
			From:       0,
			MaxResults: 0,
			Source:     nil,
		}
	}
	defer func() {
		client.params = nil
	}()

	res, err := client.doRequest("POST", url, client.params)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	response := &SearchResponse{}
	err = json.Unmarshal(res, response)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	return response, nil
}

func (client *ZincClient) SearchByParams(indexName string, params *Params) (*SearchResponse, error) {
	if client.params == nil {
		client.params = params
	}
	return client.Search(indexName)
}

//
//// 新建 Match 查询，支持链式调用
//func (client *ZincClient) Match(field, value string) *ZincClient {
//	if client.params == nil {
//		client.params = &Params{
//			Query: map[string]interface{}{},
//		}
//	}
//	client.params.Merge(map[string]interface{}{
//		Match: map[string]interface{}{
//			field: value,
//		},
//	})
//	return client
//}
//
//// 新建 MultiMatch 查询，支持链式调用
//func (client *ZincClient) MultiMatch(query string, fields ...string) *ZincClient {
//	matchQuery := map[string]interface{}{
//		MultiMatch: map[string]interface{}{
//			"params": query,
//		},
//	}
//
//	if len(fields) > 0 {
//		matchQuery[MultiMatch].(map[string]interface{})["fields"] = fields
//	}
//
//	client.params = &Params{
//		Query: map[string]interface{}{
//			"params": matchQuery,
//		},
//	}
//	return client
//}
//
//// 新建 Range 查询，支持链式调用
//func (client *ZincClient) Range(field string, gte, lte interface{}) *ZincClient {
//	client.params = &Params{
//		Query: map[string]interface{}{
//			Range: map[string]interface{}{
//				field: map[string]interface{}{
//					"gte": gte,
//					"lte": lte,
//				},
//			},
//		},
//	}
//	return client
//}
//
//// 新建 MatchAll 查询，支持链式调用
//func (client *ZincClient) MatchAll() *ZincClient {
//	client.params = &Params{
//		Query: map[string]interface{}{
//			MatchAll: struct{}{},
//		},
//	}
//	return client
//}
//
//// 新建 And 条件查询，支持链式调用
//func (client *ZincClient) And(conditions map[string]string) *ZincClient {
//	matchQuery := make([]map[string]interface{}, 0)
//	for field, value := range conditions {
//		matchQuery = append(matchQuery, map[string]interface{}{
//			Match: map[string]interface{}{
//				field: value,
//			},
//		})
//	}
//
//	client.params = &Params{
//		Query: map[string]interface{}{
//			Bool: map[string]interface{}{
//				"must": matchQuery,
//			},
//		},
//	}
//	return client
//}
//
//// 新建 Or 条件查询，支持链式调用
//func (client *ZincClient) Or(conditions map[string]string) *ZincClient {
//	matchQuery := make([]map[string]interface{}, 0)
//	for field, value := range conditions {
//		matchQuery = append(matchQuery, map[string]interface{}{
//			Match: map[string]interface{}{
//				field: value,
//			},
//		})
//	}
//
//	client.params = &Params{
//		Query: map[string]interface{}{
//			Bool: map[string]interface{}{
//				"should": matchQuery,
//			},
//		},
//	}
//	return client
//}
//
//// 新建 Term 查询，支持链式调用
//func (client *ZincClient) Term(value string) *ZincClient {
//	if client.params == nil {
//		client.params = &Params{
//			Query: map[string]interface{}{},
//		}
//	}
//	client.params.Merge(map[string]interface{}{
//		Term: value,
//	})
//
//	return client
//}
//
//// 新建 Prefix 查询，支持链式调用
//func (client *ZincClient) Prefix(field, value string) *ZincClient {
//	client.params = &Params{
//		Query: map[string]interface{}{
//			Prefix: map[string]interface{}{
//				field: value,
//			},
//		},
//	}
//	return client
//}
//
//// 新建 Wildcard 查询，支持链式调用
//func (client *ZincClient) Wildcard(field, pattern string) *ZincClient {
//	client.params = &Params{
//		Query: map[string]interface{}{
//			Wildcard: map[string]interface{}{
//				field: pattern,
//			},
//		},
//	}
//	return client
//}
//
//// 新建 Fuzzy 查询，支持链式调用
//func (client *ZincClient) Fuzzy(field, value string) *ZincClient {
//	client.params = &Params{
//		Query: map[string]interface{}{
//			Fuzzy: map[string]interface{}{
//				field: value,
//			},
//		},
//	}
//	return client
//}
//
//// 新建 Exists 查询，支持链式调用
//func (client *ZincClient) Exists(field string) *ZincClient {
//	client.params = &Params{
//		Query: map[string]interface{}{
//			Exists: map[string]interface{}{
//				"field": field,
//			},
//		},
//	}
//	return client
//}
//
//// 新建 Match Phrase 查询，支持链式调用
//func (client *ZincClient) MatchPhrase(field, value string) *ZincClient {
//	client.params = &Params{
//		Query: map[string]interface{}{
//			MatchPhrase: map[string]interface{}{
//				field: value,
//			},
//		},
//	}
//	return client
//}
