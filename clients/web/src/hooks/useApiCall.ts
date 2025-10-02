import { useState } from 'react';

export const useApiCall = () => {
  const [response, setResponse] = useState<string>('');
  const [loading, setLoading]   = useState<boolean>(false);

  const handleApiCall           = async (
    apiCall:    ()          => Promise<any>,
    formatter?: (data: any) => string
  ) => {
    setLoading(true);
    try {
      const result = await apiCall();
      setResponse(formatter ? formatter(result) : result);
    } catch (error) {
      setResponse('Error: ' + (error as Error).message);
    }
    setLoading(false);
  };

  return {
    response,
    loading,
    handleApiCall
  };
};
