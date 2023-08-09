Go2ctp
==========================================
## Sample
参考 sample 目录中的代码


## FAQ
> Q: 为什么在有了 [*pseudocodes/goctp*](github.com/pseudocodes/goctp) 项目之后又开启了 *go2ctp* 项目
> 
> A: 发现在封装使用后切换平台 (MacOS), 使用 OpenCTP TTS 系统，以及穿透式认证时(使用 CP 版本) 比较不方便，在静态编译封装的方式下一只能将不同的平台库放在不同的分支下，再加未来对 *ctp-mini* 的支持，就在本项目集中对不同平台的 CTP 库进行集中提供

> Q: 与 [*github.com/hankchan/goctp*](https://github.com/hankchan/goctp) 项目有什么区别
> 
> A: 该项目最早来自 *github.com/qerio/goctp*，是最早可用的完整 golang 版本 ctp 封装，利用 *go build* 对 *swig* 的自动化编译支持，对 CTP 的头文件进行解析并生成中间代码参与编译，一般大家使用的 goctp 80% 都来自该项目的衍生品
> 该项目的问题：
> 1. swig 展开后的中间的中间代码高达十几万行，编译比较慢，约超过 5 分钟
> 2. 早期更新 CTP 库版本时，需要修改部分 CTP 头文件内容，例如行情接口 `SubscribeMarketData` 的合约参数声明方式需要修改，否则 *swig* 在生成接口时无法生成 `__goslice__` 模式的参数，该问题直到最近 swig 才修复？
> 3. 由于(1) 中的原因，在开发时 IDE 的提示无法及时给出，同时 `swig` 生成的接口代码以及使用方式并不友好
> 4. Debug 跟踪困难
> 
> 本项目解决的问题：
> * 简洁封装代码，不采用 swig，不对原始 API 文件进行修改
> * 完全原生 go 语法类型支持，编译速度快，接口提示友好，保留原始文档注释说明
> * Debug 跟踪方便

> Q: 与 [*github.com/haifengat/goctp*](https://github.com/haifengat/goctp) 项目有什么区别
> 
> A: 海风提供的 goctp 的封装对 golang 的特性应用得最好， [*pseudocodes/goctp*](github.com/pseudocodes/goctp) 也参考了该项目的一些细节, 提供 lite 接口的封装
> 该项目的问题：
> * v1 版本的封装比较乱，使用方式与传统的调用方式有出入
> * v2 版本进行了进一步的抽象封装，提供了一些交易基础功能
> * v1 版本一个进程程序只能创建一个 mdapi 和 traderapi 实例，这是该版本最大的问题 
> * v2 版本在 v1 版本的基础上做了增强，改用一个 map 来存储实例指针与对应对象 [quote.go](https://github.com/haifengat/goctp/blob/4f6f9df33fba15d28bf55209fc8111057c9546c7/quote.go#L113C14-L113C45) 可能出于使用频度的考虑，该 map 并未上锁，而 swig 展开代码对使用的 map 都进行了并发锁处理
> * 切换平台库依然不方便，v1 与 v2 版本都是
> 本项目解决问题：
> * 封装层数低，只提供一次 cgo 桥接转换，不涉及数据深拷贝
> * 运行多个 API 实例共存，不同平台的 API 实例可以同时存在一个进程空间中，例如可以同时连接 SIMNOW 以及 OPENCTP
> * 全局指针 map 采用了标准库中的 cgo.Handle 无并发问题

> Q: 为什么不提供 Windows 平台支持
> 
> A: 开发者没有 Windows 环境，不想装虚拟机，大家可以提 PR

> Q: ctp 目录和 ctp_tts目录的封装有什么不同
> 
> A: 
> * ctp 的封装来自 [*pseudocodes/goctp*](github.com/pseudocodes/goctp), 剔除了 *lite* 和 *convert* 相关代码，原则上只提供一层封装，直接桥接 Cgo 空间回调过来的数据
> * ctp 目录属于静态编译链接，其他项目采用该目录下的代码，动态库路径会在编译构建期完成绑定，无需开发者关心，同时提供 MacOS 系统支持
> * ctp_tts 封装方式采用了读取动态库符号的方式，通过在运行时提供动态库的路径，完成 API 实例对象的生成和绑定，理论上 ctp_tts 的封装方式更灵活方便，可以通过配置的方式动态的替换 linux 下不同平台的库，包括官方生产和测评版本，rohan(荣航)，openctp-tts 平台，无需再编译项目代码. ctp-tts 同时支持 openctp 的 MacOS 版本 dylib

## TODO
* ctp-mini 支持
* 更丰富的使用样例
## 免责声明
本项目明确拒绝对产品做任何明示或暗示的担保。由于软件系统开发本身的复杂性，无法保证本项目完全没有错误。如选择使用本项目表示您同意错误和/或遗漏的存在，在任何情况下本项目对于直接、间接、特殊的、偶然的、或间接产生的、使用或无法使用本项目进行交易和投资造成的盈亏、直接或间接引起的赔偿、损失、债务或是任何交易中止均不承担责任和义务。

此声明永久有效