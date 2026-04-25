# Composable 规范

## 何时抽取 Composable

当页面 `.vue` 文件超过 **300 行**，或者某段逻辑在 **2 个以上页面复用**时，抽取为 composable。

**不要提前抽象**：不要因为"看起来可以复用"就抽，等到真的需要复用时再抽。

---

## 命名规则

- 文件名：`use` + 业务描述，例如 `useOrderList.js`、`usePagination.js`
- 函数名与文件名一致：`export function useOrderList()`
- 返回值用解构方式暴露，不要返回一个大对象

---

## 页面专属 Composable

当一个页面逻辑复杂，但该逻辑不需要复用时，放在页面目录下：

```
pages/order/
├── index.vue
└── composables/
    └── useOrderListLogic.js  # 只给 index.vue 用
```

```js
// pages/order/composables/useOrderListLogic.js
import { ref, onMounted } from 'vue'
import { getOrderList } from '@/api/order'

export function useOrderListLogic() {
  const orderList = ref([])
  const isLoading = ref(false)
  const currentTab = ref('all')  // all | pending | done

  const tabOptions = [
    { label: '全部', value: 'all' },
    { label: '待付款', value: 'pending' },
    { label: '已完成', value: 'done' }
  ]

  async function loadOrderList() {
    isLoading.value = true
    const data = await getOrderList({ status: currentTab.value })
    orderList.value = data.list
    isLoading.value = false
  }

  function handleTabChange(tab) {
    currentTab.value = tab
    loadOrderList()
  }

  onMounted(() => {
    loadOrderList()
  })

  return {
    orderList,
    isLoading,
    currentTab,
    tabOptions,
    handleTabChange,
    loadOrderList
  }
}
```

```vue
<!-- pages/order/index.vue —— 清爽，职责清晰 -->
<script setup>
import { useOrderListLogic } from './composables/useOrderListLogic'
const { orderList, isLoading, currentTab, tabOptions, handleTabChange } = useOrderListLogic()
</script>
```

---

## 全局公共 Composable

放在 `src/composables/` 下，供多个页面使用：

```js
// src/composables/usePagination.js
import { ref } from 'vue'

export function usePagination(fetchFn, options = {}) {
  const pageSize = options.pageSize || 10

  const list = ref([])
  const currentPage = ref(1)
  const isLoading = ref(false)
  const hasMore = ref(true)

  async function loadFirstPage() {
    currentPage.value = 1
    list.value = []
    hasMore.value = true
    await loadNextPage()
  }

  async function loadNextPage() {
    if (isLoading.value || !hasMore.value) return

    isLoading.value = true
    const data = await fetchFn({ page: currentPage.value, pageSize })
    list.value = [...list.value, ...data.list]
    hasMore.value = data.list.length === pageSize
    currentPage.value += 1
    isLoading.value = false
  }

  return {
    list,
    isLoading,
    hasMore,
    loadFirstPage,
    loadNextPage
  }
}
```

```js
// src/composables/useLoading.js —— 管理加载状态的简单封装
import { ref } from 'vue'

export function useLoading() {
  const isLoading = ref(false)

  async function withLoading(asyncFn) {
    isLoading.value = true
    try {
      return await asyncFn()
    } finally {
      isLoading.value = false
    }
  }

  return { isLoading, withLoading }
}
```

---

## 禁止事项

```js
// ❌ 把 UI 逻辑（toast、navigate）放进 composable
export function useOrderLogic() {
  async function submitOrder(data) {
    await createOrder(data)
    uni.showToast({ title: '下单成功' })   // ❌ UI 逻辑不属于 composable
    uni.navigateBack()                      // ❌
  }
}

// ✅ composable 只返回数据和方法，UI 由调用方处理
export function useOrderLogic() {
  async function submitOrder(data) {
    await createOrder(data)
    // 只做数据操作，不做 UI
  }
  return { submitOrder }
}

// 页面里处理 UI
async function handleSubmit() {
  await submitOrder(formData.value)
  uni.showToast({ title: '下单成功' })
  uni.navigateBack()
}
```

```js
// ❌ 一个 composable 做太多事
export function useEverything() {
  // 用户信息 + 订单列表 + 购物车 + 地址管理（不要这样）
}

// ✅ 每个 composable 职责单一
export function useUserInfo() { ... }
export function useOrderList() { ... }
```