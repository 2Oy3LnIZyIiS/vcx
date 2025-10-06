

export interface ProgressData {
  step:      number;
  total:     number;
  message:   string;
  completed?: boolean;  // Optional field for completion signal
}


class ApiService {
  private readonly baseUrl: string;


  constructor(baseUrl: string = 'http://localhost:9847') {
    this.baseUrl = baseUrl;

    // Bind methods to preserve 'this' context
    this.health            = this.health.bind(this);
    this.ping              = this.ping.bind(this);
    this.projectInit       = this.projectInit.bind(this);
    this.projectInitStream = this.projectInitStream.bind(this);
  }


  async health(): Promise<any> {
    const response = await fetch(`${this.baseUrl}/health`);
    return response.json();
  }


  async ping(): Promise<any> {
    const response = await fetch(`${this.baseUrl}/ping`);
    return response.json();
  }



  async projectInit(): Promise<string> {
    const response = await fetch(`${this.baseUrl}/api/project/init`);
    return response.text();
  }


  projectInitStream( onProgress: (data: ProgressData) => void,
                     onComplete: ()                   => void
                   ): EventSource {
    const eventSource = new EventSource(`${this.baseUrl}/api/project/init-stream`);

    function cleanup() {
      eventSource.close();
      onComplete();
    }

    eventSource.onmessage = function(event) {
      try {
        const data = JSON.parse(event.data);
        if (data.completed) {
          cleanup();
        } else {
          onProgress(data);
        }
      } catch (error) {
        console.error('Error parsing SSE data:', error);
        cleanup();
      }
    };

    eventSource.onerror = function(error) {
      console.error('EventSource error:', error);
      cleanup();
    };

    return eventSource;
  }


}


export const apiService = new ApiService();
