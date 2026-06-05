export type Color = 'red' | 'yellow' | 'blue' | 'green';

export interface Piece {
  id: number;
  color: Color;
  position: number;
}

export interface Player {
  id: string;
  name: string;
  color: Color;
  isAi: boolean;
  pieces: Piece[];
  rank: number;
  connected: boolean;
}

export interface BoardState {
  players: Player[];
}

export type MoveEventType = 'move' | 'jump' | 'fly' | 'eat' | 'home';

export interface MoveEvent {
  type: MoveEventType;
  position: number;
  target?: number;
}

// WS Messages
export interface WsMessage<T = any> {
  type: string;
  payload: T;
}

export interface RoomJoinedPayload {
  roomId: string;
  players: Player[];
  myColor: Color;
  myId: string;
  gameMode: string;
}

export interface GameStartedPayload {
  boardState: BoardState;
  currentTurn: Color;
}

export interface DiceRolledPayload {
  playerId: string;
  color: Color;
  value: number;
  movablePieces: number[];
  continuousSix: number;
}

export interface PieceMovedPayload {
  pieceId: number;
  color: Color;
  from: number;
  to: number;
  events: MoveEvent[];
}

export interface TurnChangedPayload {
  playerId: string;
  color: Color;
  extraRoll: boolean;
}

export interface GameOverPayload {
  winnerId: string;
  color: Color;
  rankings: Player[];
}
