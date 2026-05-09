<template>
  <main class="page dashboard-page">
    <PageHeader :title="t('page.dashboard.title')" :description="t('page.dashboard.description')">
      <template #actions>
        <a-button>{{ t('common.export') }}</a-button>
        <a-button type="primary">{{ t('common.view') }}</a-button>
      </template>
    </PageHeader>

    <section class="dashboard-page__metrics" :aria-label="t('page.dashboard.title')">
      <MetricCard :title="t('page.dashboard.onlineUsers')" :value="128" :description="t('status.processing')" :progress="78">
        <template #icon><UserOutlined /></template>
      </MetricCard>
      <MetricCard :title="t('page.dashboard.activeRoles')" :value="12" :description="t('status.enabled')" :progress="62">
        <template #icon><TeamOutlined /></template>
      </MetricCard>
      <MetricCard :title="t('page.dashboard.menuResources')" :value="36" :description="t('status.enabled')" :progress="86">
        <template #icon><MenuOutlined /></template>
      </MetricCard>
      <MetricCard :title="t('page.dashboard.auditEvents')" :value="942" :description="t('status.success')" :progress="54">
        <template #icon><AuditOutlined /></template>
      </MetricCard>
    </section>

    <section class="dashboard-page__grid">
      <a-card class="dashboard-card dashboard-card--trend" :bordered="false">
        <div class="dashboard-card__header">
          <div>
            <h3>{{ t('page.dashboard.trendTitle') }}</h3>
            <p>{{ t('page.dashboard.trendDescription') }}</p>
          </div>
          <a-tag color="processing">7D</a-tag>
        </div>
        <div class="trend-chart" aria-hidden="true">
          <span v-for="item in trend" :key="item.day" :style="{ height: `${item.value}%` }">
            <small>{{ item.day }}</small>
          </span>
        </div>
      </a-card>

      <a-card class="dashboard-card dashboard-card--quick" :bordered="false">
        <div class="dashboard-card__header">
          <h3>{{ t('page.dashboard.quickTitle') }}</h3>
        </div>
        <div class="quick-links">
          <router-link v-for="entry in entries" :key="entry.path" :to="entry.path">
            <component :is="entry.icon" />
            <span>{{ entry.label }}</span>
          </router-link>
        </div>
      </a-card>

      <a-card class="dashboard-card dashboard-card--todo" :bordered="false">
        <div class="dashboard-card__header">
          <h3>{{ t('page.dashboard.todoTitle') }}</h3>
          <a-badge :count="3" />
        </div>
        <a-list :data-source="todos" size="small">
          <template #renderItem="{ item }">
            <a-list-item>
              <a-list-item-meta :title="item.title" :description="item.description" />
              <StatusTag :status="item.status" />
            </a-list-item>
          </template>
        </a-list>
      </a-card>

      <a-card class="dashboard-card dashboard-card--health" :bordered="false">
        <div class="dashboard-card__header">
          <div>
            <h3>{{ t('page.dashboard.healthTitle') }}</h3>
            <p>{{ t('page.dashboard.healthDescription') }}</p>
          </div>
        </div>
        <a-progress :percent="96" status="active" />
        <div class="health-list">
          <StatusTag status="success" />
          <StatusTag status="processing" />
          <StatusTag status="enabled" />
        </div>
      </a-card>

      <DataTable
        class="dashboard-card--table"
        row-key="id"
        :title="t('page.dashboard.recentTitle')"
        :columns="activityColumns"
        :data-source="activities"
        :pagination="false"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <StatusTag :status="record.status" />
          </template>
        </template>
      </DataTable>
    </section>
  </main>
</template>

<script setup>
import { computed } from 'vue'
import { AuditOutlined, MenuOutlined, TeamOutlined, UserOutlined } from '@ant-design/icons-vue'
import { useI18n } from 'vue-i18n'
import DataTable from '@/components/common/DataTable.vue'
import MetricCard from '@/components/common/MetricCard.vue'
import PageHeader from '@/components/common/PageHeader.vue'
import StatusTag from '@/components/common/StatusTag.vue'

const { t } = useI18n()

const trend = [
  { day: 'Mon', value: 42 },
  { day: 'Tue', value: 58 },
  { day: 'Wed', value: 46 },
  { day: 'Thu', value: 72 },
  { day: 'Fri', value: 86 },
  { day: 'Sat', value: 64 },
  { day: 'Sun', value: 76 }
]

const entries = computed(() => [
  { path: '/system/user', label: t('page.dashboard.entries.users'), icon: UserOutlined },
  { path: '/system/role', label: t('page.dashboard.entries.roles'), icon: TeamOutlined },
  { path: '/system/menu', label: t('page.dashboard.entries.menus'), icon: MenuOutlined },
  { path: '/audit/operation-log', label: t('page.dashboard.entries.audits'), icon: AuditOutlined }
])

const todos = computed(() => [
  { title: t('page.role.title'), description: t('page.role.description'), status: 'processing' },
  { title: t('page.menu.title'), description: t('page.menu.description'), status: 'pending' },
  { title: t('page.loginLog.title'), description: t('page.loginLog.description'), status: 'success' }
])

const activityColumns = computed(() => [
  { title: t('field.module'), dataIndex: 'module' },
  { title: t('field.action'), dataIndex: 'action' },
  { title: t('field.status'), key: 'status', width: 120 },
  { title: t('field.createdAt'), dataIndex: 'time', width: 180 }
])

const activities = computed(() => [
  { id: 1, module: t('page.user.title'), action: t('common.edit'), status: 'success', time: '2026-05-09 09:30' },
  { id: 2, module: t('page.role.title'), action: t('common.save'), status: 'processing', time: '2026-05-09 10:12' },
  { id: 3, module: t('page.file.title'), action: t('common.upload'), status: 'success', time: '2026-05-09 11:05' }
])
</script>

<style scoped lang="scss">
.dashboard-page {
  :deep(.page-header) {
    position: relative;
    overflow: hidden;
    background: var(--color-bg-container);
  }

  :deep(.page-header::before) {
    content: "";
    position: absolute;
    left: 0;
    top: 0;
    width: 4px;
    height: 100%;
    background: var(--color-primary);
  }

  &__metrics {
    display: grid;
    grid-template-columns: repeat(4, minmax(0, 1fr));
    gap: var(--space-4);
  }

  &__grid {
    display: grid;
    grid-template-columns: repeat(12, minmax(0, 1fr));
    grid-template-areas:
      "trend trend trend trend trend trend trend trend quick quick quick quick"
      "trend trend trend trend trend trend trend trend todo todo todo todo"
      "table table table table table table table table health health health health";
    gap: var(--space-4);
    align-items: stretch;
  }
}

.dashboard-card {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-card);

  &--trend {
    grid-area: trend;
    min-height: 420px;
  }

  &--quick {
    grid-area: quick;
  }

  &--todo {
    grid-area: todo;
  }

  &--health {
    grid-area: health;
  }

  &--table {
    grid-area: table;
    min-width: 0;
  }

  &__header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: var(--space-4);
    margin-bottom: var(--space-4);
  }

  h3 {
    margin: 0;
    color: var(--color-text-heading);
    font-size: 16px;
    font-weight: 650;
  }

  p {
    margin: var(--space-1) 0 0;
    color: var(--color-text-secondary);
  }
}

.dashboard-card--table {
  :deep(.ant-table) {
    min-width: 680px;
  }
}

.trend-chart {
  height: 316px;
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  align-items: end;
  gap: var(--space-3);
  padding-top: var(--space-6);
}

.trend-chart span {
  position: relative;
  min-height: 36px;
  border-radius: var(--radius-md) var(--radius-md) 0 0;
  background: linear-gradient(180deg, #69b1ff, var(--color-primary));
}

.trend-chart small {
  position: absolute;
  left: 50%;
  bottom: -24px;
  transform: translateX(-50%);
  color: var(--color-text-secondary);
}

.quick-links {
  display: grid;
  gap: var(--space-2);
}

.quick-links a {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-lg);
  color: var(--color-text);
  text-decoration: none;
  transition: background-color 0.16s ease, border-color 0.16s ease, color 0.16s ease;
}

.quick-links a:hover,
.quick-links a:focus-visible {
  color: var(--color-primary);
  background: #f0f7ff;
  border-color: #bae0ff;
}

.health-list {
  display: flex;
  gap: var(--space-2);
  flex-wrap: wrap;
  margin-top: var(--space-4);
}

@media (max-width: 1180px) {
  .dashboard-page__metrics {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .dashboard-page__grid {
    grid-template-areas:
      "trend trend trend trend trend trend trend trend trend trend trend trend"
      "quick quick quick quick todo todo todo todo health health health health"
      "table table table table table table table table table table table table";
  }

  .dashboard-card--trend {
    min-height: auto;
  }
}

@media (max-width: 640px) {
  .dashboard-page__metrics {
    grid-template-columns: 1fr;
  }

  .dashboard-page__grid {
    grid-template-columns: 1fr;
    grid-template-areas:
      "trend"
      "quick"
      "todo"
      "health"
      "table";
  }

  .trend-chart {
    height: 220px;
  }
}
</style>
