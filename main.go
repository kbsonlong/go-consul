package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

//curl  http://consul:8500/v1/catalog/service/pci-gateway?filter=ServiceID=="<servide_id>"
// 查询服务注册到consul所在节点
func GetServiceOnNode(service string, ip string) (node_ip string) {

	consul_host := os.Getenv("CONSUL_HOST")
	chek_url := fmt.Sprintf("http://%s/v1/catalog/service/%s?filter=ServiceID==\"%s-%s-8080\"", consul_host, service, service, ip)

	client := &http.Client{}
	request, err := http.NewRequest("GET", chek_url, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(chek_url)
	res, err := client.Do(request)
	fmt.Println(res.Body)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	//服务端返回数据
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	// 解析数组
	var requestBody = make([]map[string]interface{}, 2)
	if err := json.Unmarshal([]byte(b), &requestBody); err != nil {
		fmt.Println(err)
	}
	fmt.Println(requestBody[0]["Address"])
	addr := requestBody[0]["Address"].(string)
	return addr
}

// curl -X PUT http://10.49.5.72:8500/v1/agent/service/deregister/<service_id>
// 下线节点
func main() {

	service := os.Getenv("SERVICE_NAME")

	ip := os.Getenv("MY_POD_IP")
	ipstr := strings.Replace(ip, ".", "-", -1)
	node := GetServiceOnNode(service, ipstr)
	url := fmt.Sprintf("http://%s:8500/v1/agent/service/deregister/%s-%s-8080", node, service, ipstr)
	fmt.Println(url)
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	_, err = client.Do(req)
	if err != nil {
		// handle error
		log.Fatal(err)
	}

}
