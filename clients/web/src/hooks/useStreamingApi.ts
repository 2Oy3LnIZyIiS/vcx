import { useState, useCallback } from 'react';

interface StreamState<T = any> {
  streaming: boolean;
  data:      T | null;
  error:     string | null;
}

interface UseStreamingApiReturn<T = any> {
  streaming:     boolean;
  data:          T | null;
  error:         string | null;
  executeStream: (
    streamCall: (onData: (data: T) => void, onComplete: () => void) => EventSource
  ) => EventSource;
  reset:         () => void;
}

export function useStreamingApi<T = any>(): UseStreamingApiReturn<T> {
  const [state, setState] = useState<StreamState<T>>({
    streaming: false,
    data:      null,
    error:     null
  });

  const executeStream = useCallback((
    streamCall: (onData: (data: T) => void,
    onComplete: () => void)        => EventSource
  ): EventSource => {
    setState({ streaming: true, data: null, error: null });

    function onData(data: T) {
      setState(prev => ({ ...prev, data }));
    }

    function onComplete() {
      setState({ streaming: false, data: null, error: null });
    }

    function onError(error: string) {
      setState({ streaming: false, data: null, error });
    }

    try {
      const eventSource = streamCall(onData, onComplete);
      console.log('EventSource created:', eventSource);
      return eventSource;
    } catch (error) {
      console.log('Error creating EventSource:', error);
      const errorMessage = error instanceof Error ? error.message : 'Streaming error occurred';
      onError(errorMessage);
      throw error;
    }
  }, []);

  const reset = useCallback(() => {
    setState({ streaming: false, data: null, error: null });
  }, []);

  return {
    streaming: state.streaming,
    data:      state.data,
    error:     state.error,
    executeStream,
    reset
  };
}
