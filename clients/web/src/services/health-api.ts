import { HttpClient } from './http-client';

export class HealthApi {
  constructor(private http: HttpClient) {
    this.check = this.check.bind(this);
    this.ping  = this.ping.bind(this);
  }

  async check(): Promise<any> {
    const response = await this.http.get('/health');
    return response.json();
  }

  async ping(): Promise<any> {
    const response = await this.http.get('/ping');
    return response.json();
  }
}
