import { useState, useCallback } from 'react';

interface ApiState<T = any> {
  data:    T | null;
  loading: boolean;
  error:   string | null;
}

interface UseApiCallReturn<T = any> {
  data:    T | null;
  loading: boolean;
  error:   string | null;
  execute: <R = T>(apiCall: () => Promise<R>, formatter?: (data: R) => T) => Promise<void>;
  reset:   () => void;
}

export function useApiCall<T = string>(): UseApiCallReturn<T> {
  const [state, setState] = useState<ApiState<T>>({
    data:    null,
    loading: false,
    error:   null
  });

  const execute = useCallback(async <R = T>( apiCall: () => Promise<R>,
                                             formatter?: (data: R) => T
                                           ): Promise<void> => {
    setState(prev => ({ ...prev, loading: true, error: null }));

    try {
      const result        = await apiCall();
      const formattedData = formatter ? formatter(result) : (result as unknown as T);
      setState({ data: formattedData, loading: false, error: null });
    } catch (error) {
      const errorMessage  = error instanceof Error ? error.message : 'An unknown error occurred';
      setState({ data: null, loading: false, error: errorMessage });
    }
  }, []);

  const reset = useCallback(() => {
    setState({ data: null, loading: false, error: null });
  }, []);

  return { data:    state.data,
           loading: state.loading,
           error:   state.error,
           execute,
           reset };
}
