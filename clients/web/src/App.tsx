import React from 'react';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';
import SamplePage from './pages/SamplePage';

function App() {
  return (
    <MantineProvider>
      <SamplePage />
    </MantineProvider>
  );
}

export default App;
