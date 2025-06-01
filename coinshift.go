package main

import (
	"bytes"
	"compress/gzip"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// ANSI 颜色代码
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

// 图标
const (
	IconSuccess = "✅"
	IconError   = "❌"
	IconInfo    = "ℹ️"
	IconWarning = "⚠️"
	IconStart   = "🚀"
	IconKey     = "🔑"
	IconAddress = "🏠"
	IconNetwork = "🌐"
)

// Config 定义配置文件结构
type Config struct {
	Accounts []struct {
		PrivateKey   string `json:"private_key"`
		Proxy        string `json:"proxy"`
		RefreshToken string `json:"refresh_token,omitempty"`
	} `json:"accounts"`
}

// AuthenticateRequest 定义请求结构体
type AuthenticateRequest struct {
	Message          string `json:"message"`
	Signature        string `json:"signature"`
	ChainID          string `json:"chainId"`
	WalletClientType string `json:"walletClientType"`
	ConnectorType    string `json:"connectorType"`
	Mode             string `json:"mode"`
}

// AuthenticateResponse 定义响应结构体
type AuthenticateResponse struct {
	User struct {
		ID             string `json:"id"`
		CreatedAt      int64  `json:"created_at"`
		LinkedAccounts []struct {
			Type              string `json:"type"`
			Address           string `json:"address,omitempty"`
			ChainType         string `json:"chain_type,omitempty"`
			ChainID           string `json:"chain_id,omitempty"`
			WalletClient      string `json:"wallet_client,omitempty"`
			WalletClientType  string `json:"wallet_client_type,omitempty"`
			ConnectorType     string `json:"connector_type,omitempty"`
			VerifiedAt        int64  `json:"verified_at,omitempty"`
			FirstVerifiedAt   int64  `json:"first_verified_at,omitempty"`
			LatestVerifiedAt  int64  `json:"latest_verified_at,omitempty"`
			Subject           string `json:"subject,omitempty"`
			Name              string `json:"name,omitempty"`
			Username          string `json:"username,omitempty"`
			ProfilePictureURL string `json:"profile_picture_url,omitempty"`
			Email             string `json:"email,omitempty"`
		} `json:"linked_accounts"`
		MFAMethods       []interface{} `json:"mfa_methods"`
		HasAcceptedTerms bool          `json:"has_accepted_terms"`
		IsGuest          bool          `json:"is_guest"`
	} `json:"user"`
	Token               string `json:"token"`
	PrivyAccessToken    string `json:"privy_access_token"`
	RefreshToken        string `json:"refresh_token"`
	IdentityToken       string `json:"identity_token"`
	SessionUpdateAction string `json:"session_update_action"`
	IsNewUser           bool   `json:"is_new_user"`
}

// PrivyInitRequest 定义请求结构体
type PrivyInitRequest struct {
	Address string `json:"address"`
}

// PrivyInitResponse 定义响应结构体
type PrivyInitResponse struct {
	Nonce     string `json:"nonce"`
	Address   string `json:"address"`
	ExpiresAt string `json:"expires_at"`
}

// GraphQLRequest 定义 GraphQL 请求结构体
type GraphQLRequest struct {
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
	Query         string                 `json:"query"`
}

// GraphQLResponse 定义 GraphQL 响应结构体
type GraphQLResponse struct {
	Data struct {
		UserLogin string `json:"userLogin"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}
type VerifyActivityResponse struct {
	Data struct {
		VerifyActivity struct {
			Record struct {
				ID            string `json:"id"`
				ActivityID    string `json:"activityId"`
				Status        string `json:"status"`
				Properties    any    `json:"properties"`
				CreatedAt     string `json:"createdAt"`
				RewardRecords []struct {
					ID                    string `json:"id"`
					Status                string `json:"status"`
					AppliedRewardType     string `json:"appliedRewardType"`
					AppliedRewardQuantity int    `json:"appliedRewardQuantity"`
					AppliedRewardMetadata any    `json:"appliedRewardMetadata"`
					Error                 any    `json:"error"`
					RewardID              string `json:"rewardId"`
					Reward                struct {
						ID         string   `json:"id"`
						Quantity   int      `json:"quantity"`
						Type       string   `json:"type"`
						Properties struct{} `json:"properties"`
						Typename   string   `json:"__typename"`
					} `json:"reward"`
					Typename string `json:"__typename"`
				} `json:"rewardRecords"`
				Typename string `json:"__typename"`
			} `json:"record"`
			MissionRecord any    `json:"missionRecord"`
			Typename      string `json:"__typename"`
		} `json:"verifyActivity"`
	} `json:"data"`
}

// 自定义日志函数
func logInfo(format string, v ...interface{}) {
	log.Printf(ColorCyan+IconInfo+" INFO: "+format+ColorReset, v...)
}

func logSuccess(format string, v ...interface{}) {
	log.Printf(ColorGreen+IconSuccess+" SUCCESS: "+format+ColorReset, v...)
}

func logWarning(format string, v ...interface{}) {
	log.Printf(ColorYellow+IconWarning+" WARNING: "+format+ColorReset, v...)
}

func logError(format string, v ...interface{}) {
	log.Printf(ColorRed+IconError+" ERROR: "+format+ColorReset, v...)
}

func logStart(format string, v ...interface{}) {
	log.Printf(ColorBlue+IconStart+" START: "+format+ColorReset, v...)
}

// loadConfig 加载配置文件
func loadConfig(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
}

// createHTTPClient 创建带代理的HTTP客户端
func createHTTPClient(proxyURL string) (*http.Client, error) {
	if proxyURL == "" {
		return &http.Client{Timeout: 10 * time.Second}, nil
	}

	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("解析代理URL失败: %v", err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}

	return &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}, nil
}

// SignEIP4361Message 生成 EIP-4361 签名
func SignEIP4361Message(
	privateKeyHex, domain, address, statement, uri, version, chainID, nonce, issuedAt string,
	resources []string,
) (string, string, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", "", fmt.Errorf("解析私钥失败: %v", err)
	}

	message := formatEIP4361Message(domain, address, statement, uri, version, chainID, nonce, issuedAt, resources)
	hashedMessage := hashMessage(message)

	signature, err := crypto.Sign(hashedMessage, privateKey)
	if err != nil {
		return "", "", fmt.Errorf("签名失败: %v", err)
	}

	signatureHex := hex.EncodeToString(signature)
	return "0x" + signatureHex, message, nil
}

func formatEIP4361Message(domain, address, statement, uri, version, chainID, nonce, issuedAt string, resources []string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s wants you to sign in with your Ethereum account:\n", domain))
	sb.WriteString(fmt.Sprintf("%s\n\n", address))

	if statement != "" {
		sb.WriteString(fmt.Sprintf("%s\n\n", statement))
	} else {
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("URI: %s\n", uri))
	sb.WriteString(fmt.Sprintf("Version: %s\n", version))
	sb.WriteString(fmt.Sprintf("Chain ID: %s\n", chainID))
	sb.WriteString(fmt.Sprintf("Nonce: %s\n", nonce))
	sb.WriteString(fmt.Sprintf("Issued At: %s", issuedAt))

	if len(resources) > 0 {
		sb.WriteString("\nResources:")
		for _, resource := range resources {
			sb.WriteString(fmt.Sprintf("\n- %s", resource))
		}
	}

	return sb.String()
}

func hashMessage(message string) []byte {
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(message))
	data := []byte(prefix + message)
	return crypto.Keccak256([]byte(data))
}

// GetAddressFromPrivateKey 通过私钥获取以太坊钱包地址
func GetAddressFromPrivateKey(privateKeyHex string) (string, error) {
	privateKeyBytes, err := hexutil.Decode(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("无效的私钥格式: %v", err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("私钥转换失败: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("无法获取ECDSA公钥")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address.Hex(), nil
}

// InitPrivyAuth 初始化Privy认证
func InitPrivyAuth(address, proxyURL string) (*PrivyInitResponse, error) {
	url := "https://auth.privy.io/api/v1/siwe/init"
	requestBody, err := json.Marshal(PrivyInitRequest{Address: address})
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	setRequestHeaders(req)

	client, err := createHTTPClient(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP客户端失败: %v", err)
	}

	logInfo("正在初始化 Privy 认证...")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求发送失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("非预期状态码: %d, 响应: %s", resp.StatusCode, body)
	}

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("创建gzip读取器失败: %v", err)
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	var response PrivyInitResponse
	if err := json.NewDecoder(reader).Decode(&response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &response, nil
}

// AuthenticateWithPrivy 向Privy认证服务发送请求
func AuthenticateWithPrivy(request AuthenticateRequest, proxyURL string) (*AuthenticateResponse, error) {
	url := "https://auth.privy.io/api/v1/siwe/authenticate"
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	setRequestHeaders(req)

	client, err := createHTTPClient(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP客户端失败: %v", err)
	}

	logInfo("正在向 Privy 发送认证请求...")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求发送失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("非预期状态码: %d, 响应: %s", resp.StatusCode, body)
	}

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("创建gzip读取器失败: %v", err)
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	var response AuthenticateResponse
	if err := json.NewDecoder(reader).Decode(&response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &response, nil
}

func setRequestHeaders(req *http.Request) {
	headers := map[string]string{
		"accept":             "application/json",
		"accept-encoding":    "gzip, deflate, br, zstd",
		"accept-language":    "zh-CN,zh;q=0.9,en;q=0.8",
		"content-type":       "application/json",
		"origin":             "https://campaign.coinshift.xyz",
		"priority":           "u=1, i",
		"privy-app-id":       "clphlvsh3034xjw0fvs59mrdc",
		"privy-ca-id":        "e37a03d7-0a73-423e-b427-71b288d6c199",
		"privy-client":       "react-auth:2.4.1",
		"referer":            "https://campaign.coinshift.xyz/",
		"sec-ch-ua":          `"Google Chrome";v="137", "Chromium";v="137", "Not/A)Brand";v="24"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"macOS"`,
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "cross-site",
		"user-agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36",
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

func GetCurrentTimeInISO8601() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
}

// DeformLoginRequest 向 deform.cc 发送登录请求
func DeformLoginRequest(authToken, proxyURL string) (string, error) {
	// 1. 准备请求URL
	url := "https://api.deform.cc/"

	// 2. 构造 GraphQL 请求
	requestBody := GraphQLRequest{
		OperationName: "UserLogin",
		Variables: map[string]interface{}{
			"data": map[string]string{
				"externalAuthToken": authToken,
			},
		},
		Query: `mutation UserLogin($data: UserLoginInput!) {
			userLogin(data: $data)
		}`,
	}

	// 3. 序列化请求体
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求体失败: %v", err)
	}

	// 4. 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 5. 设置请求头
	setDeformRequestHeaders(req)

	// 6. 创建HTTP客户端并发送请求
	client, err := createHTTPClient(proxyURL)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求发送失败: %v", err)
	}
	defer resp.Body.Close()

	// 7. 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("非预期状态码: %d, 响应: %s", resp.StatusCode, body)
	}

	// 8. 解析响应体
	var response GraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	// 9. 检查 GraphQL 错误
	if len(response.Errors) > 0 {
		return "", fmt.Errorf("GraphQL错误: %v", response.Errors[0].Message)
	}

	return response.Data.UserLogin, nil
}
func VerifyActivity(activityId, bearerToken, privyIdToken, proxyURL string) (string, error) {
	// 1. 准备请求URL
	uri := "https://api.deform.cc/"

	// 2. 构造 GraphQL 请求
	requestBody := GraphQLRequest{
		OperationName: "VerifyActivity",
		Variables: map[string]interface{}{
			"data": map[string]interface{}{
				"activityId": activityId,
			},
		},
		Query: `mutation VerifyActivity($data: VerifyActivityInput!) {
  verifyActivity(data: $data) {
    record {
      id
      activityId
      status
      properties
      createdAt
      rewardRecords {
        id
        status
        appliedRewardType
        appliedRewardQuantity
        appliedRewardMetadata
        error
        rewardId
        reward {
          id
          quantity
          type
          properties
          __typename
        }
        __typename
      }
      __typename
    }
    missionRecord {
      id
      missionId
      status
      createdAt
      rewardRecords {
        id
        status
        appliedRewardType
        appliedRewardQuantity
        appliedRewardMetadata
        error
        rewardId
        reward {
          id
          quantity
          type
          properties
          __typename
        }
        __typename
      }
      __typename
    }
    __typename
  }
}`,
	}

	// 3. 序列化请求体
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求体失败: %v", err)
	}

	// 4. 创建HTTP请求
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 5. 设置请求头
	setDeformRequestHeaders(req)
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Privy-Id-Token", privyIdToken)
	// 6. 创建HTTP客户端并发送请求
	client, err := createHTTPClient(proxyURL)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求发送失败: %v", err)
	}
	defer resp.Body.Close()

	// 7. 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("非预期状态码: %d, 响应: %s", resp.StatusCode, body)
	}
	// 8. 解析响应体
	var response VerifyActivityResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}
	logInfo("完成任务状态：%s", response.Data.VerifyActivity.Record.Status)
	return response.Data.VerifyActivity.Record.ActivityID, nil
}

// setDeformRequestHeaders 设置 deform.cc 请求头
func setDeformRequestHeaders(req *http.Request) {
	headers := map[string]string{
		"accept":                  "*/*",
		"accept-encoding":         "gzip, deflate, br, zstd",
		"accept-language":         "zh-CN,zh;q=0.9,en;q=0.8",
		"content-type":            "application/json",
		"origin":                  "https://campaign.coinshift.xyz",
		"priority":                "u=1, i",
		"referer":                 "https://campaign.coinshift.xyz/",
		"sec-ch-ua":               `"Google Chrome";v="137", "Chromium";v="137", "Not/A)Brand";v="24"`,
		"sec-ch-ua-mobile":        "?0",
		"sec-ch-ua-platform":      `"macOS"`,
		"sec-fetch-dest":          "empty",
		"sec-fetch-mode":          "cors",
		"sec-fetch-site":          "cross-site",
		"user-agent":              "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36",
		"x-apollo-operation-name": "UserLogin",
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

func main() {
	logStart("     Coinshift 每日签到脚本")
	logStart("欢迎关注「闲菜」矩阵账号获取深度内容：")
	logStart("公众号搜索：「闲菜web3日记」、「加密之友闲菜哥」、「闲菜解码WEB3」")
	logStart("视频号、YouTube：「加密小闲菜」")
	logStart("Twitter：「@xiancai4188391」\n")
	// 定义命令行参数，默认值为 "config.json"
	filename := flag.String("config", "config.json", "配置文件路径")
	flag.Parse()
	// 加载配置文件
	config, err := loadConfig(*filename)
	if err != nil {
		logError("加载配置失败: %v", err)
		return
	}
	logSuccess("成功加载配置文件，共 %d 个账户 ", len(config.Accounts))

	// 遍历所有账户
	for i, account := range config.Accounts {
		logInfo("处理第 %d 个账户 (代理: %s)", i+1, account.Proxy)

		// 获取地址
		address, err := GetAddressFromPrivateKey(account.PrivateKey)
		if err != nil {
			logError("获取地址失败: %v", err)
			continue
		}
		logSuccess("%s 地址: %s", IconAddress, address)

		// 初始化 Privy 认证
		initResponse, err := InitPrivyAuth(address, account.Proxy)
		if err != nil {
			logError("初始化 Privy 认证失败: %v", err)
			continue
		}
		logSuccess("成功获取 Nonce: %s", initResponse.Nonce)

		// 生成签名
		signature, msg, err := SignEIP4361Message(
			account.PrivateKey[2:],
			"campaign.coinshift.xyz",
			address,
			"By signing, you are proving you own this wallet and logging in. This does not initiate a transaction or cost any fees.",
			"https://campaign.coinshift.xyz",
			"1",
			"1",
			initResponse.Nonce,
			GetCurrentTimeInISO8601(),
			[]string{"https://privy.io"},
		)
		if err != nil {
			logError("生成签名失败: %v", err)
			continue
		}
		logSuccess("%s 签名生成成功", IconKey)

		// 认证请求
		authRequest := AuthenticateRequest{
			Message:          msg,
			Signature:        signature,
			ChainID:          "eip155:1",
			WalletClientType: "okx_wallet",
			ConnectorType:    "injected",
			Mode:             "login-or-sign-up",
		}

		authResponse, err := AuthenticateWithPrivy(authRequest, account.Proxy)
		if err != nil {
			logError("认证失败: %v", err)
			continue
		}

		// 打印结果
		logSuccess("%s 认证成功!", IconSuccess)
		logInfo("用户ID: %s", authResponse.User.ID)
		logInfo("访问Token: %s...", authResponse.Token[:30])
		logInfo("刷新Token: %s...", authResponse.RefreshToken[:10])
		logInfo("链接账户数: %d", len(authResponse.User.LinkedAccounts))
		logInfo("是否新用户: %d", authResponse.IsNewUser)

		token, err := DeformLoginRequest(authResponse.Token, account.Proxy)
		if err != nil {
			logError("登录失败: %v", err)
		}
		logSuccess("登录成功! Token: %s\n", token[30:])

		// 活动ID列表
		activityIDs := []string{
			"304a9530-3720-45c8-a778-fbd3060d5cfd",
			"e3e5f263-b471-4ef3-b285-77a66e358a69",
			"907b82a0-152f-45d7-ae35-ce01de22b481",
		}

		// 循环处理每个活动ID
		for _, activityID := range activityIDs {
			activity, err := VerifyActivity(activityID, token, authResponse.IdentityToken, account.Proxy)
			if err != nil {
				logError("活动 %s 领取失败: %v", activityID, err)
			} else {
				logSuccess("活动 %s 领取成功! %s\n", activityID, activity)
			}

			// 可选：添加延迟避免请求过于频繁
			time.Sleep(1 * time.Second)
		}
	}

	logSuccess("所有账户处理完成")
}
