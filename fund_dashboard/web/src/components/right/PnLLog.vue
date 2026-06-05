<script setup lang="ts">
import { useAnalyticsStore } from '@/stores/analytics';
import { useAnalytics } from '@/composables/useAnalytics';

const store = useAnalyticsStore();
const { formatMoney, getMatrixPnlClass } = useAnalytics();
</script>

<template>
  <div class="glass-panel p-4 mt-4">
    <h3 class="text-green mb-4">> 系统日志_</h3>
    <div class="log-container">
      <div v-for="pos in store.summary?.positions" :key="pos.fund_code" class="log-line">
        <span class="timestamp">[{{ pos.nav_date || new Date().toISOString().split('T')[0] }}]</span>
        <span class="action">ANALYSIS</span>
        <span class="target">{{ pos.fund_code }}</span>
        <span class="result" :class="getMatrixPnlClass(pos.daily_pnl)">
          PnL: {{ pos.daily_pnl > 0 ? '+' : '' }}{{ formatMoney(pos.daily_pnl) }}
        </span>
      </div>
      <div v-if="!store.summary?.positions?.length" class="log-line text-muted">
        > 等待系统输入...
      </div>
    </div>
  </div>
</template>

<style scoped>
.p-4 { padding: 1rem; }
.mb-4 { margin-bottom: 1rem; }
.mt-4 { margin-top: 1rem; }

.log-container {
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.85rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  max-height: 200px;
  overflow-y: auto;
}

.log-line {
  display: flex;
  gap: 0.5rem;
}

.timestamp { color: var(--text-muted); }
.action { color: #88ccff; }
.target { color: var(--text-main); }
</style>
