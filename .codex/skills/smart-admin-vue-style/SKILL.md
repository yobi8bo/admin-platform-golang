---
name: smart-admin-vue-style
description: Generate SmartAdmin-style Vue 3 admin pages and frontend UI. Use when building pages that should match smart-admin-web-javascript: Vue 3 Composition API, Vite, ant-design-vue, Pinia, dynamic router, query-form + card-table backend management layouts, SmartAdmin spacing classes, theme tokens, and existing /@/api import patterns.
---

# SmartAdmin Vue Style

## 技术栈

使用 Vue 3、`<script setup>`、Vite、ant-design-vue 4、Pinia、vue-router、Less。

使用 `/@/` 作为 `src` 别名。业务接口从 `/@/api/...` 引入，公共常量从 `/@/constants/...` 引入，公共工具从 `/@/lib/...` 或 `/@/utils/...` 引入。

优先使用 ant-design-vue 组件，不自行实现基础 UI：
`a-layout`、`a-form`、`a-row`、`a-col`、`a-card`、`a-table`、`a-pagination`、`a-modal`、`a-drawer`、`a-button`、`a-button-group`、`a-space`、`a-input`、`a-select`、`a-date-picker`、`a-radio-group`。

图标使用全局注册的 `@ant-design/icons-vue`，例如 `SearchOutlined`、`ReloadOutlined`、`PlusOutlined`、`DeleteOutlined`。

## 整体风格

构建克制、密集、偏企业后台的管理系统页面。

视觉特征：
- 页面背景为浅灰 `#f8f8f8`。
- 内容容器以白色 Ant Design 卡片为主。
- 默认主色为 Ant Design 蓝 `#1677ff`。
- 默认左侧菜单为深色 `#001529`。
- 字号以 Ant Design 默认 14px 为主，不使用大标题和营销型 Hero。
- 页面尽量使用表单、表格、分页、弹窗、抽屉等后台组件。

避免：
- 新增全局样式。
- 自定义大面积渐变、装饰图形、营销首页式布局。
- 用普通 `div` 重写 Ant Design 已提供的组件能力。

## 布局模型

项目支持 `side`、`side-expand`、`top`、`top-expand` 多种布局，默认是 `side`。

默认配置特征：
- 左侧菜单宽度 `sideMenuWidth: 200`。
- 头部用户栏高度 `@header-user-height: 40px`。
- 顶部页签默认开启，样式为 `chrome`。
- 面包屑、页脚、水印默认开启。
- 页面内容由动态路由挂载到 `SmartLayout` 内。

业务页面不要重建 Layout、菜单、页签、面包屑、页脚。页面只写当前业务内容。

## 页面结构

列表页优先参考 `src/views/support/config/config-list.vue`。

标准结构：
1. 查询区：`a-form class="smart-query-form"`。
2. 查询行：`a-row class="smart-query-form-row"`。
3. 查询项：`a-form-item class="smart-query-form-item"`。
4. 列表区：`a-card size="small" :bordered="false" :hoverable="true"`。
5. 表格：`a-table size="small" bordered :pagination="false"`。
6. 分页：`a-pagination` 放在 `div.smart-query-table-page` 内。
7. 操作列：`a-button type="link"`，外层使用 `div.smart-table-operate`。
8. 新增/编辑：通常拆成同目录 modal 或 drawer 组件，通过 `ref` + `defineExpose` 打开。

`smart-query-form` 是全局样式约定，不是独立组件。

## 命名约定

查询和表格变量保持脚手架风格：
- `queryFormState`
- `queryForm`
- `tableLoading`
- `tableData`
- `total`
- `columns`
- `resetQuery`
- `onSearch`
- `ajaxQuery` 或 `queryData`

弹窗和抽屉变量：
- `formModal`
- `formRef`
- `visible`
- `formDefault`
- `form`
- `rules`

## 查询区模板

```vue
<a-form class="smart-query-form">
  <a-row class="smart-query-form-row">
    <a-form-item label="关键字" class="smart-query-form-item">
      <a-input style="width: 200px" v-model:value="queryForm.searchWord" placeholder="请输入关键字" />
    </a-form-item>

    <a-form-item class="smart-query-form-item smart-margin-left10">
      <a-button-group>
        <a-button type="primary" @click="onSearch">
          <template #icon>
            <SearchOutlined />
          </template>
          查询
        </a-button>
        <a-button @click="resetQuery">
          <template #icon>
            <ReloadOutlined />
          </template>
          重置
        </a-button>
      </a-button-group>
    </a-form-item>
  </a-row>
</a-form>
```

## 表格区模板

```vue
<a-card size="small" :bordered="false" :hoverable="true">
  <a-row justify="end">
    <TableOperator class="smart-margin-bottom5" v-model="columns" :tableId="TABLE_ID_CONST.xxx" :refresh="ajaxQuery" />
  </a-row>

  <a-table
    size="small"
    :loading="tableLoading"
    bordered
    :dataSource="tableData"
    :columns="columns"
    rowKey="id"
    :pagination="false"
  >
    <template #bodyCell="{ record, column }">
      <template v-if="column.dataIndex === 'action'">
        <div class="smart-table-operate">
          <a-button type="link" @click="toEdit(record)">编辑</a-button>
          <a-button type="link" danger @click="toDelete(record)">删除</a-button>
        </div>
      </template>
    </template>
  </a-table>

  <div class="smart-query-table-page">
    <a-pagination
      showSizeChanger
      showQuickJumper
      show-less-items
      :pageSizeOptions="PAGE_SIZE_OPTIONS"
      v-model:current="queryForm.pageNum"
      v-model:pageSize="queryForm.pageSize"
      :total="total"
      @change="ajaxQuery"
      :show-total="(total) => `共${total}条`"
    />
  </div>
</a-card>
```

## 脚本模板

```vue
<script setup>
  import { onMounted, reactive, ref } from 'vue';
  import { PAGE_SIZE_OPTIONS } from '/@/constants/common-const';
  import { smartSentry } from '/@/lib/smart-sentry';
  import { xxxApi } from '/@/api/xxx/xxx-api';

  const columns = ref([
    { title: '名称', dataIndex: 'name', ellipsis: true },
    { title: '创建时间', dataIndex: 'createTime', width: 150 },
    { title: '操作', dataIndex: 'action', fixed: 'right', width: 100 },
  ]);

  const queryFormState = {
    searchWord: '',
    pageNum: 1,
    pageSize: 10,
  };
  const queryForm = reactive({ ...queryFormState });

  const tableLoading = ref(false);
  const tableData = ref([]);
  const total = ref(0);

  function resetQuery() {
    Object.assign(queryForm, queryFormState);
    ajaxQuery();
  }

  function onSearch() {
    queryForm.pageNum = 1;
    ajaxQuery();
  }

  async function ajaxQuery() {
    try {
      tableLoading.value = true;
      const res = await xxxApi.queryList(queryForm);
      tableData.value = res.data.list;
      total.value = res.data.total;
    } catch (e) {
      smartSentry.captureError(e);
    } finally {
      tableLoading.value = false;
    }
  }

  onMounted(ajaxQuery);
</script>
```

## API 约定

API 文件放在 `/@/api` 下，按模块导出对象：

```js
import { postRequest, getRequest } from '/@/lib/axios';

export const xxxApi = {
  queryList: (param) => postRequest('/xxx/query', param),
  add: (param) => postRequest('/xxx/add', param),
  update: (param) => postRequest('/xxx/update', param),
  delete: (id) => getRequest(`/xxx/delete/${id}`),
};
```

如果接口暂未实现，可以先保留 mock 注释，不要在页面里写大量临时数据逻辑。

## 交互约定

使用 `smartSentry.captureError(e)` 捕获异常。

提交、删除、导入、导出等阻塞操作使用：
`SmartLoading.show()` 和 `SmartLoading.hide()`。

成功反馈使用 `message.success()`。

危险操作使用 `Modal.confirm()`，删除按钮使用 `danger`。

权限按钮使用 `v-privilege`：

```vue
<a-button type="primary" v-privilege="'module:resource:add'">新增</a-button>
<a-button type="link" v-privilege="'module:resource:update'">编辑</a-button>
```

## 样式约定

页面样式使用：

```vue
<style lang="less" scoped>
</style>
```

优先使用现有工具类：
- `smart-width-100`
- `smart-margin-left5`
- `smart-margin-left10`
- `smart-margin-left15`
- `smart-margin-left20`
- `smart-margin-right5`
- `smart-margin-right10`
- `smart-margin-bottom5`
- `smart-margin-bottom10`
- `smart-table-operate`
- `smart-table-btn-block`
- `smart-table-operate-block`
- `smart-table-setting-block`
- `smart-query-table-page`

只有页面确实需要特殊布局时才写 scoped Less。

在布局或可复用组件中需要主题色时，使用 Ant Design token：

```js
import { theme } from 'ant-design-vue';

const { useToken } = theme;
const { token } = useToken();
```

```less
.item:hover {
  color: v-bind('token.colorPrimary');
}
```

## 新页面生成流程

1. 先找同类型页面参考，列表页优先参考 `support/config/config-list.vue`。
2. 复用查询区、卡片表格、分页、操作列结构。
3. 从 `/@/api` 引入接口对象。
4. 使用 `queryForm`、`tableData`、`tableLoading` 等固定命名。
5. 增删改查弹窗拆分为同目录组件。
6. 样式只写 scoped Less，不新增全局样式。
7. 如果需要表格设置，使用 `TableOperator` 和 `TABLE_ID_CONST`。
8. 如果需要字典或枚举，优先用项目内 `DictSelect`、`DictLabel`、`SmartEnumSelect`、`$smartEnumPlugin`。
9. 最后检查页面是否符合后台密集布局：小间距、小表格、白色卡片、右侧分页、link 操作按钮。
