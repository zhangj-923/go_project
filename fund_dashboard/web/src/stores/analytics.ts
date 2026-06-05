import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useApi } from '@/composables/useApi';
import type { PortfolioSummary, FundPosition } from '@/types';

export const useAnalyticsStore = defineStore('analytics', () => {
  const summary = ref<PortfolioSummary | null>(null);
  const loading = ref(false);
  const api = useApi();

  const fetchSummary = async () => {
    loading.value = true;
    try {
      const data = await api.get<PortfolioSummary>('/api/analytics/summary');
      summary.value = data;
    } catch (e) {
      console.error(e);
    } finally {
      loading.value = false;
    }
  };

  return { summary, loading, fetchSummary };
});
