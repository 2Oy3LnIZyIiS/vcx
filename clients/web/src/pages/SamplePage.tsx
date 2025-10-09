import React from 'react';
import { Button, Group, Stack, Title, Code, Progress, Text } from '@mantine/core';
import { apiService, ProgressData                          } from '../services/api';
import { useApiCall } from '../hooks/useApiCall';
import { useStreamingApi } from '../hooks/useStreamingApi';

function SamplePage() {
    const { data: response, loading, error, execute: handleApiCall } = useApiCall<string>();
    const { streaming, data: progress, error: streamError, executeStream } = useStreamingApi<ProgressData | string>();

    function checkHealth() {
        handleApiCall(apiService.health.check, function(data) {
            return JSON.stringify(data, null, 2);
        });
    }

    function callInit() {
        executeStream(apiService.project.initStreamSimple);
    }

    function callInitStream() {
        executeStream(apiService.project.initStream);
    }


  return (
    <Stack gap="xl" p="xl">
      <Title order={1}>VCX Control Panel</Title>

      <Group gap="md">
        <Button onClick={checkHealth} loading={loading}>
          Check Agent Health
        </Button>

        <Button onClick={callInit} loading={streaming} variant="outline">
          Initialize Project (Simple)
        </Button>

        <Button onClick={callInitStream} loading={streaming} variant="filled">
          Initialize Project (Progress)
        </Button>
      </Group>

      {progress && (
        <Stack gap="sm">
          {typeof progress === 'string' ? (
            <Text size="sm">{progress}</Text>
          ) : (
            <>
              <Text size="sm">{progress.message}</Text>
              <Progress
                value={(progress.step / progress.total) * 100}
                size="lg"
                animated
              />
              <Text size="xs" c="dimmed">
                Step {progress.step} of {progress.total}
              </Text>
            </>
          )}
        </Stack>
      )}

      {(error || streamError) && (
        <Code block color="red">
          {error || streamError}
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
