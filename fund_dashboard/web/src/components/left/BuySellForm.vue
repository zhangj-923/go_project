<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useFundStore } from '@/stores/fund';
import { useTransactionStore } from '@/stores/transaction';

const fundStore = useFundStore();
const transactionStore = useTransactionStore();

const form = ref({
  fund_code: '',
  type: 'buy' as 'buy' | 'sell',
  date: new Date().toISOString().split('T')[0],
  shares: 0,
  price: 0,
  amount: 0,
  fee: 0,
});

onMounted(() => {
  if (fundStore.funds.length === 0) {
    fundStore.fetchFunds();
  }
});

const submitForm = async () => {
  if (!form.value.fund_code) return alert('请输入或选择基金代码');
  
  await transactionStore.addTransaction({
    ...form.value,
    shares: Number(form.value.shares),
    price: Number(form.value.price),
    amount: Number(form.value.amount),
    fee: Number(form.value.fee),
  });
  
  // reset
  form.value.shares = 0;
  form.value.price = 0;
  form.value.amount = 0;
  form.value.fee = 0;
};
</script>

<template>
  <div class="glass-panel p-4 mb-4">
    <h3 class="text-green mb-4">> 交易指令输入_</h3>
    <form @submit.prevent="submitForm" class="form-grid">
      <div class="form-group">
        <label>基金代码</label>
        <input v-model="form.fund_code" type="text" placeholder="例如: 000001" required />
      </div>
      <div class="form-group">
        <label>交易类型</label>
        <select v-model="form.type">
          <option value="buy">买入 [BUY]</option>
          <option value="sell">卖出 [SELL]</option>
        </select>
      </div>
      <div class="form-group">
        <label>交易日期</label>
        <input v-model="form.date" type="date" required />
      </div>
      <div class="form-group">
        <label>确认净值</label>
        <input v-model="form.price" type="number" step="0.0001" min="0" required />
      </div>
      <div class="form-group">
        <label>确认份额</label>
        <input v-model="form.shares" type="number" step="0.01" min="0" required />
      </div>
      <div class="form-group">
        <label>发生金额</label>
        <input v-model="form.amount" type="number" step="0.01" min="0" required />
      </div>
      <div class="form-group">
        <label>手续费</label>
        <input v-model="form.fee" type="number" step="0.01" min="0" required />
      </div>
      <div class="form-group full-width">
        <button type="submit" class="submit-btn" :class="{ danger: form.type === 'sell' }">
          > 确认执行 [ENTER]
        </button>
      </div>
    </form>
  </div>
</template>

<style scoped>
.p-4 { padding: 1rem; }
.mb-4 { margin-bottom: 1rem; }
.mb-4 { margin-bottom: 1rem; }
h3 { font-size: 1.2rem; }

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.full-width {
  grid-column: 1 / -1;
}

label {
  font-size: 0.9rem;
  color: var(--text-muted);
}

input, select {
  width: 100%;
}

.submit-btn {
  margin-top: 0.5rem;
  padding: 10px;
  width: 100%;
}
</style>
