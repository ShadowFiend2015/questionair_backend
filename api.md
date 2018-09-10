# questionair_backend 接口文档

## API 域名

请求API示例：
```
http://{IP}:{PORT}/api/scope/check
```

## 格式说明
返回结果格式
---
| 字段 | 类型 | 说明 |
| ------ | ------ | ------- |
| code   | int    | 错误码   |
| msg    | string | 错误消息 |
| body   | struct | 具体数据 |

统一时间格式：`2000-01-01 12:00:00`, 且秒数为`00`
统一日期格式：`2000-01-01`

##	错误码
code | 说明
---- | ---
90001 | 内部错误
90002 | 接口不存在
90003 | 未鉴权
90004 | 鉴权失败
90005 | 请求参数错误
90006 | 没有操作权限
90007 | 超过限制
90008 | 重复操作
90009 | 登录失败
91001 | 数据读取错误
91002 | 数据插入错误
91003 | 数据修改错误
91004 | 数据删除错误
91005 | 数据修改无效
91006 | 无此数据


详细接口信息如下：


## 用户

### 一、用户登录
---
/user/login

支持格式
---
JSON

HTTP Method
---
POST

请求参数
---
| 字段 | 必选 | 类型 | 说明 |
| ----- | ---- | ---- | ---- |
| account | true | string | 用户帐号 |
| password | true | string | 用户密码 |

返回字段
---
| 字段 | 必返 | 类型 | 说明 |
| ----- | ---- | ---- | ---- |
| pass | true | bool | 是否成功，true/false |
| token | true | string | 登录成功返回的用户 token |

返回结果示例
----
```
{
    "code": 0,
    "msg": "success",
    "body": {
        "pass": true,
        "token": "test_user_token"
    }
}
```
*注意*
> `code` 与 `msg` 字段只代表请求方式的正确与否，是否成功登录要看 `body` 里的 `pass` 字段。


## 建筑类型

### 一、查询(所有)
---
/api/scope/read

支持格式
---
QueryString

HTTP Method
---
GET

请求参数
---
| 字段 | 必选 | 类型 | 说明 |
| ----- | ---- | ---- | ---- |
| token | true | string | 用户登陆时产生的 token，放在 header 里 |

返回字段
---
| 字段 | 必选 | 类型 | 说明 |
| ----- | ---- | ---- | ---- |
| id | true | int | 类型 id |
| name | true | string | 类型名称 |
| code | true | string | 类型编码 |

返回结果示例
----
```
{
    "code": 0,
    "msg": "success",
    "body": {
        "total": 9,
        "data": [
            {
                "id": 1,
                "name": "电网工程",
                "code": "DW"
            },
            {
                "id": 2,
                "name": "建筑工程",
                "code": "JZ"
            },
            {
                "id": 3,
                "name": "铁路工程",
                "code": "TL"
            },
            {
                "id": 4,
                "name": "公路工程",
                "code": "GL"
            },
            {
                "id": 5,
                "name": "水利工程",
                "code": "SL"
            },
            {
                "id": 6,
                "name": "民航工程",
                "code": "MH"
            },
            {
                "id": 7,
                "name": "石油管道",
                "code": "YQ"
            },
            {
                "id": 8,
                "name": "城市轨道",
                "code": "GD"
            },
            {
                "id": 9,
                "name": "机械制造",
                "code": "JX"
            }
        ]
    }
}
```

### 二、查询(所有除了给定的类型)
---
/api/scope/other/read

支持格式
---
QueryString

HTTP Method
---
GET

请求参数
---
| 字段 | 必选 | 类型 | 说明 |
| ----- | ---- | ---- | ---- |
| token | true | string | 用户登陆时产生的 token，放在 header 里 |
| scope_name | true | string | 用户登陆时产生的 token，放在 header 里 |

返回字段
---
| 字段 | 必选 | 类型 | 说明 |
| ----- | ---- | ---- | ---- |
| id | true | int | 类型 id |
| name | true | string | 类型名称 |
| code | true | string | 类型编码 |

返回结果示例
----
同 `一、查询`


## 关联数据

### 一、创建关联数据
---
/api/link/create

支持格式
---
QueryString

HTTP Method
---
POST

请求参数
---
| 字段 | 必选 | 类型 | 说明 |
| ----- | ---- | ---- | ---- |
| token | true | string | 用户登陆时产生的 token，放在 header 里 |
| scope_name_1 | true | string | 领域中文名称1 |
| scope_name_2 | true | string | 领域中文名称2 |
| element_name_1 | true | string | 构件中文名称1 |
| element_name_2 | true | string | 构件中文名称2 |
| host_scope | true | string | 主体领域中文名称 |

返回字段
---
| 字段 | 必选 | 类型 | 说明 |
| ----- | ---- | ---- | ---- |
| id | true | int | 关联数据 id |
| code | true | string | 构建编码 |
| name | true | string | 构建名称 |
| link_code | true | string | 相似构建的编码 |
| link_name | true | string | 相似构建的名称 |
| link_full_name | true | string | 相似构建的名称(包含领略，以":"隔开) |
| status | true | int | 确认状态，0 - 未确认，1 - 已确认 |

返回结果示例
----
```

```

### 二、查询关联数据(根据领域查询所有)
---
/api/link/scope/read

支持格式
---
QueryString

HTTP Method
---
GET

请求参数
---
| 字段 | 必选 | 类型 | 说明 |
| ----- | ---- | ---- | ---- |
| token | true | string | 用户登陆时产生的 token，放在 header 里 |
| scope_name | true | string | 领域中文名称 |
| order | false | string | "asc" 或 "desc"，默认为 asc |
| order_by | false | string | 排序依据，选项："code"，默认为 code |

返回字段
---
| 字段 | 必选 | 类型 | 说明 |
| ----- | ---- | ---- | ---- |
| id | true | int | 关联数据 id |
| code | true | string | 构建编码 |
| name | true | string | 构建名称 |
| link_code | true | string | 相似构建的编码 |
| link_name | true | string | 相似构建的名称 |
| link_full_name | true | string | 相似构建的名称(包含领略，以":"隔开) |
| status | true | int | 确认状态，0 - 未确认，1 - 已确认 |

返回结果示例
----
```

```

### 三、确认关联数据
---
/api/link/confirm

支持格式
---
QueryString

HTTP Method
---
POST

请求参数
---
| 字段 | 必选 | 类型 | 说明 |
| ----- | ---- | ---- | ---- |
| token | true | string | 用户登陆时产生的 token，放在 header 里 |
| id | true | int | 关联数据 id |
| host_scope | true | string | 主体领域中文名称 |
| agree | false | bool | 确认关联数据是否正确，true - 正确，false - 错误，默认为错误 |

返回字段
---
无

返回结果示例
----
```

```
