<script setup lang="ts">
import { computed } from 'vue';
import { useAnalyticsStore } from '@/stores/analytics';

const store = useAnalyticsStore();

const bars = computed(() => {
  if (!store.summary || !store.summary.positions.length) return [];
  const total = store.summary.total_value;
  return store.summary.positions.map(p => ({
    name: p.fund_name,
    percent: (p.current_value / total) * 100,
    value: p.current_value
  })).sort((a, b) => b.percent - a.percent);
});
</script>

<template>
  <div class="glass-panel p-4 mt-4">
    <h3 class="text-green mb-4">> 资产分布_</h3>
    <div class="distribution-list">
      <div v-for="bar in bars" :key="bar.name" class="bar-row">
        <div class="bar-info">
          <span>{{ bar.name }}</span>
          <span>{{ bar.percent.toFixed(2) }}%</span>
        </div>
        <div class="bar-track">
          <div class="bar-fill" :style="{ width: bar.percent + '%' }"></div>
        </div>
      </div>
      <div v-if="!bars.length" class="text-muted text-center italic">
        [无数据...]
      </div>
    </div>
  </div>
</template>

<style scoped>
.p-4 { padding: 1rem; }
.mb-4 { margin-bottom: 1rem; }
.mt-4 { margin-top: 1rem; }

.distribution-list {
  display: flex;
  flex-direction: column;
  gap: 0.8rem;
}

.bar-row {
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
}

.bar-info {
  display: flex;
  justify-content: space-between;
  font-size: 0.85rem;
}

.bar-track {
  height: 8px;
  background: rgba(0, 255, 65, 0.1);
  border-radius: 4px;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  background: var(--primary-green);
  box-shadow: 0 0 10px var(--primary-green-glow);
  transition: width 0.5s ease-out;
}

.text-center { text-align: center; }
.italic { font-style: italic; }
</style>
