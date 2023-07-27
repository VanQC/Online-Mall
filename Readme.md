# CloudRestaurant 项目

本项目主要实现了用户的注册和登录功能，具体包括：

+ 使用 手机号+密码+短信验证码 进行用户注册（使用的是阿里云短信验证码服务api）
+ 实现 手机号+短信验证码 登录功能
+ 实现 手机号+密码登录+图片验证码 登录功能
+ 提供用户头像上传、用户信息查询功能
+ 实现跨域访问中间件

+ 实现 获取食品类别、批量添加食品信息 api

# v1.0.1
用户头像上传、用户信息查询功能需要进行用户登录状态认证，在第一版中：
用户头像上传 使用的是 session 进行登录状态认证
用户信息查询 使用的是 cookie 进行登录状态认证
（目的是为了进行 session 和 cookie 的简单学习，在第二版中将全部使用 jwt 中间件进行身份认证）