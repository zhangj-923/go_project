import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useApi } from '@/composables/useApi';
import type { Fund } from '@/types';

export const useFundStore = defineStore('fund', () => {
  const funds = ref<Fund[]>([]);
  const loading = ref(false);
  const api = useApi();

  const fetchFunds = async () => {
    loading.value = true;
    try {
      const data = await api.get<Fund[]>('/api/funds');
      funds.value = data || [];
    } catch (e) {
      console.error(e);
    } finally {
      loading.value = false;
    }
  };

  const addFund = async (fund: Partial<Fund>) => {
    await api.post('/api/funds', fund);
    await fetchFunds();
  };

  return { funds, loading, fetchFunds, addFund };
});
