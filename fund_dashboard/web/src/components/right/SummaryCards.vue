<script setup lang="ts">
import { useAnalyticsStore } from '@/stores/analytics';
import { useAnalytics } from '@/composables/useAnalytics';

const store = useAnalyticsStore();
const { formatMoney, formatRate, getMatrixPnlClass } = useAnalytics();
</script>

<template>
  <div class="summary-cards" v-if="store.summary">
    <div class="card glass-panel fade-in">
      <div class="label">[SYSTEM] 总资产_</div>
      <div class="value text-green">{{ formatMoney(store.summary.total_value) }}</div>
    </div>
    <div class="card glass-panel fade-in" style="animation-delay: 0.1s">
      <div class="label">[SYSTEM] 总投入_</div>
      <div class="value">{{ formatMoney(store.summary.total_cost) }}</div>
    </div>
    <div class="card glass-panel fade-in" style="animation-delay: 0.2s">
      <div class="label">[SYSTEM] 今日盈亏_</div>
      <div class="value" :class="getMatrixPnlClass(store.summary.daily_pnl)">
        {{ store.summary.daily_pnl > 0 ? '+' : '' }}{{ formatMoney(store.summary.daily_pnl) }}
        <span class="rate">({{ formatRate(store.summary.daily_pnl_rate) }})</span>
      </div>
    </div>
    <div class="card glass-panel fade-in" style="animation-delay: 0.3s">
      <div class="label">[SYSTEM] 累计盈亏_</div>
      <div class="value" :class="getMatrixPnlClass(store.summary.total_pnl)">
        {{ store.summary.total_pnl > 0 ? '+' : '' }}{{ formatMoney(store.summary.total_pnl) }}
        <span class="rate">({{ formatRate(store.summary.total_pnl_rate) }})</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.summary-cards {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
  margin-bottom: 1rem;
}

.card {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.label {
  font-size: 0.9rem;
  color: var(--text-muted);
}

.value {
  font-size: 1.5rem;
  font-weight: bold;
}

.rate {
  font-size: 1rem;
  opacity: 0.8;
}
</style>
