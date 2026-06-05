<script setup lang="ts">
import { onMounted } from 'vue';
import { useAnalyticsStore } from '@/stores/analytics';
import SummaryCards from './SummaryCards.vue';
import ValueDistribution from './ValueDistribution.vue';
import PnLCalendar from './PnLCalendar.vue';
import FundDetails from './FundDetails.vue';
import PnLLog from './PnLLog.vue';

const store = useAnalyticsStore();

onMounted(() => {
  store.fetchSummary();
});
</script>

<template>
  <div class="analytics-panel">
    <SummaryCards />
    <div class="grid-2">
      <ValueDistribution />
      <PnLCalendar />
    </div>
    <FundDetails />
    <PnLLog />
  </div>
</template>

<style scoped>
.analytics-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow-y: auto;
  padding-right: 0.5rem;
}

.grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}
</style>
