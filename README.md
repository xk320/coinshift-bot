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

以下是在不同操作系统上安装 Go 的常用方法（请参考官方文档获取最新和最详细的指引）：

**对于 Linux:**

1.  **下载 Go 安装包:**
    前往 [https://go.dev/dl/](https://go.dev/dl/) 下载最新的 Go 1.23+ Linux tarball (例如 `go1.23.x.linux-amd64.tar.gz`)。
    ```bash
    # 示例：下载 Go (请替换为最新版本)
    # 例如，若最新版为 go1.23.0:
    wget https://go.dev/dl/go1.24.3.linux-amd64.tar.gz
    ```

2.  **解压安装包:**
    建议将其解压到 `/usr/local`。
    ```bash
    # 确保替换为下载的实际文件名
    sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.24.3.linux-amd64.tar.gz
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

**对于 Windows:**

1.  **下载官方安装程序:**
    前往 [https://go.dev/dl/](https://go.dev/dl/) 下载适用于 Windows 的 `.msi` 安装程序（例如 `go1.23.x.windows-amd64.msi`）。

2.  **运行安装程序:**
    双击下载的 `.msi` 文件，并按照安装向导的提示完成安装。默认情况下，Go 会被安装在 `C:\Program Files\Go` 或 `C:\Go` (较新版本倾向于后者)，并且安装程序会自动将 `C:\Go\bin` 添加到系统环境变量 `Path` 中。

3.  **验证安装:**
    打开新的命令提示符 (Command Prompt) 或 PowerShell 窗口，输入：
    ```cmd
    go version
    ```
    应显示您安装的 Go 版本。如果命令未被识别，请检查环境变量 `Path` 是否已正确配置并重启命令提示符/PowerShell。

---
## 安装与设置

1.  **克隆仓库**
    ```bash
    git clone [https://github.com/xk320/coinshift-bot.git](https://github.com/xk320/coinshift-bot.git)
    cd coinshift-bot
    ```
    (在 Windows 上，您可以使用 Git Bash, PowerShell 或命令提示符执行这些命令，前提是已安装 Git。)

2.  **编译**
    在项目根目录下执行编译命令：
    ```bash
    go build
    ```
    * 在 Linux 和 macOS 上，编译成功后，会生成名为 `coinshift` (或与您项目设置相关的名称) 的可执行文件。
    * 在 Windows 上，编译成功后，会生成名为 `coinshift.exe` (或 `项目名.exe`) 的可执行文件。

---
## 配置

1.  **复制配置文件**
    将配置文件示例复制为正式配置文件：
    ```bash
    cp config.json.example config.json
    ```
    (在 Windows 命令提示符下，使用 `copy` 命令: `copy config.json.example config.json`)

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
### 使用 screen 在后台运行 (Linux / macOS)  স্ক্রিন

如果您需要在服务器上或希望关闭终端后机器人仍能持续运行，可以使用 `screen`。

1.  **安装 screen (如果尚未安装):** 🛠️
    * Debian/Ubuntu: `sudo apt update && sudo apt install screen -y`
    * CentOS/RHEL: `sudo yum install screen -y`
    * macOS (通常自带，或通过 Homebrew): `brew install screen`

2.  **创建并进入一个新的 screen 会话:** 🚀
    您可以为会话指定一个名称，方便后续管理，例如 `coinshift_session`。
    ```bash
    screen -S coinshift_session
    ```

3.  **在 screen 会话中运行机器人:** ▶️
    在新的 screen 终端中，导航到项目目录并启动机器人。
    ```bash
    # 示例：假设您的项目在 ~/coinshift-bot
    # cd ~/coinshift-bot
    ./coinshift
    ```

4.  **分离 (Detach) screen 会话:** 💨
    机器人开始运行后，您可以按下 `Ctrl+A` 然后再按 `D`键。这样会话会继续在后台运行，您可以关闭当前终端窗口。

5.  **重新连接 (Reattach) screen 会话:** 🔗
    如果您想查看机器人的输出或停止它，可以重新连接到之前的会话：
    ```bash
    screen -r coinshift_session
    ```
    如果只有一个 screen 会话在运行，也可以简单使用 `screen -r`。

6.  **列出所有 screen 会话:** 📋
    ```bash
    screen -ls
    ```

7.  **终止 screen 会话 (及内部程序):** ⏹️
    重新连接到会话后，使用 `Ctrl+C` 停止机器人，然后输入 `exit` 来关闭 screen 会话。或者，从外部直接杀死会话（不推荐，除非无法重连）：
    ```bash
    screen -X -S coinshift_session quit
    ```

---
### 在 Windows 上后台运行

Windows 没有内置与 `screen` 完全相同的工具，但有几种方法可以实现后台运行：

1.  **使用 `start` 命令 (命令提示符):**
    这个命令可以启动一个新进程而不等待它完成，`/B` 参数表示在不创建新窗口的情况下启动应用程序。但输出可能不会显示，且管理相对困难。
    ```cmd
    start /B .\coinshift.exe
    ```

2.  **使用 PowerShell `Start-Process`:**
    这提供了更多控制，例如可以将窗口隐藏。
    ```powershell
    # 启动进程，不创建新窗口 (输出可能在当前PowerShell窗口，或根据程序设计被重定向)
    Start-Process -FilePath ".\coinshift.exe" -NoNewWindow

    # 或者，如果程序是控制台应用，可以尝试将其输出重定向到文件：
    Start-Process -FilePath ".\coinshift.exe" -RedirectStandardOutput "output.log" -RedirectStandardError "error.log" -NoNewWindow
    ```
    若要完全隐藏窗口并后台运行，且程序本身不创建窗口，`-WindowStyle Hidden` 可能有用 (主要对 GUI 应用或可配置窗口行为的应用)。对于控制台应用，`-NoNewWindow` 通常是使其在当前控制台会话之外运行的方式。

3.  **使用任务计划程序 (Task Scheduler):**
    对于需要长期可靠运行的后台任务，Windows 任务计划程序是更健壮的选择。您可以配置任务在系统启动时运行，或按计划运行，并管理其重启策略等。配置步骤相对复杂，请查阅 Windows 官方文档。

4.  **使用第三方工具:**
    有一些第三方工具（如 NSSM - Non-Sucking Service Manager）可以将普通的可执行文件包装成 Windows 服务，这对于需要像服务一样运行的应用程序非常有用。

**注意:** 后台运行控制台应用程序时，如何处理其输出（日志）非常重要。确保您的应用程序有良好的日志记录机制，或者将输出重定向到文件。

---
## 注意事项 ⚠️

* **安全警告**: 私钥非常敏感，请确保您的运行环境安全，不要在不信任的计算机上运行此机器人，也不要将包含私钥的配置文件提交到公共代码仓库。
* 代理配置请根据您的实际网络环境填写。

---
## 免责声明 📜

本项目仅供学习和技术研究使用，作者不对任何因使用本项目代码造成的直接或间接损失负责。**请您务必在充分了解相关风险后谨慎使用。**

---
## 贡献 (可选) 🤝

欢迎提交 Pull Request 或 Issue 来改进本项目。

---
## 许可证 (可选) 📄

(如果您的项目有许可证，例如 MIT, Apache 2.0 等，请在此处说明)
