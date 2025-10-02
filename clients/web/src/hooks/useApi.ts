import { useState } from 'react';
import { api } from '../services/api';

export const useApi = () => {
  const [response, setResponse] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);

  const callInit = async () => {
    setLoading(true);
    try {
      const text = await api.projectInit();
      setResponse(text);
    } catch (error) {
      setResponse('Error: ' + (error as Error).message);
    }
    setLoading(false);
  };

  const checkHealth = async () => {
    setLoading(true);
    try {
      const data = await api.health();
      setResponse(JSON.stringify(data, null, 2));
    } catch (error) {
      setResponse('Error: ' + (error as Error).message);
    }
    setLoading(false);
  };

  return {
    response,
    loading,
    callInit,
    checkHealth
  };
};