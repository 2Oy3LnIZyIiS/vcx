export class HttpClient {
  constructor(private baseUrl: string = 'http://localhost:9847') {}

  async get(endpoint: string): Promise<Response> {
    return fetch(`${this.baseUrl}${endpoint}`);
  }

  async post(endpoint: string, data?: any): Promise<Response> {
    return fetch(`${this.baseUrl}${endpoint}`,
        { method:  'POST',
          headers: { 'Content-Type': 'application/json' },
          body:    data ? JSON.stringify(data) : undefined
    });
  }

  createEventSource(endpoint: string): EventSource {
    const url         = `${this.baseUrl}${endpoint}`;
    const eventSource = new EventSource(url);
    return eventSource;
  }
}
