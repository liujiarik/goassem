# Go 集成打包工具：goassem
goassem只需要一个`package`指令就可以将脚本、配置文件、文档、可执行文件等等打包一个压缩包，同时goassem支持多个平台的交叉编译，便于Release不同平台bin包。

## 安装
```
go get -u github.com/liujiarik/goassem
```

 查看简要使用说明
```
goassem -help
```
## 初始化
goassem 使用 assembly.json 文件对项目进行集成打包。类似maven中pom.xml，我们需要在Go项目根目录下，新建一个assembly.json文件。
在Go项目根路径中，使用
```
goassem init
```

会在当前目录创建一个模板assembly.json文件。Goer 可以在assembly.json中自定义打包逻辑。如果assembly.json已经存在，则提示已经init过。

## assembly.json 配置说明

assembly.json 描述了一组打包逻辑。文件是一个json数组，它支持多重打包，因为在Go项目我们常常希望build不同的Go可执行文件。
常见assembly.json的样例如下，这个assembly.json只打了一个包，如果需要打多个，只需要在后面添加类似的配置结构

``` json
[
  {
    "name": "demo",
    "version": "1.0.1",
    "format": "zip",
    "main": "main",
    "binDir": "bin",
    "buildArgs": [],
    "fileSets": [
      {
        "directory": "conf",
        "outputDirectory": "conf",
        "includes": [
          "ss.yaml"
        ]
      },
      {
        "directory": "sh",
        "outputDirectory": "/",
        "includes": [
          "start.sh",
          "stop.sh"
        ]
      },
      {
        "directory": "./",
        "outputDirectory": "/",
        "includes": [
          "README"
        ]
      }
    ],
    "platforms": [
      {
        "arch": "amd64",
        "os": "linux"
      },
      {
        "arch": "amd64",
        "os": "darwin"
      }
    ]
  }
]
```


参数说明：
**name**：包名
**version**：版本
**format**：压缩类型（目前支持zip,以后会支持其他类型）
**main**：main函数入口Go文件， go build 的目标文件。文件地址相对于 **Go工程根目录**。**注意不要加`.go`后缀**。

 **binDir** ：可执行文件的地址，文件地址相对于**包目录**
**buildArgs** ：go build 的编译参数

**fileSets** ：文件打包处理命令，是一个数组格式，可以处理多文件。
* directory：源文件目录 ，文件地址相对于 **Go工程根目录**。如果是Go工程根目录请使用“./”
* outputDirectory：输出目录，文件地址相对于**包目录**。如果是包目录，请使用“/”
* includes：需要处理的文件。
fileSets 将directory 下 includes包含的文件拷贝到outputDirectory目录下。通过fileSets可以将文档，脚本，配置文件打包到包内的任意地址。
**platforms** ：交叉编译支持，是一个数组格式，可以处理多个平台。
* arch ： 目标架构
* os  ：目标操作系统
不同的platform会生成不同的包。如果platforms为空，则不会进行交叉编译，

而最终生成的包文件名为:
$name-$version-$os-$arch.zip


## 打包
```
goassem package
```

goassem读取 assembly.json 配置对Go项目进行打包。所有的输出在`_out`文件夹中

## 清除
```
goassem clear
```
goassem会删除`_out` 文件夹下的所有内容
