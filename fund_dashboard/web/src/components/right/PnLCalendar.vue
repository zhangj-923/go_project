<script setup lang="ts">
// Placeholder for PnL Calendar
</script>

<template>
  <div class="glass-panel p-4 mt-4">
    <h3 class="text-green mb-4">> 盈亏日历_</h3>
    <div class="calendar-grid">
      <div class="day" v-for="i in 30" :key="i" :class="i % 3 === 0 ? 'loss' : i % 2 === 0 ? 'profit' : 'neutral'">
        <div class="date">{{ i }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.p-4 { padding: 1rem; }
.mb-4 { margin-bottom: 1rem; }
.mt-4 { margin-top: 1rem; }

.calendar-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
}

.day {
  aspect-ratio: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(0, 255, 65, 0.2);
  font-size: 0.8rem;
  opacity: 0.7;
}

.day.profit {
  background: rgba(0, 255, 65, 0.2);
  color: var(--primary-green);
}

.day.loss {
  background: rgba(255, 68, 68, 0.2);
  color: var(--loss-red);
  border-color: rgba(255, 68, 68, 0.2);
}

.day.neutral {
  background: rgba(255, 255, 255, 0.05);
}

.day:hover {
  opacity: 1;
  border-color: var(--primary-green);
  cursor: crosshair;
}
</style>
