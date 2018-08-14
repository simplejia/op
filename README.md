# 通用op（后台运营）系统

通用op（后台运营）系统，后台开发工程师打造的针对后台服务的通用配置管理系统，带web操作界面

做这套系统的缘由：
> 由于专业web前端开发人员支持不足，后台服务，没有web管理界面可用，管理后台服务很不方便，尤其是微服务化后，后台大大小小服务特别多，急需一个简单可视化的界面系统统一管理和利用这些服务

这套系统的最初设想：
1. 界面足够简单，html，js，css足够简单，因为后续维护开发都是由后台开发人员执行（非专业web前端）
2. 界面简单，但功能完善，支持列表，支持翻页，支持数据更新，删除等操作，支持外部接口调用，支持输入参数校验，支持输入参数按单选或者复选的方式辅助输入，支持按输入参数过滤数据等
3. 支持基本的后台服务上线，回滚，重启等操作，支持文件上传、下载
4. 集中式管理，支持区分账号的数据保存，比如支持操作历史保存和复用，支持基于账号的权限控制等

这套系统目前具有的功能：
> 简单说，我对这套系统最初的设想，目前已全部实现，早在几年前，我的心里就萌生了要做成这套系统的想法，我做后台服务开发至今有10多年了，一直就想有这样的一套系统来管理项目里的各种后台服务，这个服务一定要满足：通用，实用，简单，现在这个最初的想法终于实现了。

TODO:
1. 支持更细粒度的权限控制，比如某些配置项只能由创建者和管理员操作等，理想中这一功能需要一个通用的后台账号服务支持
2. 接口类型是list的接口支持对返回的参数做定制化配置，比如返回错误码定义，返回列表字段名等，也就是格式支持自定义，而不仅仅是目前内定好的格式，目前格式支持如下：json格式，需满足：{"ret": 1, "data": {"list": [...]}}，ret存储返回码，1代表成功，其它为失败，data.list存储数组结构列表数据
3. 调用后端服务不仅仅限post方式，也不仅仅限json格式的body，比如支持form表单，get，put等
4. 支持针对输入输出字段配置描述信息
5. 支持接口类型是update, delete的接口的操作历史保存
6. 为方便审计，支持操作历史记录永久存储备份，而不像现在针对每个人，每个接口，最多只保存一定数目的记录
7. 支持以现有配置为模板复制生成新的配置项

这套系统现在长什么样？
1. 有账号认证，如下图：

![1](https://github.com/simplejia/nothing/raw/master/1.tiff)

希望账号管理、权限认证系统是一个统一的，独立的后台账号管理服务，目前开源出去的是把这一步省掉了，大家可以根据自己的情况自行添加处理代码（filter/auth.go）

2. 首页，如下图：

![2](https://github.com/simplejia/nothing/raw/master/2.tiff)

列表里展示的是已配置的可用服务

3. 点击新建，进入服务配置页面，如下图：

![3](https://github.com/simplejia/nothing/raw/master/3.tiff)
￼

4. 查看某一个配置项

![4](https://github.com/simplejia/nothing/raw/master/4.tiff)
￼

5. 进入某一项（默认进入类型是list的action，如果没有配，就会进入类型是customer的列表）

![5](https://github.com/simplejia/nothing/raw/master/5.tiff)
￼

如上，是因为配置了cid这个必填字段，limit和offset是可选字段，其中limit配置的默认值是20，执行后，显示如下：

![6](https://github.com/simplejia/nothing/raw/master/6.tiff)

￼
以上是列表信息，一共返回1条数据，继续执行会把cid, limit, offset, total字段发给服务端作为输入参数

![7](https://github.com/simplejia/nothing/raw/master/7.tiff)

￼
以上是没有配置list类型的action，进入某一项会返回类型是customer的列表（或者点击list类型的页面的“其它”进入）

![8](https://github.com/simplejia/nothing/raw/master/8.tiff)

￼
以上是点击“更新”或“删除”进入的页面，每一项数据都可以修改，点击执行后会把以上数据post给服务端（具体执行接口是来自于update或delete类型的action配置）

![9](https://github.com/simplejia/nothing/raw/master/9.tiff)

￼
以上是点击某一customer类型的action进入的页面，上面一部分是传给后端的输入参数，可以删除一些字段，不用传，下面一部分是执行过的历史，可以恢复记录重新执行。

安装使用：
1. 下载源代码
> go get github.com/simplejia/op
2. 配置数据库
> 目前的配置信息存储在mongo db，需要修改配置文件：mongo/op.json
3. 使用
> 进入op目录，启动编译好的op程序，比如：./op -env dev，打开浏览器，输入网址，如果是本地测试运行，请输入：127.0.0.1:8336


最佳实践：
1. 基本功能演示（增删改查）

这是一个推送服务，以下是配置项：

![10](https://github.com/simplejia/nothing/raw/master/10.tiff)

![11](https://github.com/simplejia/nothing/raw/master/11.tiff)
￼
￼

这是配置项的进入页（列表页）：

![12](https://github.com/simplejia/nothing/raw/master/12.tiff)

![13](https://github.com/simplejia/nothing/raw/master/13.tiff)
￼

￼
这是点击列表页的更新/删除：

![14](https://github.com/simplejia/nothing/raw/master/14.tiff)
￼

2. 服务上线

这是一个php的上线功能，以下是配置项：

![15](https://github.com/simplejia/nothing/raw/master/15.tiff)
￼

这是配置项的进入页：

![16](https://github.com/simplejia/nothing/raw/master/16.tiff)
￼

点击/online/trans_cmd:

![17](https://github.com/simplejia/nothing/raw/master/17.tiff)
￼

点击/online/trans_file:

![18](https://github.com/simplejia/nothing/raw/master/18.tiff)
￼

这是一个go服务的上线功能，以下是配置项：

![19](https://github.com/simplejia/nothing/raw/master/19.tiff)
￼

这是配置项的进入页：

![20](https://github.com/simplejia/nothing/raw/master/20.tiff)

￼
点击/online/trans_cmd:

![21](https://github.com/simplejia/nothing/raw/master/21.tiff)

￼
点击/online/trans_file:

![22](https://github.com/simplejia/nothing/raw/master/22.tiff)
￼

注：上线服务依赖：github.com/simplejia/online
online项目用于提供远程文件上传及远程执行命令功能，类似运维工具：ansible

3. 下载数据

这是一个提供数据下载的配置项：

![23](https://github.com/simplejia/nothing/raw/master/23.tiff)

￼
注意下类型：transparent，这个类型表示后端接口返回什么数据，页面上直接展示，不做任何处理

以下是配置项的进入页，点击执行后的效果：（提示有文件正在下载）

![24](https://github.com/simplejia/nothing/raw/master/24.tiff)
￼

4. 复杂功能演示

这是一个提供视频处理功能的配置项：（提供按条件过滤的功能）

![25](https://github.com/simplejia/nothing/raw/master/25.tiff)
￼

这是配置项的进入页：

![27](https://github.com/simplejia/nothing/raw/master/27.tiff)
￼

注意“remote_ip”这个字段，配置“数据源”是“从URL”，此ip列表是调用配置的url接口返回的结果

点击“执行”后的运行结果：（部分结果如下）

![28](https://github.com/simplejia/nothing/raw/master/28.tiff)
￼


这是一个提供视频处理功能的配置项：（提供更好的报表展示功能）

![26](https://github.com/simplejia/nothing/raw/master/26.tiff)
￼

这是配置项的进入页：

![29](https://github.com/simplejia/nothing/raw/master/29.tiff)
￼

点击“执行”后的运行结果：（部分结果如下）

![30](https://github.com/simplejia/nothing/raw/master/30.tiff)
￼

注：此表格是调用的后端接口直接吐出来的html代码


## 依赖
    wsp: github.com/simplejia/wsp
    clog: github.com/simplejia/clog
    utils: github.com/simplejia/utils
    namecli: github.com/simplejia/namecli
    mongo: gopkg.in/mgo.v2

## 注意
    如果在controller里修改了路由，编译前需执行go generate，实际是运行了wsp这个工具，所以需要提前go get github.com/simplejia/wsp
