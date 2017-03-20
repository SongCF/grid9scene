# 景和场景服务（九宫格）API

**API 开发作者 宋长福**

**API [Zookeeper] 注册信息：http://139.198.2.84:8011/get_service?path=/platform/scene/http**


# Group API接口

## 应用操作 [/api/v1/app/{APP_ID}?user_id={USER_ID}&token={TOKEN}]
应用相关操作，通过该接口可以创建或删除应用，`APP_ID`是要操作的应用标识，用户需要传入`USER_ID`和`TOKEN`做权限验证。

+ Parameters

    + APP_ID: app101 (string) - 应用唯一标识.

    + USER_ID: 1001 (number) - 用户唯一标识，创建者.

    + TOKEN: af618f49ae7c318ec2a (string) - 用户校验码，用于验证用户是否正确且有权限.

### 创建应用 [POST]

+ Request 

    + Headers

            Accept: application/json

+ Response 200 

    + Headers

            Content-Type: application/json

    + Body

			{
			  "code": 10000000,
			  "msg": "Success"
			}

### 删除应用 [DELETE]

+ Request 

    + Headers

            Accept: application/json

+ Response 200 

    + Headers

            Content-Type: application/json

    + Body

			{
			  "code": 10000000,
			  "msg": "Success"
			}

## 场景操作 [/api/v1/app/{APP_ID}/space/{SPACE_ID}?user_id={USER_ID}&token={TOKEN}&grid_width={GRID_WIDTH}&grid_height={GRID_HEIGHT}]

### 创建场景 [POST]

+ Parameters

    + APP_ID: app101 (string) - 应用唯一标识.

    + SPACE_ID: city0280 (string) - 应用唯一标识.

    + USER_ID: 1001 (number) - 用户唯一标识，创建者.

    + TOKEN: af618f49ae7c318ec2a (string) - 用户校验码，用于验证用户是否正确且有权限.

    + GRID_WIDTH: 100 (number) - 场景进行九宫格拆分的宽.

    + GRID_HEIGHT: 100 (number) - 场景进行九宫格拆分的宽.

+ Request 

    + Headers

            Accept: application/json

+ Response 200 

    + Headers

            Content-Type: application/json

    + Body

			{
			  "code": 10000000,
			  "msg": "Success"
			}

## 场景操作 [/api/v1/app/{APP_ID}/space/{SPACE_ID}?user_id={USER_ID}&token={TOKEN}]

### 删除场景 [DELETE]

+ Parameters

    + APP_ID: app101 (string) - 应用唯一标识.

    + SPACE_ID: city0280 (string) - 场景唯一标识.

    + USER_ID: 1001 (number) - 用户唯一标识.

    + TOKEN: af618f49ae7c318ec2a (string) - 用户校验码，用于验证用户是否正确且有权限.

+ Request 

    + Headers

            Accept: application/json

+ Response 200 

    + Headers

            Content-Type: application/json

    + Body

			{
			  "code": 10000000,
			  "msg": "Success"
			}

## 查询用户位置 [/api/v1/app/{APP_ID}/user/{TARGET_USER_ID}/pos?user_id={USER_ID}&token={TOKEN}]

### 查询用户位置 [GET]

+ Parameters

    + APP_ID: app101 (string) - 应用唯一标识.

    + TARGET_USER_ID: 2002 (number) - 被查询者.

    + USER_ID: 1001 (number) - 查询者.

    + TOKEN: af618f49ae7c318ec2a (string) - 用户校验码，用于验证用户是否正确且有权限.

+ Request 

    + Headers

            Accept: application/json

+ Response 200 

    + Headers

            Content-Type: application/json

    + Body

			{
			  "code": 10000000,
			  "msg": "Success",
			  "response": {
				"space_id": "city_0280",
				"pos_x": 20.2,
				"pos_y": 120.9,
				"angle": 360
			  }
			}


# Group 错误码

## 错误码定义

|名字|说明|
|:---|:------|
|10000000|请求成功，正常返回|
|10070001|服务器错误，返回服务器错误堆栈|
|10070002|request timeout|
|10070101|has not login|
|10070102|duplicate login|
|10070103|user not found|
|10070104|error user|
|10070201|app not exist|
|10070301|space not exist|
|10070302|space already exist|
|10070401|error message format|
|10070402|unknown message|
|10070403|Missing parameter: xxx|
|10070404|cmd not support|

## 错误JSON回复示例

``` js
{
  "code": 10070403,
  "msg": "Missing parameter: user_id"
}
```
