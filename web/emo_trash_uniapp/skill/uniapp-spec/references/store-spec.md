# Store（Pinia）规范

## 原则

- 每个业务模块一个 store 文件，禁止所有状态集中在一个 store
- Store 只存**跨页面共享**的状态，页面私有状态用 `ref/reactive`
- Action 直接写业务逻辑，不要再包一层 service

---

## 标准 Store 结构

```js
// src/stores/user.js
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, getUserProfile } from '@/api/user'

export const useUserStore = defineStore('user', () => {
  // State（响应式数据）
  const token = ref(uni.getStorageSync('token') || '')
  const userInfo = ref(null)

  // Getters（计算属性）
  const isLoggedIn = computed(() => !!token.value)
  const nickname = computed(() => userInfo.value?.nickname || '未登录')

  // Actions（方法）
  async function loginWithPhone(phone, code) {
    const data = await login({ phone, code })
    token.value = data.token
    userInfo.value = data.userInfo
    uni.setStorageSync('token', data.token)
  }

  async function fetchUserInfo() {
    if (!token.value) return
    const data = await getUserProfile()
    userInfo.value = data
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    uni.removeStorageSync('token')
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    nickname,
    loginWithPhone,
    fetchUserInfo,
    logout
  }
})
```

---

## 各模块 Store 职责

```js
// src/stores/cart.js —— 购物车状态
export const useCartStore = defineStore('cart', () => {
  const items = ref([])
  const totalCount = computed(() => items.value.reduce((sum, item) => sum + item.quantity, 0))
  const totalPrice = computed(() => items.value.reduce((sum, item) => sum + item.price * item.quantity, 0))

  function addItem(product, quantity = 1) {
    const existing = items.value.find(item => item.id === product.id)
    if (existing) {
      existing.quantity += quantity
    } else {
      items.value.push({ ...product, quantity })
    }
  }

  function removeItem(productId) {
    items.value = items.value.filter(item => item.id !== productId)
  }

  function clearCart() {
    items.value = []
  }

  return { items, totalCount, totalPrice, addItem, removeItem, clearCart }
})
```

```js
// src/stores/app.js —— 全局应用状态（系统信息、全局配置）
export const useAppStore = defineStore('app', () => {
  const systemInfo = ref(null)
  const networkType = ref('unknown')

  function initSystemInfo() {
    systemInfo.value = uni.getSystemInfoSync()
  }

  return { systemInfo, networkType, initSystemInfo }
})
```

---

## 什么状态不该放 Store

```js
// ❌ 页面私有的加载状态，不应该放 store
// stores/order.js
const isOrderListLoading = ref(false)  // 只有 order/index.vue 用，放页面里

// ✅ 跨页面共享的状态，才放 store
// stores/order.js
const unreadOrderCount = ref(0)  // tabbar badge 用，多处使用
```

**判断依据**：这个状态是否需要在 2 个以上的页面/组件里访问？不是 → 放页面，是 → 放 Store。

---

## 禁止事项

```js
// ❌ 把所有状态写进一个 store
// stores/index.js （不要这样做）
export const useStore = defineStore('main', () => {
  const userInfo = ref(null)
  const cartItems = ref([])
  const orderList = ref([])
  const productList = ref([])
  // ... 100 个状态
})

// ✅ 按业务模块拆分
// stores/user.js、stores/cart.js、stores/order.js
```

```js
// ❌ 在 store 里做 UI 操作
async function loadOrders() {
  uni.showLoading({ title: '加载中' })  // ❌ UI 逻辑不属于 store
  const data = await getOrderList()
  orderList.value = data
  uni.hideLoading()  // ❌
}

// ✅ UI 逻辑在调用方（页面）处理
async function loadOrders() {
  const data = await getOrderList()
  orderList.value = data
}
```