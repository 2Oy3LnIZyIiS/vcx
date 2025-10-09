
import { HttpClient } from './http-client';
import { ProjectApi } from './project-api';
import { HealthApi  } from './health-api';

class ApiService {
  public readonly project: ProjectApi;
  public readonly health: HealthApi;

  constructor(baseUrl?: string) {
    const http   = new HttpClient(baseUrl);

    this.project = new ProjectApi(http);
    this.health  = new HealthApi(http);
  }
}

export const apiService = new ApiService();
export * from './project-api';  // Export types
export * from './health-api';
export * from './http-client';
