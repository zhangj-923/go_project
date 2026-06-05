<script setup lang="ts">
import { useAnalyticsStore } from '@/stores/analytics';
import { useAnalytics } from '@/composables/useAnalytics';

const store = useAnalyticsStore();
const { formatMoney, formatRate, getMatrixPnlClass } = useAnalytics();
</script>

<template>
  <div class="glass-panel p-4 mt-4">
    <h3 class="text-green mb-4">> 持仓明细_</h3>
    <div class="table-wrapper">
      <table v-if="store.summary?.positions?.length">
        <thead>
          <tr>
            <th>基金名称</th>
            <th>市值/占比</th>
            <th>持仓成本</th>
            <th>最新净值</th>
            <th>今日盈亏</th>
            <th>累计盈亏</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="pos in store.summary.positions" :key="pos.fund_code" class="fade-in">
            <td>
              <div class="code text-muted">{{ pos.fund_code }}</div>
              <div>{{ pos.fund_name }}</div>
            </td>
            <td>
              <div>{{ formatMoney(pos.current_value) }}</div>
              <div class="text-muted">{{ pos.total_shares.toFixed(2) }} 份</div>
            </td>
            <td>
              <div>{{ formatMoney(pos.total_cost) }}</div>
              <div class="text-muted">均价 {{ pos.average_cost.toFixed(4) }}</div>
            </td>
            <td>
              <div>{{ pos.current_nav.toFixed(4) }}</div>
              <div class="text-muted">{{ pos.nav_date }}</div>
            </td>
            <td :class="getMatrixPnlClass(pos.daily_pnl)">
              <div>{{ pos.daily_pnl > 0 ? '+' : '' }}{{ formatMoney(pos.daily_pnl) }}</div>
              <div>{{ formatRate(pos.daily_pnl_rate) }}</div>
            </td>
            <td :class="getMatrixPnlClass(pos.total_pnl)">
              <div>{{ pos.total_pnl > 0 ? '+' : '' }}{{ formatMoney(pos.total_pnl) }}</div>
              <div>{{ formatRate(pos.total_pnl_rate) }}</div>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-else class="empty-state text-muted">
        [暂无持仓...]
      </div>
    </div>
  </div>
</template>

<style scoped>
.p-4 { padding: 1rem; }
.mb-4 { margin-bottom: 1rem; }
.mt-4 { margin-top: 1rem; }

.table-wrapper {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9rem;
}

th, td {
  padding: 0.75rem 0.5rem;
  text-align: left;
  border-bottom: 1px solid rgba(0, 255, 65, 0.1);
}

th {
  color: var(--text-muted);
  font-weight: normal;
}

tr:hover {
  background: rgba(0, 255, 65, 0.05);
}

.code {
  font-size: 0.8rem;
}

.empty-state {
  text-align: center;
  padding: 2rem;
  font-style: italic;
}
</style>
