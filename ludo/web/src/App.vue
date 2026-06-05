<script setup lang="ts">
import { RouterView } from 'vue-router'
</script>

<template>
  <div class="app-container">
    <RouterView />
  </div>
</template>

<style>
:root {
  --bg-color: #0f172a;
  --text-color: #f8fafc;
  --primary-color: #a855f7;
  --primary-hover: #9333ea;
  --red-color: #ef4444;
  --yellow-color: #eab308;
  --blue-color: #3b82f6;
  --green-color: #22c55e;
}

body {
  margin: 0;
  padding: 0;
  font-family: 'Inter', system-ui, Avenir, Helvetica, Arial, sans-serif;
  background-color: var(--bg-color);
  color: var(--text-color);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.app-container {
  width: 100vw;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
}

button {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

button:hover {
  background-color: var(--primary-hover);
  transform: translateY(-2px);
}

button:active {
  transform: translateY(0);
}

button:disabled {
  background-color: #475569;
  cursor: not-allowed;
  transform: none;
}
</style>
