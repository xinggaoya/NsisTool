## 一键生成 Windows 安装程序

### 简介

本项目提供一个简单易用的工具，可以根据配置文件快速生成专业的 Windows 安装程序。无需编写复杂的 NSIS 脚本，只需配置简单的 JSON 文件，即可生成包含以下功能的安装程序：

-   自定义安装路径
-   创建桌面快捷方式
-   写入注册表信息
-   生成卸载程序
-   设置程序版本号

### 快速使用
1. 直接下载编译好的安装程序，解压后直接运行即可。
2. 为程序配置PATH环境变量，在项目新建nsis.json文件，配置好后运行即可。

### 编译方法

1.  **安装依赖:**

    -   确保系统已安装 [NSIS](https://nsis.sourceforge.io/Download)  (Nullsoft Scriptable Install System)。
    -   安装 Go 语言环境。

2.  **准备配置文件:**

    -   创建一个名为 `nsis.json` 的文件，并按照以下格式配置程序信息：
    ```json
    {
      "name": "程序名称",
      "company": "公司名称",
      "version": "程序版本",
      "icon_path": "图标文件路径",
      "program_path": "程序主文件路径"
    }
    ```

3.  **运行程序:**

    -   使用命令行进入项目目录。
    -   执行 `go run main.go`  命令。

4.  **获取安装程序:**

    -   程序运行完成后，将在项目目录下生成一个可执行的安装程序文件。

### 示例

```json
{
  "name": "MyProgram",
  "company": "MyCompany",
  "version": "1.0.0",
  "icon_path": "path/to/icon.ico",
  "program_path": "path/to/my_program.exe"
}
```

### 注意

-   确保配置文件中的路径信息正确。
-   图标文件必须为 `.ico` 格式。

### 联系方式

如有任何问题或建议，请联系 xinggaoya@qq.com。 
