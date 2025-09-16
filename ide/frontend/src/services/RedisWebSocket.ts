import { RedisResponse, ConnectionStatus } from '../types';

export class RedisWebSocket {
  private ws: WebSocket | null = null;
  private url: string;
  private reconnectInterval: number = 5000;
  private reconnectAttempts: number = 0;
  private maxReconnectAttempts: number = 10;
  private commandQueue: Map<string, (response: RedisResponse) => void> = new Map();
  private commandId: number = 0;
  
  // Event handlers
  public onConnectionChange: ((status: ConnectionStatus) => void) | null = null;
  public onResponse: ((response: RedisResponse) => void) | null = null;
  public onError: ((error: string) => void) | null = null;

  constructor() {
    // Use environment variable if available, otherwise fall back to current behavior
    if (process.env.REACT_APP_WEBSOCKET_URL) {
      this.url = process.env.REACT_APP_WEBSOCKET_URL;
    } else {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const host = window.location.host;
      this.url = `${protocol}//${host}/ws`;
    }
  }

  public connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      try {
        this.ws = new WebSocket(this.url);
        
        this.ws.onopen = () => {
          console.log('Connected to PathwayDB Redis WebSocket');
          this.reconnectAttempts = 0;
          this.notifyConnectionChange(true);
          resolve();
        };

        this.ws.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data);
            this.handleMessage(data);
          } catch (error) {
            console.error('Failed to parse WebSocket message:', error);
          }
        };

        this.ws.onclose = () => {
          console.log('Disconnected from PathwayDB Redis WebSocket');
          this.notifyConnectionChange(false);
          this.attemptReconnect();
        };

        this.ws.onerror = (error) => {
          console.error('WebSocket error:', error);
          this.onError?.('WebSocket connection error');
          reject(new Error('Failed to connect to WebSocket'));
        };

      } catch (error) {
        reject(error);
      }
    });
  }

  public disconnect(): void {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  public async executeCommand(command: string, args: string[] = []): Promise<RedisResponse> {
    return new Promise((resolve, reject) => {
      if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
        reject(new Error('WebSocket not connected'));
        return;
      }

      const id = (++this.commandId).toString();
      const message = {
        id,
        command,
        args,
        timestamp: Date.now()
      };

      // Store the resolver for this command
      this.commandQueue.set(id, resolve);

      // Send command
      this.ws.send(JSON.stringify(message));

      // Set timeout for command
      setTimeout(() => {
        if (this.commandQueue.has(id)) {
          this.commandQueue.delete(id);
          reject(new Error('Command timeout'));
        }
      }, 10000);
    });
  }

  private handleMessage(data: any): void {
    if (data.id && this.commandQueue.has(data.id)) {
      // This is a response to a specific command
      const resolver = this.commandQueue.get(data.id);
      this.commandQueue.delete(data.id);
      
      const response: RedisResponse = {
        type: data.type || 'string',
        value: data.value,
        timestamp: Date.now()
      };
      
      resolver?.(response);
    } else {
      // This is a broadcast message or notification
      const response: RedisResponse = {
        type: data.type || 'string',
        value: data.value,
        timestamp: Date.now()
      };
      
      this.onResponse?.(response);
    }
  }

  private attemptReconnect(): void {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('Max reconnection attempts reached');
      return;
    }

    this.reconnectAttempts++;
    console.log(`Attempting to reconnect... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);

    setTimeout(() => {
      this.connect().catch(() => {
        // Reconnection failed, will try again
      });
    }, this.reconnectInterval);
  }

  private notifyConnectionChange(connected: boolean): void {
    const [host, port] = this.url.replace('ws://', '').replace('/ws', '').split(':');
    const status: ConnectionStatus = {
      connected,
      host,
      port: parseInt(port),
      lastPing: connected ? Date.now() : undefined
    };
    
    this.onConnectionChange?.(status);
  }

  public isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }

  // Convenience methods for common Redis commands
  public async ping(): Promise<string> {
    const response = await this.executeCommand('PING');
    return response.value;
  }

  public async info(): Promise<string> {
    const response = await this.executeCommand('INFO');
    return response.value;
  }

  // Graph commands
  public async createGraph(name: string, description?: string): Promise<void> {
    const args = description ? [name, description] : [name];
    await this.executeCommand('GRAPH.CREATE', args);
  }

  public async listGraphs(): Promise<string[]> {
    const response = await this.executeCommand('GRAPH.LIST');
    return response.value || [];
  }

  public async getGraph(name: string): Promise<any> {
    const response = await this.executeCommand('GRAPH.GET', [name]);
    return response.value;
  }

  // Node commands
  public async listNodes(graphId: string): Promise<string[]> {
    const response = await this.executeCommand('NODE.LIST', [graphId]);
    return response.value || [];
  }

  public async getNode(graphId: string, nodeId: string): Promise<any> {
    const response = await this.executeCommand('NODE.GET', [graphId, nodeId]);
    return response.value;
  }

  // Edge commands
  public async listEdges(graphId: string): Promise<string[]> {
    const response = await this.executeCommand('EDGE.LIST', [graphId]);
    return response.value || [];
  }

  public async getEdge(graphId: string, edgeId: string): Promise<any> {
    const response = await this.executeCommand('EDGE.GET', [graphId, edgeId]);
    return response.value;
  }
}
