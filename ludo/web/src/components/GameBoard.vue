<script setup lang="ts">
import { computed } from 'vue'
import { useGameStore } from '../stores/gameStore'
import { useWebSocket } from '../composables/useWebSocket'
import PieceToken from './PieceToken.vue'

const store = useGameStore()
const ws = useWebSocket()

const CELL_SIZE = 40
const BOARD_MARGIN = 40

// 52 coordinate positions on the track (15x15 board grid)
const trackCells = [
  // 0: Red start. Right arm bottom row moving left
  {x:14,y:8},{x:13,y:8},{x:12,y:8},{x:11,y:8},{x:10,y:8},{x:9,y:8}, // 0-5
  // Turn DOWN: Bottom arm right row
  {x:8,y:9},{x:8,y:10},{x:8,y:11},{x:8,y:12},{x:8,y:13},{x:8,y:14}, // 6-11
  // Tip
  {x:7,y:14}, // 12
  // 13: Green start. Bottom arm left row moving up
  {x:6,y:14},{x:6,y:13},{x:6,y:12},{x:6,y:11},{x:6,y:10},{x:6,y:9}, // 13-18
  // Turn LEFT: Left arm bottom row
  {x:5,y:8},{x:4,y:8},{x:3,y:8},{x:2,y:8},{x:1,y:8},{x:0,y:8}, // 19-24
  // Tip
  {x:0,y:7}, // 25
  // 26: Yellow start. Left arm top row moving right
  {x:0,y:6},{x:1,y:6},{x:2,y:6},{x:3,y:6},{x:4,y:6},{x:5,y:6}, // 26-31
  // Turn UP: Top arm left row
  {x:6,y:5},{x:6,y:4},{x:6,y:3},{x:6,y:2},{x:6,y:1},{x:6,y:0}, // 32-37
  // Tip
  {x:7,y:0}, // 38
  // 39: Blue start. Top arm right row moving down
  {x:8,y:0},{x:8,y:1},{x:8,y:2},{x:8,y:3},{x:8,y:4},{x:8,y:5}, // 39-44
  // Turn RIGHT: Right arm top row
  {x:9,y:6},{x:10,y:6},{x:11,y:6},{x:12,y:6},{x:13,y:6},{x:14,y:6}, // 45-50
  // Tip
  {x:14,y:7} // 51
]

// Get color of cell based on standard position pattern (Red, Green, Yellow, Blue)
const getFillColor = (idx: number) => {
  // Exception cells (Tip cells)
  if (idx === 12) return 'var(--green-color)'  // Green tip
  if (idx === 25) return 'var(--yellow-color)' // Yellow tip
  if (idx === 38) return 'var(--blue-color)'   // Blue tip
  if (idx === 51) return 'var(--red-color)'    // Red tip

  const colors = ['var(--red-color)', 'var(--green-color)', 'var(--yellow-color)', 'var(--blue-color)']
  return colors[idx % 4]
}

// Corner triangle centers for piece rendering (relative to margin +40)
const triangleCenters: Record<number, {x: number, y: number}> = {
  31: { x: 266, y: 306 }, // Top-Left corner (Blue triangle)
  32: { x: 294, y: 266 }, // Top-Left corner (Red triangle)
  44: { x: 374, y: 266 }, // Top-Right corner (Red triangle)
  45: { x: 414, y: 306 }, // Top-Right corner (Green triangle)
  5:  { x: 414, y: 374 }, // Bottom-Right corner (Green triangle)
  6:  { x: 374, y: 414 }, // Bottom-Right corner (Blue triangle)
  18: { x: 294, y: 426 }, // Bottom-Left corner (Blue triangle)
  19: { x: 266, y: 374 }  // Bottom-Left corner (Yellow triangle)
}

// Calculate absolute center coords for any track/home position
const getCellCoord = (pos: number) => {
  if (pos >= 0 && pos <= 51) {
    if (triangleCenters[pos]) {
      return { x: triangleCenters[pos].x, y: triangleCenters[pos].y }
    }
    const c = trackCells[pos]
    return {
      x: (c.x + 0.5) * CELL_SIZE + BOARD_MARGIN,
      y: (c.y + 0.5) * CELL_SIZE + BOARD_MARGIN
    }
  }
  
  // Home stretch coordinate calculation
  if (pos >= 100 && pos <= 105) { // Red home stretch (right arm, cols 9 to 13)
    const step = pos - 100
    if (step === 5) return { x: 7.5 * CELL_SIZE + BOARD_MARGIN, y: 7.5 * CELL_SIZE + BOARD_MARGIN }
    return { x: (13 - step + 0.5) * CELL_SIZE + BOARD_MARGIN, y: 7.5 * CELL_SIZE + BOARD_MARGIN }
  }
  if (pos >= 110 && pos <= 115) { // Green home stretch (bottom arm, rows 9 to 13)
    const step = pos - 110
    if (step === 5) return { x: 7.5 * CELL_SIZE + BOARD_MARGIN, y: 7.5 * CELL_SIZE + BOARD_MARGIN }
    return { x: 7.5 * CELL_SIZE + BOARD_MARGIN, y: (13 - step + 0.5) * CELL_SIZE + BOARD_MARGIN }
  }
  if (pos >= 120 && pos <= 125) { // Yellow home stretch (left arm, cols 1 to 5)
    const step = pos - 120
    if (step === 5) return { x: 7.5 * CELL_SIZE + BOARD_MARGIN, y: 7.5 * CELL_SIZE + BOARD_MARGIN }
    return { x: (step + 1.5) * CELL_SIZE + BOARD_MARGIN, y: 7.5 * CELL_SIZE + BOARD_MARGIN }
  }
  if (pos >= 130 && pos <= 135) { // Blue home stretch (top arm, rows 1 to 5)
    const step = pos - 130
    if (step === 5) return { x: 7.5 * CELL_SIZE + BOARD_MARGIN, y: 7.5 * CELL_SIZE + BOARD_MARGIN }
    return { x: 7.5 * CELL_SIZE + BOARD_MARGIN, y: (step + 1.5) * CELL_SIZE + BOARD_MARGIN }
  }

  return { x: 300 + BOARD_MARGIN, y: 300 + BOARD_MARGIN }
}

// Get center coordinates of the 4 circles in the player hangar/base (matching 200x200 sizes)
const getBaseCoord = (color: string, id: number) => {
  const localId = id % 4
  const isLeft = (localId % 2) === 0
  const isTop = localId < 2
  
  const dx = isLeft ? 50 : 150
  const dy = isTop ? 50 : 150

  if (color === 'yellow') { // Top-Left
    return { x: 40 + dx, y: 40 + dy }
  }
  if (color === 'blue') { // Top-Right
    return { x: 440 + dx, y: 40 + dy }
  }
  if (color === 'green') { // Bottom-Left
    return { x: 40 + dx, y: 440 + dy }
  }
  if (color === 'red') { // Bottom-Right
    return { x: 440 + dx, y: 440 + dy }
  }
  return { x: 340, y: 340 }
}

const pieceCoords = computed(() => {
  const coords: any[] = []
  const posCounts: Record<number, number> = {}

  if (!store.board || !store.board.players) return []

  store.board.players.forEach(p => {
    p.pieces.forEach(piece => {
      let x, y
      if (piece.position === -1) {
        const coord = getBaseCoord(p.color, piece.id)
        x = coord.x
        y = coord.y
      } else {
        const coord = getCellCoord(piece.position)
        x = coord.x
        y = coord.y
        posCounts[piece.position] = (posCounts[piece.position] || 0) + 1
      }
      
      const isMovable = store.currentTurn === store.myColor && 
                       p.color === store.myColor && 
                       store.movablePieces.includes(piece.id)

      coords.push({
        ...piece,
        x, y,
        isMovable,
        playerColor: p.color
      })
    })
  })

  const stackCounters: Record<number, number> = {}
  return coords.map(c => {
    if (c.position !== -1 && posCounts[c.position] > 1) {
      stackCounters[c.position] = (stackCounters[c.position] || 0) + 1
      return {
        ...c,
        isStacked: true,
        stackIndex: stackCounters[c.position] - 1,
        stackTotal: posCounts[c.position]
      }
    }
    return { ...c, isStacked: false }
  })
})

const movePiece = (pieceId: number) => {
  ws.send('move_piece', { pieceId })
}

// Get the coordinates for polygon of corner triangle cells
const getCornerTrianglePoints = (idx: number) => {
  // Top-Left Corner
  if (idx === 31) return "240,320 280,280 280,320"
  if (idx === 32) return "280,280 320,240 320,280"
  // Top-Right Corner
  if (idx === 44) return "360,240 360,280 400,280"
  if (idx === 45) return "400,280 440,320 400,320"
  // Bottom-Right Corner
  if (idx === 5) return "400,360 440,360 400,400"
  if (idx === 6) return "360,400 400,400 360,440"
  // Bottom-Left Corner
  if (idx === 18) return "280,400 320,440 280,440"
  if (idx === 19) return "240,360 280,360 280,400"
  return ""
}
</script>

<template>
  <div class="board-wrapper">
    <!-- viewBox 0 0 680 680 to accommodate 600x600 grid plus 40px margins on all sides -->
    <svg width="680" height="680" viewBox="0 0 680 680" class="game-board">
      <defs>
        <!-- Dotted fly path arrowheads -->
        <marker id="arrow-red" viewBox="0 0 10 10" refX="6" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse"><path d="M 0 1.5 L 8 5 L 0 8.5 z" fill="var(--red-color)"/></marker>
        <marker id="arrow-yellow" viewBox="0 0 10 10" refX="6" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse"><path d="M 0 1.5 L 8 5 L 0 8.5 z" fill="var(--yellow-color)"/></marker>
        <marker id="arrow-blue" viewBox="0 0 10 10" refX="6" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse"><path d="M 0 1.5 L 8 5 L 0 8.5 z" fill="var(--blue-color)"/></marker>
        <marker id="arrow-green" viewBox="0 0 10 10" refX="6" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse"><path d="M 0 1.5 L 8 5 L 0 8.5 z" fill="var(--green-color)"/></marker>
      </defs>

      <!-- Cream Board Background -->
      <rect x="5" y="5" width="670" height="670" rx="16" fill="#FAF6EE" stroke="#333333" stroke-width="4"/>

      <!-- White main grid backdrop -->
      <rect x="40" y="40" width="600" height="600" fill="#ffffff" stroke="#333333" stroke-width="4"/>

      <!-- Base Areas (4 corners - 200x200 squares with 40px white margins, matching template) -->
      <!-- Yellow Base (Top Left) -->
      <g transform="translate(40, 40)">
        <rect width="200" height="200" fill="var(--yellow-color)" stroke="#333333" stroke-width="3"/>
        <rect x="15" y="15" width="170" height="170" fill="none" stroke="#333333" stroke-width="2"/>
        <circle cx="50" cy="50" r="24" fill="#ffffff" stroke="#333333" stroke-width="2"/>
        <circle cx="150" cy="50" r="24" fill="#ffffff" stroke="#333333" stroke-width="2"/>
        <circle cx="50" cy="150" r="24" fill="#ffffff" stroke="#333333" stroke-width="2"/>
        <circle cx="150" cy="150" r="24" fill="#ffffff" stroke="#333333" stroke-width="2"/>
      </g>

      <!-- Blue Base (Top Right) -->
      <g transform="translate(440, 40)">
        <rect width="200" height="200" fill="var(--blue-color)" stroke="#333333" stroke-width="3"/>
        <rect x="15" y="15" width="170" height="170" fill="none" stroke="#333333" stroke-width="2"/>
        <circle cx="50" cy="50" r="24" fill="#ffffff" stroke="#333333" stroke-width="2"/>
        <circle cx="150" cy="50" r="24" fill="#ffffff" stroke="#333333" stroke-width="2"/>
        <circle cx="50" cy="150" r="24" fill="#ffffff" stroke="#333333" stroke-width="2"/>
        <circle cx="150" cy="150" r="24" fill="#ffffff" stroke="#333333" stroke-width="2"/>
      </g>

      <!-- Green Base (Bottom Left) -->
      <g transform="translate(40, 440)">
        <rect width="200" height="200" fill="var(--green-color)" stroke="#333333" stroke-width="3"/>
        <rect x="15" y="15" width="170" height="170" fill="none" stroke="#333333" stroke-width="2"/>
        <circle cx="50" cy="50" r="24" fill="#ffffff" stroke="#333333" stroke-width="2"/>
        <circle cx="150" cy="50" r="24" fill="#ffffff" stroke="#333333" stroke-width="2"/>
        <circle cx="50" cy="150" r="24" fill="#ffffff" stroke="#333333" stroke-width="2"/>
        <circle cx="150" cy="150" r="24" fill="#ffffff" stroke="#333333" stroke-width="2"/>
      </g>

      <!-- Red Base (Bottom Right) -->
      <g transform="translate(440, 440)">
        <rect width="200" height="200" fill="var(--red-color)" stroke="#333333" stroke-width="3"/>
        <rect x="15" y="15" width="170" height="170" fill="none" stroke="#333333" stroke-width="2"/>
        <circle cx="50" cy="50" r="24" fill="#ffffff" stroke="#333333" stroke-width="2"/>
        <circle cx="150" cy="50" r="24" fill="#ffffff" stroke="#333333" stroke-width="2"/>
        <circle cx="50" cy="150" r="24" fill="#ffffff" stroke="#333333" stroke-width="2"/>
        <circle cx="150" cy="150" r="24" fill="#ffffff" stroke="#333333" stroke-width="2"/>
      </g>

      <!-- Center Square with 4 Triangles -->
      <g transform="translate(280, 280)">
        <!-- Top: Blue -->
        <polygon points="0,0 120,0 60,60" fill="var(--blue-color)" stroke="#333333" stroke-width="2"/>
        <!-- Left: Yellow -->
        <polygon points="0,0 0,120 60,60" fill="var(--yellow-color)" stroke="#333333" stroke-width="2"/>
        <!-- Bottom: Green -->
        <polygon points="0,120 120,120 60,60" fill="var(--green-color)" stroke="#333333" stroke-width="2"/>
        <!-- Right: Red -->
        <polygon points="120,0 120,120 60,60" fill="var(--red-color)" stroke="#333333" stroke-width="2"/>
        
        <!-- Center "终点" (Finish) Circle -->
        <circle cx="60" cy="60" r="24" fill="var(--blue-color)" stroke="#ffffff" stroke-width="3"/>
        <text x="60" y="65" text-anchor="middle" font-size="13" fill="#ffffff" font-weight="bold">终点</text>
      </g>

      <!-- Dotted Fly Lines (Cross-board flying paths) -->
      <g opacity="0.85">
        <!-- Red fly line: cell 16 to 28 (Bottom arm left row to Left arm top row) -->
        <path d="M 300 500 L 140 300" stroke="var(--red-color)" stroke-width="4.5" stroke-dasharray="8,8" fill="none" marker-end="url(#arrow-red)"/>
        <!-- Green fly line: cell 29 to 41 (Left arm top row to Top arm right row) -->
        <path d="M 180 300 L 380 140" stroke="var(--green-color)" stroke-width="4.5" stroke-dasharray="8,8" fill="none" marker-end="url(#arrow-green)"/>
        <!-- Yellow fly line: cell 42 to 2 (Top arm right row to Right arm bottom row) -->
        <path d="M 380 180 L 540 380" stroke="var(--yellow-color)" stroke-width="4.5" stroke-dasharray="8,8" fill="none" marker-end="url(#arrow-yellow)"/>
        <!-- Blue fly line: cell 3 to 15 (Right arm bottom row to Bottom arm left row) -->
        <path d="M 480 380 L 300 540" stroke="var(--blue-color)" stroke-width="4.5" stroke-dasharray="8,8" fill="none" marker-end="url(#arrow-blue)"/>
      </g>

      <!-- Home Stretch Paths (White cells with colored chevrons) -->
      <!-- Red (Right arm, cols 9 to 13, row 7) -->
      <g v-for="step in 5" :key="'hr'+step" :transform="`translate(${(14 - step) * 40 + 40}, ${7 * 40 + 40})`">
        <rect width="40" height="40" fill="#ffffff" stroke="#333333" stroke-width="2"/>
        <path d="M 24 10 L 14 20 L 24 30" stroke="var(--red-color)" stroke-width="3" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
      </g>
      <!-- Green (Bottom arm, col 7, rows 9 to 13) -->
      <g v-for="step in 5" :key="'hg'+step" :transform="`translate(${7 * 40 + 40}, ${(14 - step) * 40 + 40})`">
        <rect width="40" height="40" fill="#ffffff" stroke="#333333" stroke-width="2"/>
        <path d="M 10 24 L 20 14 L 30 24" stroke="var(--green-color)" stroke-width="3" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
      </g>
      <!-- Yellow (Left arm, row 7, cols 1 to 5) -->
      <g v-for="step in 5" :key="'hy'+step" :transform="`translate(${(step) * 40 + 40}, ${7 * 40 + 40})`">
        <rect width="40" height="40" fill="#ffffff" stroke="#333333" stroke-width="2"/>
        <path d="M 16 10 L 26 20 L 16 30" stroke="var(--yellow-color)" stroke-width="3" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
      </g>
      <!-- Blue (Top arm, col 7, rows 1 to 5) -->
      <g v-for="step in 5" :key="'hb'+step" :transform="`translate(${7 * 40 + 40}, ${(step) * 40 + 40})`">
        <rect width="40" height="40" fill="#ffffff" stroke="#333333" stroke-width="2"/>
        <path d="M 10 16 L 20 26 L 30 16" stroke="var(--blue-color)" stroke-width="3" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
      </g>

      <!-- 52 Main Track Cells -->
      <g v-for="(cell, i) in trackCells" :key="'t'+i">
        <!-- Render as Triangle if it's a corner cell, otherwise render as standard square -->
        <g v-if="getCornerTrianglePoints(i)" :transform="`translate(0, 0)`">
          <polygon :points="getCornerTrianglePoints(i)" :fill="getFillColor(i)" stroke="#333333" stroke-width="2"/>
          <circle :cx="triangleCenters[i].x" :cy="triangleCenters[i].y" r="11" fill="#ffffff" stroke="#333333" stroke-width="1.5" stroke-opacity="0.15"/>
        </g>
        <g v-else :transform="`translate(${cell.x * 40 + 40}, ${cell.y * 40 + 40})`">
          <!-- Main Cell Colored Base -->
          <rect width="40" height="40" :fill="getFillColor(i)" stroke="#333333" stroke-width="2"/>
          <!-- Inner white circle standard for Ludo -->
          <circle cx="20" cy="20" r="12" fill="#ffffff" stroke="#333333" stroke-width="1.5" stroke-opacity="0.15"/>
        </g>
      </g>

      <!-- Start Circles outside the grid with arrows pointing in (strictly aligned) -->
      <!-- Yellow Start Circle & Arrow (row 6, left) -->
      <g>
        <circle cx="20" cy="300" r="17" fill="var(--yellow-color)" stroke="#333333" stroke-width="2"/>
        <text x="20" y="304" text-anchor="middle" font-size="11" fill="#333333" font-weight="bold">起点</text>
        <path d="M 37 300 L 52 300 M 46 295 L 52 300 L 46 305" stroke="var(--yellow-color)" stroke-width="3" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
      </g>

      <!-- Blue Start Circle & Arrow (col 8, top) -->
      <g>
        <circle cx="380" cy="20" r="17" fill="var(--blue-color)" stroke="#333333" stroke-width="2"/>
        <text x="380" y="24" text-anchor="middle" font-size="11" fill="#ffffff" font-weight="bold">起点</text>
        <path d="M 380 37 L 380 52 M 375 46 L 380 52 L 385 46" stroke="var(--blue-color)" stroke-width="3" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
      </g>

      <!-- Red Start Circle & Arrow (row 8, right) -->
      <g>
        <circle cx="660" cy="380" r="17" fill="var(--red-color)" stroke="#333333" stroke-width="2"/>
        <text x="660" y="384" text-anchor="middle" font-size="11" fill="#ffffff" font-weight="bold">起点</text>
        <path d="M 643 380 L 628 380 M 634 375 L 628 380 L 634 385" stroke="var(--red-color)" stroke-width="3" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
      </g>

      <!-- Green Start Circle & Arrow (col 6, bottom) -->
      <g>
        <circle cx="300" cy="660" r="17" fill="var(--green-color)" stroke="#333333" stroke-width="2"/>
        <text x="300" y="664" text-anchor="middle" font-size="11" fill="#ffffff" font-weight="bold">起点</text>
        <path d="M 300 643 L 300 628 M 295 634 L 300 628 L 305 634" stroke="var(--green-color)" stroke-width="3" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
      </g>

      <!-- Render Pieces on top of the board layout -->
      <PieceToken
        v-for="p in pieceCoords"
        :key="`${p.playerColor}-${p.id}`"
        :piece="p"
        :x="p.x"
        :y="p.y"
        :movable="p.isMovable"
        :is-stacked="p.isStacked"
        :stack-index="p.stackIndex"
        :stack-total="p.stackTotal"
        @move="movePiece"
      />
    </svg>
  </div>
</template>

<style scoped>
.board-wrapper {
  background: #fdf2e9;
  padding: 15px;
  border-radius: 16px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.35);
  display: flex;
  justify-content: center;
  align-items: center;
}

.game-board {
  max-width: 100%;
  max-height: 100%;
}
</style>
