import { HttpClient } from './http-client';

export interface ProgressData {
  step: number;
  total: number;
  message: string;
  completed?: boolean;
}

export class ProjectApi {
  constructor(private http: HttpClient) {
    // console.log('ProjectApi constructor called with http:', http);
    this.initStreamSimple = this.initStreamSimple.bind(this);
    this.initStream = this.initStream.bind(this);
    this.init = this.init.bind(this);
  }

  async init(): Promise<string> {
    const response = await this.http.get('/api/project/init');
    return response.text();
  }

  initStreamSimple(
    onData: (data: string) => void,
    onComplete: () => void
  ): EventSource {
    const eventSource = this.http.createEventSource('/api/project/init');

    const cleanup = () => {
      eventSource.close();
      onComplete();
    };

    eventSource.onmessage = (event) => {
      console.log('Raw SSE message:', event.data);
      try {
        const data = JSON.parse(event.data);
        if (data.completed) {
          cleanup();
        } else {
          // If it's JSON but not completion, treat as string
          onData(event.data);
        }
      } catch (error) {
        // Not JSON, treat as string message
        onData(event.data);
      }
    };

    eventSource.onerror = (error) => {
      console.error('EventSource error:', error);
      cleanup();
    };

    return eventSource;
  }

  initStream(
    onData: (data: ProgressData) => void,
    onComplete: () => void
  ): EventSource {
    const eventSource = this.http.createEventSource('/api/project/init-stream');

    const cleanup = () => {
      eventSource.close();
      onComplete();
    };

    eventSource.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.completed) {
          cleanup();
        } else {
          onData(data);
        }
      } catch (error) {
        console.error('Error parsing SSE data:', error);
        cleanup();
      }
    };

    eventSource.onerror = (error) => {
      console.error('EventSource error:', error);
      cleanup();
    };

    return eventSource;
  }
}
