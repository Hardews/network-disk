## network-disk 网盘
### 架构设计

**|-- network-disk**
    **|-- api**
    **|   |-- admin.go**
    **|   |-- router.go**
    **|   |-- upload.go**
    **|   |-- user.go**
    **|-- cmd**
    **|   |-- main.go**
    **|-- dao**
    **|   |-- dao.go**
    **|   |-- upload.go**
    **|   |-- user.go**
    **|-- middleware**
    **|   |-- cors.go**
    **|   |-- jwt.go**
    **|-- model**
    **|   |-- user.go**
    **|-- service**
    **|   |-- admin.go**
    **|   |-- encryption.go**
    **|   |-- upload.go**
    **|   |-- user.go**
    **|-- tool**
    **|   |-- resp.go**
    **|-- uploadFile**
        **|-- breakPoint**
        **|-- jpg**
        **|-- mp3**
        **|-- mp4**
        **|-- png**
        **|-- zip**

### 实现的功能

- [x] **文件上传**

- [x] **文件修改（路径修改、重命名、权限更改）**

- [x] **文件删除**

- [x] **文件下载**

- [x] **文件指纹（通过文件内容进行md5编码获取唯一的文件指纹，不同用户上传多个相同文件只会保存一个）**

- [x] **获取用户的文件信息(所有文件、某个目录的所有文件)**

- [x] **登录注册**

- [x] **文件分享 (可设置分享的时间)**

- [x] **权限管理（所有人可下载、仅获得链接的人可下载、仅自己见）**

- [x] **断点续传**

- [x] **加密分享链接**

- [x] **二维码分享链接**

- [x] **下载限速**

- [ ] **日志管理**

- [x] **管理员模式 (新增用户，查看用户的文件，更改用户的非法文件)**

### 接口说明

####  Login

`POST`  `/login`

`form-data`

`登录`

| 请求参数 | 类型   | 说明           |
| -------- | ------ | -------------- |
| username | string | `必选`，用户名 |
| password | string | `必选`，密码   |

| 返回参数 | 类型   | 说明     |
| -------- | ------ | -------- |
| data     | string | 返回消息 |

| status  | data             | 说明             |
| ------- | ---------------- | ---------------- |
| `false` | "输入的账号为空" | `username`为空   |
| `false` | "输入的密码为空" | `password`为空   |
| `false` | "账号不存在"     | `username`不存在 |
| `false` | "密码错误"       | `password`不正确 |
| `ture`  | token            | token            |



#### Register

`POST`   ` /register`

`form-data`

`注册`

| 请求参数 | 类型   | 说明           |
| -------- | ------ | -------------- |
| username | string | `必选`，用户名 |
| password | string | `必选`，密码   |

| 返回参数 | 类型   | 说明     |
| -------- | ------ | -------- |
| data     | string | 返回消息 |

| status  | data             | 说明             |
| ------- | ---------------- | ---------------- |
| `false` | "输入的账号为空" | `username`为空   |
| `false` | "输入的密码为空" | `password`为空   |
| `false` | "账号已存在"     | `username`已存在 |
| `ture`  | "成功"           |                  |

#### Upload

`POST`   ` /upload`

`form-data`

`上传文件`

| 请求参数      | 类型   | 说明                                                         |
| ------------- | ------ | ------------------------------------------------------------ |
| Authorization | Header | `必选`，鉴权token                                            |
| attribute     | string | `可选`，文件操作权限，不上传默认为公开，Public     = "0" ， Private    =  "1"  ，Permission = "2" |
| category      | string | `必选`,   存储的文件夹名称。主要为 “all”                     |
| path          | string | `必选`,   存储文件的路径                                     |
| file          | file   | `必选`,    用户上传的文件                                    |

| 返回参数 | 类型   | 说明     |
| -------- | ------ | -------- |
| data     | string | 返回消息 |

| status  | data               | 说明           |
| ------- | ------------------ | -------------- |
| `false` | "上传失败，请重试" | 服务器错误     |
| `false` | "服务器错误"       | 服务器错误     |
| `false` | "未指定文件夹"     | `category`为空 |
| `false` | "未指定路径"       | `path`为空     |
| `false` | "文件为空"         | `file`为空     |
| `ture`  | "成功"             |                |

#### Delete

`DELETE`   ` /upload`

`form-data`

`删除文件`

| 请求参数      | 类型   | 说明                                     |
| ------------- | ------ | ---------------------------------------- |
| Authorization | Header | `必选`，鉴权token                        |
| filename      | string | `必选`，文件存储名称                     |
| category      | string | `必选`,   存储的文件夹名称。主要为 “all” |
| path          | string | `必选`,   存储文件的路径                 |

| 返回参数 | 类型   | 说明     |
| -------- | ------ | -------- |
| data     | string | 返回消息 |

| status  | data         | 说明         |
| ------- | ------------ | ------------ |
| `false` | "没有该文件" | 用户无该文件 |
| `ture`  | "成功"       |              |

#### Download

`通过链接下载`  `所有人都可以下载` `仅自己可见` `加密分享`

一般来说这几种链接是通过用户的share返回的，这里只是说怎么处理返回的链接

| 响应头              | value    |
| ------------------- | -------- |
| Content-Disposition | 文件名   |
| Content-Type        | 文件类型 |
| Content-Length      | 文件长度 |

`GET`   ` /download/:filename`

`通过链接下载文件`

| 请求参数      | 类型   | 说明                                                 |
| ------------- | ------ | ---------------------------------------------------- |
| Authorization | Header | `必选`，鉴权token                                    |
| filename      | string | `必选`，用户分享时资源名称经过base64编码得到的字符串 |



`GET`   ` /download/:username/:filename`

`所有人都可下载`

| 请求参数      | 类型   | 说明                                  |
| ------------- | ------ | ------------------------------------- |
| Authorization | Header | `必选`，鉴权token                     |
| :username     | string | `必选`，资源拥有者的用户名            |
| :filename     | string | `必选`，该资源的文件名                |
| path          | query  | `必选`， 该资源所在路径的base64编码值 |
| category      | query  | `必选`， 该资源所在文件夹的名称       |

| 返回参数 | 类型   | 说明     |
| -------- | ------ | -------- |
| data     | string | 返回消息 |

| status  | data                 | 说明     |
| ------- | -------------------- | -------- |
| `false` | "没有该文件"         | 无该文件 |
| `false` | "没有权限下载该文件" | 非法访问 |
| `ture`  |                      |          |



`POST`   ` /download/:username/:filename`

`加密文件下载`

| 请求参数      | 类型   | 说明                                  |
| ------------- | ------ | ------------------------------------- |
| Authorization | Header | `必选`，鉴权token                     |
| :username     | string | `必选`，资源拥有者的用户名            |
| :filename     | string | `必选`，该资源的文件名                |
| path          | query  | `必选`， 该资源所在路径的base64编码值 |
| category      | query  | `必选`， 该资源所在文件夹的名称       |

| 返回参数 | 类型   | 说明     |
| -------- | ------ | -------- |
| data     | string | 返回消息 |

| status  | data                 | 说明     |
| ------- | -------------------- | -------- |
| `false` | "没有该文件"         | 无该文件 |
| `false` | "没有权限下载该文件" | 非法访问 |
| `ture`  |                      |          |



`GET`   ` /user/download/:filename`

`用户下载自己的文件`

| 请求参数      | 类型   | 说明                                  |
| ------------- | ------ | ------------------------------------- |
| Authorization | Header | `必选`，鉴权token                     |
| :filename     | string | `必选`，该资源的文件名                |
| path          | query  | `必选`， 该资源所在路径的base64编码值 |
| category      | query  | `必选`， 该资源所在文件夹的名称       |

| 返回参数 | 类型   | 说明     |
| -------- | ------ | -------- |
| data     | string | 返回消息 |

| status  | data         | 说明     |
| ------- | ------------ | -------- |
| `false` | "没有该文件" | 无该文件 |
| `ture`  |              |          |

#### GetInfo

`获取用户文件信息`

`GET`   ` /user/resource/all`

`获取用户的所有文件`

| 请求参数      | 类型   | 说明              |
| ------------- | ------ | ----------------- |
| Authorization | Header | `必选`，鉴权token |

| 返回参数 | 类型   | 说明     |
| -------- | ------ | -------- |
| data     | string | 返回消息 |

| status | data   | 说明               |
| ------ | ------ | ------------------ |
| `ture` | null   | 该用户没有保存文件 |
| `ture` | 见示例 |                    |

``` 返回示例
{
    "data": [
        {
            "Folder": "all",
            "Path": "我的资源",
            "Filename": "test20220731_114230.png",
            "ResourceName": "./uploadFile/png/a5b22fcc5fd8bbac7aca85e948f83c73.png",
            "Permission": "0",
            "DownloadAddr": "http://127.0.0.1:8080/user/download/test.png?path=5oiR55qE6LWE5rqQ&category=all",
            "CreateAt": "2022-07-31 23:42:30.5198287 +0800 CST m=+33.734813901"
        },
        {
            "Folder": "all",
            "Path": "我的资源",
            "Filename": "test20220731_114228.png",
            "ResourceName": "./uploadFile/png/a5b22fcc5fd8bbac7aca85e948f83c73.png",
            "Permission": "0",
            "DownloadAddr": "http://127.0.0.1:8080/user/download/test.png?path=5oiR55qE6LWE5rqQ&category=all",
            "CreateAt": "2022-07-31 23:42:28.1459814 +0800 CST m=+31.360966601"
        }
    ]
}
```

| 返回参数     | 类型   | 说明               |
| ------------ | ------ | ------------------ |
| Folder       | string | 用户存储的文件夹   |
| Path         | string | 用户存储的路径     |
| Filename     | string | 文件名             |
| ResourceName | string | 在服务器存储的位置 |
| Permission   | string | 权限设置           |
| DownloadAddr | string | 用户自己下载的地址 |
| CreatAt      | string | 文件上传时间       |



`GET`   ` /user/resource`

`获取用户该路径的所有文件`

| 请求参数      | 类型   | 说明                                  |
| ------------- | ------ | ------------------------------------- |
| Authorization | Header | `必选`，鉴权token                     |
| path          | query  | `必选`， 该资源所在路径的base64编码值 |
| category      | query  | `必选`， 该资源所在文件夹的名称       |

| 返回参数 | 类型   | 说明     |
| -------- | ------ | -------- |
| data     | string | 返回消息 |

| status | data   | 说明               |
| ------ | ------ | ------------------ |
| `ture` | null   | 该用户没有保存文件 |
| `ture` | 与上同 |                    |

#### Update

`修改文件`

`PUT`   ` /user/resource`

`form/data`

`修改用户的指定文件`

| 请求参数      | 类型   | 说明                                                         |
| ------------- | ------ | ------------------------------------------------------------ |
| Authorization | Header | `必选`，鉴权token                                            |
| filename      | string | `必选`,    需要修改文件的文件名                              |
| path          | string | `必选`， 该资源所在路径                                      |
| category      | string | `必选`， 该资源所在文件夹的名称                              |
| chose         | string | `必选`,     选择的操作,"1"为修改文件名，"2"为修改文件权限, "3"为修改文件的路径 |
| new           | string | `必选`,      修改的新值                                      |

| 返回参数 | 类型   | 说明     |
| -------- | ------ | -------- |
| data     | string | 返回消息 |

| status  | data         | 说明                             |
| ------- | ------------ | -------------------------------- |
| `false` | "没有该文件" | 用户没有该文件                   |
| `false` | "文件名重复" | 用户修改的文件名与已有文件名重复 |
| `false` | "更新失败"   | 服务器错误                       |
| `ture`  | "成功"       | 修改成功                         |

#### Share

`分享文件`

`正常分享` `二维码分享` `加密分享`



`GET`   ` /user/share/normal/:filename`

`正常分享`

| 请求参数      | 类型   | 说明                                                         |
| ------------- | ------ | ------------------------------------------------------------ |
| Authorization | Header | `必选`，鉴权token                                            |
| :filename     | string | `必选`,    用户想要分享的文件名                              |
| time          | query  | `必选`,    用户分享链接的时效性。 0 表示永久有效,1表示一天，7表示七天，其他值无效. |
| path          | query  | `必选`， 该资源所在路径的base64编码值                        |
| category      | query  | `必选`， 该资源所在文件夹的名称                              |

| 返回参数 | 类型   | 说明     |
| -------- | ------ | -------- |
| data     | string | 返回消息 |

| status  | data                                     | 说明                               |
| ------- | ---------------------------------------- | ---------------------------------- |
| `false` | "链接无效或已过期"                       | 用户输入的链接非法或文件分享已过期 |
| `false` | "分享失败，您以将该文件设置为仅自己可见" | 该用户设置该文件为仅自己可见       |
| `false` | "没有该文件"                             | 该用户没有该文件                   |
| `ture`  | url                                      | 返回一串指向下载地址的url          |



`GET`   ` /user/share/QrCode/:filename`

`二维码分享`

| 请求参数      | 类型   | 说明                                                         |
| ------------- | ------ | ------------------------------------------------------------ |
| Authorization | Header | `必选`，鉴权token                                            |
| :filename     | string | `必选`,    用户想要分享的文件名                              |
| time          | query  | `必选`,    用户分享链接的时效性。 0 表示永久有效,1表示一天，7表示七天，其他值无效. |
| permission    | query  | `可选`， 用户是否修改文件权限                                |
| path          | query  | `必选`， 该资源所在路径的base64编码值                        |
| category      | query  | `必选`， 该资源所在文件夹的名称                              |

| 返回参数 | 类型   | 说明     |
| -------- | ------ | -------- |
| data     | string | 返回消息 |

| status  | data                                     | 说明                                  |
| ------- | ---------------------------------------- | ------------------------------------- |
| `false` | "链接无效或已过期"                       | 用户输入的链接非法或文件分享已过期    |
| `false` | "分享失败，您以将该文件设置为仅自己可见" | 该用户设置该文件为仅自己可见          |
| `false` | "没有该文件"                             | 该用户没有该文件                      |
| `ture`  | image                                    | 返回一张二维码图片,内容为对应下载地址 |



`GET`   ` /user/share/encryption/:filename`

`加密分享`

| 请求参数      | 类型   | 说明                                                         |
| ------------- | ------ | ------------------------------------------------------------ |
| Authorization | Header | `必选`，鉴权token                                            |
| :filename     | string | `必选`,    用户想要分享的文件名                              |
| time          | query  | `必选`,    用户分享链接的时效性。 0 表示永久有效,1表示一天，7表示七天，其他值无效. |
| password      | query  | `可选`， 用户是否指定密码，否则自动生成四位密码并返回        |
| permission    | query  | `可选`， 用户是否修改文件权限                                |
| path          | query  | `必选`， 该资源所在路径的base64编码值                        |
| category      | query  | `必选`， 该资源所在文件夹的名称                              |

| 返回参数 | 类型   | 说明     |
| -------- | ------ | -------- |
| data     | string | 返回消息 |

| status  | data                                     | 说明                               |
| ------- | ---------------------------------------- | ---------------------------------- |
| `false` | "链接无效或已过期"                       | 用户输入的链接非法或文件分享已过期 |
| `false` | "分享失败，您以将该文件设置为仅自己可见" | 该用户设置该文件为仅自己可见       |
| `false` | "没有该文件"                             | 该用户没有该文件                   |
| `ture`  | 见示例                                   |                                    |

``` 返回示例
{
    "password": "QkNQ",
    "path": "http://127.0.0.1:8080/encryption/NTA5YmMyN2I3NWE3OWEzZTAyY2QyOGE5NTJhZWQ0NmFfdGVzdC5wbmc="
}
```

#### Admin

`POST`  `/admin/register`

`form-data`

`管理员注册`

| 请求参数      | 类型   | 说明              |
| ------------- | ------ | ----------------- |
| Authorization | Header | `必选`，鉴权token |
| username      | string | `必选`，用户名    |
| password      | string | `必选`，密码      |

| 返回参数 | 类型   | 说明     |
| -------- | ------ | -------- |
| data     | string | 返回消息 |

| status  | data                   | 说明             |
| ------- | ---------------------- | ---------------- |
| `false` | "非管理员，无权限操作" | 非法访问         |
| `false` | "输入的账号为空"       | `username`为空   |
| `false` | "输入的密码为空"       | `password`为空   |
| `false` | "账号已存在"           | `username`已存在 |
| `ture`  | "成功"                 |                  |

`获取用户文件信息`

`GET`   ` /admin/resource/all`

`获取用户的所有文件`

| 请求参数      | 类型   | 说明                     |
| ------------- | ------ | ------------------------ |
| Authorization | Header | `必选`，鉴权token        |
| username      | query  | `必选`，用户文件的用户名 |

| 返回参数 | 类型   | 说明     |
| -------- | ------ | -------- |
| data     | string | 返回消息 |

| status | data   | 说明               |
| ------ | ------ | ------------------ |
| `ture` | null   | 该用户没有保存文件 |
| `ture` | 见示例 |                    |

``` 返回示例
{
    "data": [
        {
            "Folder": "all",
            "Path": "我的资源",
            "Filename": "test20220731_114230.png",
            "ResourceName": "./uploadFile/png/a5b22fcc5fd8bbac7aca85e948f83c73.png",
            "Permission": "0",
            "DownloadAddr": "http://127.0.0.1:8080/user/download/test.png?path=5oiR55qE6LWE5rqQ&category=all",
            "CreateAt": "2022-07-31 23:42:30.5198287 +0800 CST m=+33.734813901"
        },
        {
            "Folder": "all",
            "Path": "我的资源",
            "Filename": "test20220731_114228.png",
            "ResourceName": "./uploadFile/png/a5b22fcc5fd8bbac7aca85e948f83c73.png",
            "Permission": "0",
            "DownloadAddr": "http://127.0.0.1:8080/user/download/test.png?path=5oiR55qE6LWE5rqQ&category=all",
            "CreateAt": "2022-07-31 23:42:28.1459814 +0800 CST m=+31.360966601"
        }
    ]
}
```

`PUT  `/admin/resource`

`form-data`

`管理员删除非法文件`

| 请求参数      | 类型   | 说明                            |
| ------------- | ------ | ------------------------------- |
| Authorization | Header | `必选`，鉴权token               |
| username      | string | `必选`，资源拥有方的用户名      |
| filename      | string | `必选`,    需要删除文件的文件名 |
| path          | string | `必选`， 该资源所在路径         |
| category      | string | `必选`， 该资源所在文件夹的名称 |
| file          | file   | `必选`， 更换的文件             |

| 返回参数 | 类型   | 说明     |
| -------- | ------ | -------- |
| data     | string | 返回消息 |

| status  | data                   | 说明             |
| ------- | ---------------------- | ---------------- |
| `false` | "非管理员，无权限操作" | 非法访问         |
| `false` | "输入的账号为空"       | `username`为空   |
| `false` | "输入的密码为空"       | `password`为空   |
| `false` | "账号已存在"           | `username`已存在 |
| `ture`  | "成功"                 |                  |

### 一些补充

因为用了文件指纹，当管理员更改一个非法文件时，所有拥有这个非法文件的都会被更改。

文件指纹是百度网盘这样存的，然后结合自己的一些想法用的md5扫描文件内容命名文件。

因为没有和前端真正的合作过，接口的设计方面可能有点问题。

用的是redis存储的用户数据（指文件的类型之类的，不是账号密码），在面对高用户量或者多文件存储可能会有性能问题。

链接的过期用的是redis存储，然后如果是永久有效的，就用mysql存储。

判断链接过期用的是中间件判断，如果不存在链接就会停止这个请求。

用postman测过了，如果按照接口传应该没啥问题？