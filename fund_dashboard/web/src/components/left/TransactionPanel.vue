<script setup lang="ts">
import BuySellForm from './BuySellForm.vue';
import TransactionList from './TransactionList.vue';
</script>

<template>
  <div class="transaction-panel">
    <BuySellForm />
    <TransactionList />
  </div>
</template>

<style scoped>
.transaction-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
}
</style>
