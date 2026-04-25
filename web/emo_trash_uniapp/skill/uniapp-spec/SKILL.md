---
name: uniapp-spec
description: UniApp 框架项目规范，专为 AI 辅助开发设计。当用户在 UniApp 项目中让 AI 编写代码、生成页面、创建组件、设计架构、或者询问如何组织 UniApp 项目时，必须触发此 skill。也适用于："帮我写一个 UniApp 页面"、"创建一个组件"、"怎么组织项目结构"、"写接口请求"、"写 store"等场景。AI 必须严格遵守本规范，防止将所有代码堆入单一文件、过度封装、以及创造难以阅读的抽象层。
---

# UniApp AI 开发规范

> 本规范专为 **AI 辅助开发** 设计，核心目标：代码可读、文件职责单一、禁止过度封装。

---

## 核心原则（AI 必须遵守）

### 原则一：一个文件只做一件事

❌ 禁止行为：
- 把多个页面逻辑写进同一个 `.vue` 文件
- 把所有 API 请求写进一个 `api.js`
- 把不相关的工具函数堆在同一个 `utils.js`
- 把所有 store 状态写进一个 `store/index.js`

✅ 正确做法：
- 每个页面独立一个目录，页面逻辑 > 200 行时拆分 composable
- API 按业务模块分文件：`api/user.js`、`api/order.js`
- Store 按业务模块分 store：`stores/user.js`、`stores/cart.js`
- 工具函数按用途分文件：`utils/format.js`、`utils/validate.js`

---

### 原则二：禁止过度封装

**判断标准**：一个封装是否必要，看它有没有减少重复 or 隐藏真正的复杂性。

❌ 禁止的过度封装：
```js
// 没有意义的二次封装，只是改了个名字
function myGet(url) {
  return request({ method: 'GET', url })
}
function myPost(url, data) {
  return request({ method: 'POST', url, data })
}

// 没必要的 wrapper 组件（只是传了 props）
// MyButton.vue 内部只是 <uni-button v-bind="$attrs" />
```

✅ 允许的封装：
```js
// 封装真正的业务逻辑（拦截器、token、错误处理）
const http = createRequest({
  baseURL: config.apiBase,
  interceptors: { request: addToken, response: handleError }
})
```

**封装层数上限**：最多 2 层（业务层 → 基础层），禁止为了"优雅"而增加层数。

---

### 原则三：代码可读性优先于简洁性

❌ 禁止为了少写代码而牺牲可读性：
```js
// 难以理解
const result = data?.list?.reduce((a, b) => ({ ...a, [b.id]: b }), {}) ?? {}

// 应写成
const list = data?.list ?? []
const resultMap = {}
list.forEach(item => {
  resultMap[item.id] = item
})
```

✅ 变量命名必须是自解释的，不允许无意义的缩写：
```js
// ❌
const u = getUser()
const fn = () => {}
const d = new Date()

// ✅
const currentUser = getUser()
const handleSubmit = () => {}
const today = new Date()
```

---

## 项目目录结构

```
project/
├── src/
│   ├── pages/               # 页面，每页一个目录
│   │   ├── home/
│   │   │   ├── index.vue    # 页面主文件（< 300 行）
│   │   │   ├── composables/ # 该页面专属逻辑（可选）
│   │   │   │   └── useHomeList.js
│   │   │   └── components/  # 该页面专属组件（可选）
│   │   │       └── HomeCard.vue
│   │   └── order/
│   │       ├── index.vue
│   │       └── detail.vue
│   │
│   ├── components/          # 全局公共组件
│   │   ├── BaseButton.vue
│   │   └── BaseModal.vue
│   │
│   ├── composables/         # 全局公共 composable
│   │   ├── useLoading.js
│   │   └── usePagination.js
│   │
│   ├── stores/              # Pinia stores，按业务模块分文件
│   │   ├── user.js
│   │   ├── cart.js
│   │   └── app.js
│   │
│   ├── api/                 # API 请求，按业务模块分文件
│   │   ├── request.js       # axios/uni.request 基础封装（仅此一层）
│   │   ├── user.js
│   │   ├── order.js
│   │   └── product.js
│   │
│   ├── utils/               # 工具函数，按用途分文件
│   │   ├── format.js        # 日期、金额、字符串格式化
│   │   ├── validate.js      # 表单校验
│   │   └── storage.js       # 本地存储封装
│   │
│   ├── constants/           # 常量
│   │   ├── routes.js
│   │   └── enums.js
│   │
│   ├── styles/              # 全局样式
│   │   ├── variables.scss
│   │   └── mixins.scss
│   │
│   ├── App.vue
│   ├── main.js
│   └── pages.json
```

---

## 文件行数限制（硬性要求）

| 文件类型 | 最大行数 | 超出时的处理 |
|---------|---------|------------|
| 页面 `.vue` | 300 行 | 拆分 composable 或子组件 |
| 组件 `.vue` | 200 行 | 拆分更小的子组件 |
| Composable `.js` | 150 行 | 按职责拆分多个 composable |
| Store `.js` | 120 行 | 按业务模块拆分 |
| API 文件 `.js` | 100 行 | 按子业务拆分 |
| 工具函数 `.js` | 80 行 | 按用途拆分 |

---

## Vue 单文件规范

### 模板结构顺序（固定）
```vue
<template>
  <!-- 保持简洁，逻辑放 script -->
</template>

<script setup>
// 1. 导入
// 2. Props & Emits 定义
// 3. Store
// 4. 响应式数据
// 5. 计算属性
// 6. 方法
// 7. 生命周期
</script>

<style lang="scss" scoped>
/* 只写当前组件的样式 */
</style>
```

### 详细示例
```vue
<template>
  <view class="user-profile">
    <view v-if="isLoading" class="loading">
      <uni-load-more status="loading" />
    </view>
    <view v-else>
      <UserAvatar :url="userInfo.avatar" :size="80" />
      <text class="username">{{ userInfo.nickname }}</text>
      <button class="edit-btn" @click="handleEdit">编辑资料</button>
    </view>
  </view>
</template>

<script setup>
// 1. 导入
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { getUserProfile } from '@/api/user'
import UserAvatar from './components/UserAvatar.vue'

// 2. Props & Emits
const props = defineProps({
  userId: {
    type: String,
    required: true
  }
})

// 3. Store
const userStore = useUserStore()

// 4. 响应式数据
const isLoading = ref(false)
const userInfo = ref({})

// 5. 计算属性
const isCurrentUser = computed(() => {
  return props.userId === userStore.currentUserId
})

// 6. 方法
async function loadUserProfile() {
  isLoading.value = true
  try {
    const data = await getUserProfile(props.userId)
    userInfo.value = data
  } catch (error) {
    uni.showToast({ title: '加载失败', icon: 'error' })
  } finally {
    isLoading.value = false
  }
}

function handleEdit() {
  uni.navigateTo({ url: '/pages/profile/edit' })
}

// 7. 生命周期
onMounted(() => {
  loadUserProfile()
})
</script>

<style lang="scss" scoped>
.user-profile {
  padding: 32rpx;
}
.username {
  font-size: 32rpx;
  font-weight: 600;
  margin-top: 16rpx;
}
.edit-btn {
  margin-top: 24rpx;
}
</style>
```

---

## API 请求规范

详见 → `references/api-spec.md`

---

## Store（Pinia）规范

详见 → `references/store-spec.md`

---

## Composable 规范

详见 → `references/composable-spec.md`

---

## 组件规范

详见 → `references/component-spec.md`

---

## AI 生成代码的检查清单

在生成每段代码前，AI 必须自问：

- [ ] 这段逻辑放在当前文件合适吗？还是应该新建文件？
- [ ] 当前文件会超过行数限制吗？
- [ ] 有没有创造不必要的抽象层？
- [ ] 变量/函数命名是否自解释？
- [ ] 是否把不相关的功能混在一起？
- [ ] 新建文件了吗？有没有在回复里说明新建了哪些文件？

**当 AI 需要新建文件时，必须明确告知用户**："我新建了以下文件：xxx"。