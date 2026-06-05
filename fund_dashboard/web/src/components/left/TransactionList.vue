<script setup lang="ts">
import { onMounted } from 'vue';
import { useTransactionStore } from '@/stores/transaction';
import { useAnalytics } from '@/composables/useAnalytics';

const transactionStore = useTransactionStore();
const { formatMoney } = useAnalytics();

onMounted(() => {
  transactionStore.fetchTransactions();
});
</script>

<template>
  <div class="glass-panel p-4 list-container">
    <h3 class="text-green mb-4">> 历史指令日志_</h3>
    <div class="table-wrapper">
      <table v-if="!transactionStore.loading && transactionStore.transactions.length">
        <thead>
          <tr>
            <th>日期</th>
            <th>代码</th>
            <th>类型</th>
            <th>金额</th>
            <th>份额</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="tx in transactionStore.transactions" :key="tx.id" class="fade-in">
            <td>{{ tx.date }}</td>
            <td>{{ tx.fund_code }}</td>
            <td :class="tx.type === 'buy' ? 'text-green' : 'text-red'">
              {{ tx.type === 'buy' ? '买入' : '卖出' }}
            </td>
            <td>{{ formatMoney(tx.amount) }}</td>
            <td>{{ tx.shares.toFixed(2) }}</td>
          </tr>
        </tbody>
      </table>
      <div v-else class="empty-state text-muted">
        [无交易记录...]
      </div>
    </div>
  </div>
</template>

<style scoped>
.p-4 { padding: 1rem; }
.mb-4 { margin-bottom: 1rem; }
.list-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.table-wrapper {
  overflow-y: auto;
  flex: 1;
}

table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9rem;
}

th, td {
  padding: 0.5rem;
  text-align: left;
  border-bottom: 1px solid rgba(0, 255, 65, 0.1);
}

th {
  color: var(--text-muted);
  position: sticky;
  top: 0;
  background: var(--bg-panel);
  backdrop-filter: blur(5px);
}

tr:hover {
  background: rgba(0, 255, 65, 0.05);
}

.empty-state {
  text-align: center;
  padding: 2rem;
  font-style: italic;
}
</style>
