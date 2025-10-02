import React, { useState } from 'react';
import { MantineProvider, Button, Group, Stack, Title, Code, Progress, Text } from '@mantine/core';
import '@mantine/core/styles.css';
import { api } from './services/api';
import { useApiCall } from './hooks/useApiCall';

function App() {
  const { response, loading, handleApiCall } = useApiCall();
  const [initProgress, setInitProgress] = useState<{step: number, total: number, message: string} | null>(null);
  const [isInitializing, setIsInitializing] = useState(false);

  const checkHealth = () => handleApiCall(api.health, (data) => JSON.stringify(data, null, 2));
  const callInit = () => handleApiCall(api.projectInit);
  const callInitStream = () => {
    setIsInitializing(true);
    setInitProgress(null);

    api.projectInitStream(
      (data) => {
        setInitProgress(data);
      },
      () => {
        setIsInitializing(false);
        setInitProgress(null);
      }
    );
  };


  return (
    <MantineProvider>
      <Stack gap="xl" p="xl">
        <Title order={1}>VCX Control Panel</Title>

        <Group gap="md">
          <Button onClick={checkHealth} loading={loading}>
            Check Agent Health
          </Button>

          <Button onClick={callInit} loading={loading} variant="outline">
            Initialize Project (Simple)
          </Button>

          <Button onClick={callInitStream} loading={isInitializing} variant="filled">
            Initialize Project (Progress)
          </Button>
        </Group>

        {initProgress && (
          <Stack gap="sm">
            <Text size="sm">{initProgress.message}</Text>
            <Progress
              value={(initProgress.step / initProgress.total) * 100}
              size="lg"
              animated
            />
            <Text size="xs" c="dimmed">
              Step {initProgress.step} of {initProgress.total}
            </Text>
          </Stack>
        )}

        {response && (
          <Code block>
            {response}
          </Code>
        )}
      </Stack>
    </MantineProvider>
  );
}

export default App;
