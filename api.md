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

详细接口信息如下：


## 关联数据

### 一、查询关联数据(根据领域查询所有)
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
| scope | true | string | 领域中文名称 |
| order | false | string | "asc" 或 "desc"，默认为 asc |
| order_by | false | string | 排序依据，选项："code"，默认为 code |

返回字段
---
| 字段 | 必选 | 类型 | 说明 |
| ----- | ---- | ---- | ---- |
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
