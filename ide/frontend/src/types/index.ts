export interface GraphNode {
  id: string;
  type: string;
  attributes: Record<string, any>;
  position?: { x: number; y: number };
  expiresAt?: string;
}

export interface GraphEdge {
  id: string;
  source: string;
  target: string;
  type: string;
  attributes: Record<string, any>;
  expiresAt?: string;
}

export interface Graph {
  id: string;
  name: string;
  description: string;
  nodes: GraphNode[];
  edges: GraphEdge[];
  hasCycles?: boolean;
}

export interface RedisCommand {
  command: string;
  args: string[];
  timestamp: number;
}

export interface RedisResponse {
  type: 'string' | 'int' | 'array' | 'bulk' | 'null' | 'error';
  value: any;
  timestamp: number;
}

export interface ConsoleEntry {
  id: string;
  command: string;
  response?: RedisResponse;
  timestamp: number;
  status: 'pending' | 'success' | 'error';
}

export interface ConnectionStatus {
  connected: boolean;
  host: string;
  port: number;
  lastPing?: number;
}
