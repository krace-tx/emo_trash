# API 请求规范

## 基础封装（仅此一层，禁止再套一层）

```js
// src/api/request.js
import { useUserStore } from '@/stores/user'

const BASE_URL = import.meta.env.VITE_API_BASE_URL

function request(options) {
  return new Promise((resolve, reject) => {
    const userStore = useUserStore()

    uni.request({
      url: BASE_URL + options.url,
      method: options.method || 'GET',
      data: options.data,
      header: {
        'Content-Type': 'application/json',
        Authorization: userStore.token ? `Bearer ${userStore.token}` : ''
      },
      success(res) {
        if (res.statusCode === 401) {
          userStore.logout()
          uni.navigateTo({ url: '/pages/login/index' })
          reject(new Error('未登录'))
          return
        }
        if (res.statusCode !== 200) {
          uni.showToast({ title: res.data?.message || '请求失败', icon: 'error' })
          reject(new Error(res.data?.message || '请求失败'))
          return
        }
        if (res.data.code !== 0) {
          uni.showToast({ title: res.data.message || '操作失败', icon: 'error' })
          reject(new Error(res.data.message))
          return
        }
        resolve(res.data.data)
      },
      fail(err) {
        uni.showToast({ title: '网络异常，请重试', icon: 'error' })
        reject(err)
      }
    })
  })
}

export default request
```

---

## 业务 API 文件（按模块分文件）

```js
// src/api/user.js
import request from './request'

// 获取用户信息
export function getUserProfile(userId) {
  return request({
    url: `/user/${userId}`,
    method: 'GET'
  })
}

// 更新用户信息
export function updateUserProfile(data) {
  return request({
    url: '/user/profile',
    method: 'PUT',
    data
  })
}

// 用户登录
export function login(params) {
  return request({
    url: '/auth/login',
    method: 'POST',
    data: params
  })
}
```

```js
// src/api/order.js
import request from './request'

export function getOrderList(params) {
  return request({
    url: '/orders',
    method: 'GET',
    data: params
  })
}

export function getOrderDetail(orderId) {
  return request({
    url: `/orders/${orderId}`,
    method: 'GET'
  })
}

export function createOrder(data) {
  return request({
    url: '/orders',
    method: 'POST',
    data
  })
}

export function cancelOrder(orderId) {
  return request({
    url: `/orders/${orderId}/cancel`,
    method: 'POST'
  })
}
```

---

## API 调用规范

### 在页面中调用（直接调用，不要再封装一层）

```vue
<script setup>
import { getUserProfile, updateUserProfile } from '@/api/user'

// ✅ 直接调用，清晰
async function loadProfile() {
  const data = await getUserProfile(userId)
  userInfo.value = data
}

// ❌ 不要再包一层
function fetchUser() {  // 这一层没有任何意义
  return getUserProfile(userId)
}
</script>
```

### 错误处理规则

- **通用错误**（网络异常、服务器错误）：在 `request.js` 统一处理，页面无需关心
- **业务错误**（需要特殊处理）：在调用处用 `try/catch` 单独处理
- **不要** 在每个 API 函数里重复写错误提示

```js
// ✅ 需要特殊处理时才 try/catch
async function handleSubmit() {
  try {
    await createOrder(orderData)
    uni.showToast({ title: '下单成功' })
    uni.navigateBack()
  } catch (error) {
    // request.js 已经 showToast 了，这里只做特殊处理
    // 比如重置表单状态
    isSubmitting.value = false
  }
}

// ✅ 普通加载，不需要 try/catch（request.js 已处理）
async function loadData() {
  isLoading.value = true
  const data = await getOrderList({ page: currentPage.value })
  orderList.value = data.list
  isLoading.value = false
}
```

---

## 文件命名规则

| 场景 | 文件名 |
|------|--------|
| 用户相关 | `api/user.js` |
| 订单相关 | `api/order.js` |
| 商品相关 | `api/product.js` |
| 购物车 | `api/cart.js` |
| 首页数据 | `api/home.js` |
| 文件上传 | `api/upload.js` |

**原则**：文件名 = 业务模块名，不要叫 `apiService.js`、`httpHelper.js` 这种无意义的名字。