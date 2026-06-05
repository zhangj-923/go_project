import { ref } from 'vue'
import type { GameMessage } from '../types/game'

// 单例模式 - 全局共享一个 WebSocket 连接
let socket: WebSocket | null = null
let isConnected = ref(false)
let reconnectTimer: number | null = null
let isConnecting = false

export function useWebSocket() {
  const messages = ref<GameMessage[]>([])

  const connect = () => {
    // 防止重复连接
    if (isConnecting || (socket && socket.readyState === WebSocket.CONNECTING)) {
      return
    }

    // 如果已连接，不重复连接
    if (socket && socket.readyState === WebSocket.OPEN) {
      return
    }

    isConnecting = true

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.host}/ws`

    socket = new WebSocket(wsUrl)

    socket.onopen = () => {
      isConnected.value = true
      isConnecting = false
      console.log('WebSocket connected')
    }

    socket.onclose = (event) => {
      isConnected.value = false
      isConnecting = false
      socket = null
      console.log('WebSocket disconnected', event.code)

      // 只在非正常关闭时重连
      if (event.code !== 1000) {
        if (reconnectTimer) {
          clearTimeout(reconnectTimer)
        }
        reconnectTimer = window.setTimeout(connect, 3000)
      }
    }

    socket.onerror = (error) => {
      console.error('WebSocket error:', error)
      isConnecting = false
    }

    socket.onmessage = (event) => {
      try {
        const msg: GameMessage = JSON.parse(event.data)
        messages.value.push(msg)
        handleMessage(msg)
      } catch (e) {
        console.error('Failed to parse message:', e)
      }
    }
  }

  const handleMessage = (msg: GameMessage) => {
    window.dispatchEvent(new CustomEvent('game-message', { detail: msg }))
  }

  const send = (msg: GameMessage) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify(msg))
    } else {
      console.error('WebSocket is not connected, attempting to reconnect...')
      connect()
    }
  }

  const disconnect = () => {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    if (socket) {
      socket.close(1000) // 正常关闭
      socket = null
    }
  }

  // 自动连接
  if (!socket && !isConnecting) {
    connect()
  }

  return {
    isConnected,
    messages,
    send,
    disconnect,
  }
}
