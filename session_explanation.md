# Session 管理详解

## 1. Session 数据结构

```
服务器端存储：
SessionStore = {
    "token_abc123": {                    // Session Token (存在Cookie中)
        "user_id": 1,                    // 用户ID
        "flash": "Welcome!",             // 闪现消息
        "reservation": {...},            // 预订信息
        "error": "Something wrong",      // 错误信息
        "_csrf_token": "xyz789",         // CSRF Token
        "_created_at": "2026-04-10",     // 创建时间
        "_last_access": "2026-04-10"     // 最后访问时间
    },
    "token_def456": {
        "user_id": 2,
        ...
    }
}

客户端Cookie：
session=token_abc123
```

## 2. 核心方法说明

### 2.1 Put(ctx, key, value)
**作用**：在session中存储数据

```go
// 示例
repo.App.SessionManager.Put(r.Context(), "user_id", 1)

// 服务器端变化：
// Before: sessions["token_abc123"] = {}
// After:  sessions["token_abc123"] = {"user_id": 1}
```

### 2.2 Get(ctx, key)
**作用**：从session中读取数据

```go
// 示例
userID := repo.App.SessionManager.Get(r.Context(), "user_id").(int)

// 从 sessions["token_abc123"]["user_id"] 读取值
```

### 2.3 Remove(ctx, key)
**作用**：删除session中的某个键值对

```go
// 示例
repo.App.SessionManager.Remove(r.Context(), "user_id")

// 服务器端变化：
// Before: sessions["token_abc123"] = {"user_id": 1, "flash": "Hi"}
// After:  sessions["token_abc123"] = {"flash": "Hi"}
```

### 2.4 Destroy(ctx)
**作用**：销毁整个session数据，但保留session token

```go
// 示例
repo.App.SessionManager.Destroy(r.Context())

// 服务器端变化：
// Before: sessions["token_abc123"] = {"user_id": 1, "flash": "Hi"}
// After:  sessions["token_abc123"] = {}  // 数据清空，但token还在
```

**注意**：Cookie中的session token仍然存在！

### 2.5 RenewToken(ctx)
**作用**：生成新的session token，保留session数据

```go
// 示例
repo.App.SessionManager.RenewToken(r.Context())

// 服务器端变化：
// Before: 
//   sessions["token_abc123"] = {"user_id": 1}
//   Cookie: session=token_abc123
// 
// After:
//   sessions["token_xyz999"] = {"user_id": 1}  // 新token，数据复制过来
//   sessions["token_abc123"] 被删除
//   Cookie: session=token_xyz999  // Cookie更新为新token
```

**为什么需要RenewToken？**
- 防止Session Fixation攻击
- 在登录/登出等敏感操作时更换token

### 2.6 Exists(ctx, key)
**作用**：检查某个key是否存在

```go
// 示例
if repo.App.SessionManager.Exists(r.Context(), "user_id") {
    // 用户已登录
}

// 检查 sessions["token_abc123"]["user_id"] 是否存在
```

## 3. 典型使用场景

### 场景1：用户登录
```go
func (repo *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
    // 1. 验证用户名密码
    id, _, err := repo.DB.Authenticate(email, password)
    if err != nil {
        // 认证失败
        return
    }
    
    // 2. 更换session token（安全考虑）
    _ = repo.App.SessionManager.RenewToken(r.Context())
    
    // 3. 存储用户ID到session
    repo.App.SessionManager.Put(r.Context(), "user_id", id)
    
    // 4. 存储flash消息
    repo.App.SessionManager.Put(r.Context(), "flash", "logged in successfully")
    
    // Session状态：
    // sessions["new_token"] = {
    //     "user_id": 1,
    //     "flash": "logged in successfully"
    // }
}
```

### 场景2：用户登出（方案A - 完全清除）
```go
func (repo *Repository) Logout(w http.ResponseWriter, r *http.Request) {
    // 销毁所有session数据
    _ = repo.App.SessionManager.Destroy(r.Context())
    
    // Session状态：
    // sessions["token_abc123"] = {}  // 完全清空
    
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
```

### 场景3：用户登出（方案B - 保留flash消息）
```go
func (repo *Repository) Logout(w http.ResponseWriter, r *http.Request) {
    // 1. 更换token（安全）
    _ = repo.App.SessionManager.RenewToken(r.Context())
    
    // 2. 只删除user_id
    repo.App.SessionManager.Remove(r.Context(), "user_id")
    
    // 3. 添加flash消息
    repo.App.SessionManager.Put(r.Context(), "flash", "logged out")
    
    // Session状态：
    // sessions["new_token"] = {
    //     "flash": "logged out"
    // }
    // user_id已被删除，但session还在，可以显示flash消息
    
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
```

### 场景4：检查用户是否登录
```go
func IsAuthenticated(r *http.Request) bool {
    // 检查session中是否有user_id
    exists := app.SessionManager.Exists(r.Context(), "user_id")
    return exists
    
    // 等价于：
    // sessions["token_abc123"]["user_id"] != nil
}
```

## 4. 常见问题

### Q1: Destroy后为什么还能Put数据？
```go
repo.App.SessionManager.Destroy(r.Context())
repo.App.SessionManager.Put(r.Context(), "flash", "logged out")
```

**答案**：
- `Destroy` 清空了session数据，但session token还在
- `Put` 会在同一个token下重新创建数据
- 结果：session被"重置"了，只有新Put的数据

### Q2: RenewToken和Destroy的区别？
```go
// RenewToken: 换token，保留数据
Before: sessions["old_token"] = {"user_id": 1, "flash": "Hi"}
After:  sessions["new_token"] = {"user_id": 1, "flash": "Hi"}

// Destroy: 保留token，清空数据
Before: sessions["token_abc"] = {"user_id": 1, "flash": "Hi"}
After:  sessions["token_abc"] = {}
```

### Q3: 为什么登录时要RenewToken？
**防止Session Fixation攻击**：

```
攻击场景：
1. 攻击者访问网站，获得 session=evil_token
2. 攻击者诱导受害者使用这个token访问网站
3. 受害者用这个token登录
4. 攻击者用同一个token就能访问受害者的账户

防御方法：
登录时调用RenewToken()，生成新token，攻击者的旧token失效
```

### Q4: 为什么登出时要RenewToken？
**防止Session Replay攻击**：

```
攻击场景：
1. 用户登录，session=user_token
2. 攻击者窃取了这个token
3. 用户登出，但token还是user_token
4. 攻击者可能还能用旧token做一些操作

防御方法：
登出时RenewToken()，旧token完全失效
```

## 5. 最佳实践

### 登录流程
```go
func Login(w http.ResponseWriter, r *http.Request) {
    // 1. 验证用户
    id, _, err := repo.DB.Authenticate(email, password)
    if err != nil {
        return
    }
    
    // 2. 更换token（安全）
    _ = repo.App.SessionManager.RenewToken(r.Context())
    
    // 3. 存储用户信息
    repo.App.SessionManager.Put(r.Context(), "user_id", id)
    repo.App.SessionManager.Put(r.Context(), "flash", "Welcome!")
    
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
```

### 登出流程（推荐）
```go
func Logout(w http.ResponseWriter, r *http.Request) {
    // 方案1：完全清除（最安全）
    _ = repo.App.SessionManager.Destroy(r.Context())
    _ = repo.App.SessionManager.RenewToken(r.Context())
    
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// 或

func Logout(w http.ResponseWriter, r *http.Request) {
    // 方案2：保留flash消息
    _ = repo.App.SessionManager.RenewToken(r.Context())
    repo.App.SessionManager.Remove(r.Context(), "user_id")
    repo.App.SessionManager.Put(r.Context(), "flash", "Goodbye!")
    
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
```

## 6. Session生命周期

```
1. 用户首次访问
   └─> 创建空session，生成token，设置Cookie

2. 用户登录
   └─> RenewToken() → Put("user_id", 1)

3. 用户浏览网站
   └─> 每次请求携带Cookie，服务器读取session数据

4. 用户登出
   └─> Destroy() + RenewToken() 或 Remove("user_id")

5. Session过期
   └─> 服务器自动清理过期session（根据配置的过期时间）
```

## 7. Session存储位置

Session数据可以存储在不同地方：

```go
// 1. 内存存储（默认，重启丢失）
sessionManager := scs.New()

// 2. Redis存储（推荐生产环境）
sessionManager.Store = redisstore.New(redisPool)

// 3. 数据库存储
sessionManager.Store = postgresstore.New(db)

// 4. Cookie存储（不推荐，数据在客户端）
sessionManager.Store = cookiestore.New()
```

## 8. 总结

**核心概念**：
- Session = Token (客户端Cookie) + Data (服务器存储)
- Token是钥匙，Data是保险箱里的内容

**关键方法**：
- `Put/Get/Remove`: 操作session数据
- `Destroy`: 清空数据，保留token
- `RenewToken`: 换token，保留数据
- `Exists`: 检查key是否存在

**安全原则**：
- 登录时：RenewToken（防止fixation）
- 登出时：Destroy + RenewToken（防止replay）
- 敏感操作：定期RenewToken
