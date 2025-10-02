const API_BASE = 'http://localhost:9847';

interface ProgressData {
  step: number;
  total: number;
  message: string;
}

export const api = {
  health: async () => {
    const response = await fetch(`${API_BASE}/health`);
    return response.json();
  },

  projectInit: async () => {
    const response = await fetch(`${API_BASE}/api/project/init`);
    return response.text();
  },

  projectInitStream: (onProgress: (data: any) => void, onComplete: () => void) => {
    const eventSource = new EventSource(`${API_BASE}/api/project/init-stream`);
    
    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.completed) {
        eventSource.close();
        onComplete();
      } else {
        onProgress(data);
      }
    };
    
    eventSource.onerror = () => {
      eventSource.close();
      onComplete();
    };
    
    return eventSource;
  },

  ping: async () => {
    const response = await fetch(`${API_BASE}/ping`);
    return response.json();
  }
};