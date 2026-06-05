<script setup lang="ts">
import TransactionPanel from './components/left/TransactionPanel.vue';
import AnalyticsPanel from './components/right/AnalyticsPanel.vue';
</script>

<template>
  <div class="dashboard-container">
    <header class="dashboard-header glass-panel">
      <h1 class="text-green">> FUND_TERMINAL_V1.0.0</h1>
      <div class="status text-muted">
        [CONNECTION: OK] | [ENCRYPTION: ENABLED]
      </div>
    </header>

    <main class="dashboard-main">
      <div class="left-panel">
        <TransactionPanel />
      </div>
      <div class="right-panel">
        <AnalyticsPanel />
      </div>
    </main>
  </div>
</template>

<style scoped>
.dashboard-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  padding: 1rem;
  gap: 1rem;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
}

h1 {
  font-size: 1.5rem;
  letter-spacing: 2px;
  margin: 0;
}

.status {
  font-size: 0.9rem;
}

.dashboard-main {
  display: grid;
  grid-template-columns: 350px 1fr;
  gap: 1rem;
  flex: 1;
  overflow: hidden;
}

.left-panel {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.right-panel {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
</style>
