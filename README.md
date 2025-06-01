# coinshift-bot

一个用于自动化相关操作的机器人。 (您可以根据实际功能修改此描述)

---
## 功能特性

* 完成coinshift项目日常任务

---
## 环境要求

* **Go >= 1.23**

### 安装 Go (>= 1.23)

您可以从 Go 官方网站下载并安装：[https://go.dev/dl/](https://go.dev/dl/)

以下是在 Linux 和 macOS 上安装 Go 的常用方法（请参考官方文档获取最新和最详细的指引）：

**对于 Linux:**

1.  **下载 Go 安装包:**
    前往 [https://go.dev/dl/](https://go.dev/dl/) 下载最新的 Go 1.23+ Linux tarball (例如 `go1.23.x.linux-amd64.tar.gz`)。
    ```bash
    # 示例：下载 Go (请替换为最新版本)
    # 例如，若最新版为 go1.23.0:
    wget [https://go.dev/dl/go1.23.0.linux-amd64.tar.gz](https://go.dev/dl/go1.23.0.linux-amd64.tar.gz)
    ```

2.  **解压安装包:**
    建议将其解压到 `/usr/local`。
    ```bash
    # 确保替换为下载的实际文件名
    sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
    ```

3.  **配置环境变量:**
    将 Go 的二进制文件路径添加到您的 `PATH` 环境变量中。通常在 `~/.bashrc`, `~/.zshrc` 或 `~/.profile` 文件中添加如下行：
    ```bash
    export PATH=$PATH:/usr/local/go/bin
    export GOPATH=$HOME/go # 可选，Go Modules 已成为主流
    export PATH=$PATH:$GOPATH/bin # 如果设置了 GOPATH
    ```
    然后使配置生效：
    ```bash
    source ~/.bashrc # 或者 source ~/.zshrc, source ~/.profile
    ```

4.  **验证安装:**
    ```bash
    go version
    ```
    应显示您安装的 Go 版本。

**对于 macOS:**

1.  **使用 Homebrew (推荐):**
    ```bash
    brew update
    brew install go
    ```
    Homebrew 会自动处理 PATH 设置。

2.  **下载官方安装包:**
    前往 [https://go.dev/dl/](https://go.dev/dl/) 下载适用于 macOS 的 `.pkg` 安装程序并按提示安装。安装程序通常会自动配置 PATH。

3.  **验证安装:**
    ```bash
    go version
    ```

---
## 安装与设置

1.  **克隆仓库**
    ```bash
    git clone [https://github.com/xk320/coinshift-bot.git](https://github.com/xk320/coinshift-bot.git)
    cd coinshift-bot
    ```

2.  **编译**
    在项目根目录下执行编译命令：
    ```bash
    go build
    ```
    编译成功后，会生成名为 `coinshift` (或与您项目设置相关的名称) 的可执行文件。

---
## 配置

1.  **复制配置文件**
    将配置文件示例复制为正式配置文件：
    ```bash
    cp config.json.example config.json
    ```

2.  **修改配置文件**
    编辑 `config.json` 文件，填入您的账户信息。

    * `private_key`: 您的账户私钥。**请务必妥善保管您的私钥，切勿泄露！**
    * `proxy`: 您的 HTTP 代理地址，例如 `http://username:password@host:port` 或 `http://host:port`。如果不需要代理，可以留空或移除该字段 (请根据实际程序逻辑确认是否支持留空)。

    配置文件结构示例：
    ```json
    {
      "accounts": [
        {
          "private_key": "0xyour_private_key_here_xxxxxxx",
          "proxy": "http://your_proxy_server:port"
        }
        // 如果需要管理多个账户，可以在 accounts 数组中添加更多对象
        // ,
        // {
        //   "private_key": "0xanother_private_key_here_yyyyyyy",
        //   "proxy": "http://another_proxy_server:port"
        // }
      ]
    }
    ```

---
## 运行

完成编译和配置后，在项目根目录下启动机器人。

```bash
./coinshift
