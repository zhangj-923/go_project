import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useApi } from '@/composables/useApi';
import type { Transaction } from '@/types';
import { useAnalyticsStore } from './analytics';

export const useTransactionStore = defineStore('transaction', () => {
  const transactions = ref<Transaction[]>([]);
  const loading = ref(false);
  const api = useApi();
  const analyticsStore = useAnalyticsStore();

  const fetchTransactions = async () => {
    loading.value = true;
    try {
      const data = await api.get<Transaction[]>('/api/transactions');
      transactions.value = data || [];
    } catch (e) {
      console.error(e);
    } finally {
      loading.value = false;
    }
  };

  const addTransaction = async (tx: Partial<Transaction>) => {
    await api.post('/api/transactions', tx);
    await fetchTransactions();
    await analyticsStore.fetchSummary();
  };

  return { transactions, loading, fetchTransactions, addTransaction };
});
