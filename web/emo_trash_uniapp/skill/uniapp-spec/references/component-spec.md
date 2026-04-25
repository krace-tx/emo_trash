# 组件规范

## 组件分类

| 类型 | 位置 | 命名 | 说明 |
|------|------|------|------|
| 全局公共组件 | `src/components/` | `Base` 前缀 | 多业务复用，无业务逻辑 |
| 页面专属组件 | `pages/xxx/components/` | 业务名 | 只在该页面用 |

---

## 全局公共组件

职责：UI 展示，不包含任何业务逻辑，不直接调用 API，不使用 store。

```vue
<!-- src/components/BaseEmpty.vue —— 空状态组件 -->
<template>
  <view class="base-empty">
    <image class="empty-icon" :src="icon" mode="aspectFit" />
    <text class="empty-text">{{ text }}</text>
    <button v-if="buttonText" class="empty-btn" @click="$emit('action')">
      {{ buttonText }}
    </button>
  </view>
</template>

<script setup>
defineProps({
  icon: {
    type: String,
    default: '/static/images/empty.png'
  },
  text: {
    type: String,
    default: '暂无数据'
  },
  buttonText: {
    type: String,
    default: ''
  }
})

defineEmits(['action'])
</script>

<style lang="scss" scoped>
.base-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 80rpx 0;
}
.empty-icon {
  width: 200rpx;
  height: 200rpx;
}
.empty-text {
  font-size: 28rpx;
  color: #999;
  margin-top: 24rpx;
}
.empty-btn {
  margin-top: 32rpx;
  padding: 16rpx 48rpx;
  background-color: #007aff;
  color: white;
  border-radius: 40rpx;
  font-size: 28rpx;
}
</style>
```

---

## 页面专属组件

只在特定页面使用，可以包含该业务的 props 类型，但**不直接调用 API**（数据由父页面传入）。

```vue
<!-- pages/order/components/OrderCard.vue -->
<template>
  <view class="order-card" @click="$emit('click', order.id)">
    <view class="order-header">
      <text class="order-no">订单号：{{ order.orderNo }}</text>
      <text class="order-status" :class="`status-${order.status}`">
        {{ statusText }}
      </text>
    </view>
    <view class="order-products">
      <view v-for="item in order.items" :key="item.id" class="product-item">
        <image :src="item.image" class="product-image" mode="aspectFill" />
        <view class="product-info">
          <text class="product-name">{{ item.name }}</text>
          <text class="product-price">¥{{ item.price }}</text>
        </view>
      </view>
    </view>
    <view class="order-footer">
      <text class="total-price">合计：¥{{ order.totalPrice }}</text>
    </view>
  </view>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  order: {
    type: Object,
    required: true
  }
})

defineEmits(['click'])

const STATUS_MAP = {
  pending: '待付款',
  paid: '待发货',
  shipped: '已发货',
  done: '已完成',
  cancelled: '已取消'
}

const statusText = computed(() => STATUS_MAP[props.order.status] || '未知')
</script>
```

---

## Props 规范

```js
// ✅ 必须定义 type，重要的 prop 标注 required
defineProps({
  userId: {
    type: String,
    required: true
  },
  size: {
    type: Number,
    default: 60
  },
  variant: {
    type: String,
    default: 'default',
    validator: (value) => ['default', 'primary', 'danger'].includes(value)
  }
})

// ❌ 不要裸写 prop 类型
defineProps(['userId', 'size'])  // 无法知道类型和是否必填
```

---

## 组件通信规范

### 父 → 子：Props
```vue
<OrderCard :order="item" />
```

### 子 → 父：Emits
```vue
<!-- 子组件 -->
<script setup>
const emit = defineEmits(['submit', 'cancel'])
function handleSubmit() {
  emit('submit', formData.value)
}
</script>

<!-- 父组件 -->
<OrderForm @submit="handleFormSubmit" @cancel="showForm = false" />
```

### 跨层级：Store（仅跨页面共享数据时使用）

---

## 禁止事项

```vue
<!-- ❌ 组件内直接调用 API -->
<script setup>
// 组件不应该自己调数据，数据应该由父组件通过 props 传入
const data = await getOrderDetail(props.orderId)
</script>

<!-- ❌ 组件内直接操作路由（除非是导航组件） -->
<script setup>
function handleClick() {
  uni.navigateTo({ url: '/pages/detail' })  // 业务组件不要做这种事
}
</script>

<!-- ✅ 通知父组件处理 -->
<script setup>
const emit = defineEmits(['click'])
function handleClick() {
  emit('click', itemId)  // 父组件决定跳转逻辑
}
</script>
```

---

## 样式规范

```scss
// ✅ 组件内样式必须 scoped
<style lang="scss" scoped>
// ✅ 用 rpx 做单位（自动适配屏幕）
.card {
  padding: 24rpx;
  border-radius: 16rpx;
}

// ✅ 颜色使用变量（定义在 styles/variables.scss）
.title {
  color: var(--color-text-primary);
  font-size: 32rpx;
}
</style>

// ❌ 不要写行内样式（难以维护）
<view style="padding: 24rpx; color: #333;">
```