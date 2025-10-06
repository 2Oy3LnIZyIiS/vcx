import React, { useState                                   } from 'react';
import { Button, Group, Stack, Title, Code, Progress, Text } from '@mantine/core';
import { apiService, ProgressData                          } from '../services/api';
import { useApiCall                                        } from '../hooks/useApiCall';


function SamplePage() {
    const { data: response, loading, error, execute: handleApiCall } = useApiCall<string>();
    const [initProgress, setInitProgress] = useState<ProgressData | null>(null);
    const [isInitializing, setIsInitializing] = useState(false);

    function checkHealth() {
        handleApiCall(apiService.health, function(data) {
            return JSON.stringify(data, null, 2);
        });
    }

    function callInit() {
        handleApiCall(apiService.projectInit);
    }

    function callInitStream() {
        setIsInitializing(true);
        setInitProgress(null);

        function onProgress(data: ProgressData) {
            setInitProgress(data);
        }
        function onComplete() {
            setIsInitializing(false);
            setInitProgress(null);
        }

        const eventSource = apiService.projectInitStream(onProgress, onComplete);
    }


  return (
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

      {error && (
        <Code block color="red">
          {error}
        </Code>
      )}

      {response && (
        <Code block>
          {response}
        </Code>
      )}
    </Stack>
  );
}

export default SamplePage;
