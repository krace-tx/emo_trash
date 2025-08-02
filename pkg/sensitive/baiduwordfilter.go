package sensitive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func token() {
	url := "https://aip.baidubce.com/oauth/2.0/token?client_id=fDrnSbMmFzNs2UCmm8ElL6cZ&client_secret=CVk7SPHwVvh6qPehjHH43M5Vs1CKvDDV&grant_type=client_credentials"
	payload := strings.NewReader(``)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func Text() {
	accessToken := "24.b22487c02d0c2e295b3588e5523c9851.2592000.1728459472.282335-59603361"
	requestURL := "https://aip.baidubce.com/rest/2.0/solution/v1/text_censor/v2/user_defined?access_token=" + accessToken

	params := map[string]string{
		"text": "抄你妈",
	}

	// 将参数转换为 URL 编码格式
	data := ""
	for key, value := range params {
		if len(data) > 0 {
			data += "&"
		}
		data += fmt.Sprintf("%s=%s", key, value)
	}

	// 创建请求
	req, err := http.NewRequest("POST", requestURL, bytes.NewBufferString(data))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 执行请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error executing request:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// 打印响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return
	}
	fmt.Println(result)
}
